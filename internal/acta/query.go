package acta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	defaultQueryLimit = 20
	maximumQueryLimit = 1000
)

// ErrQueryInvalid identifies an invalid retrieval filter.
var ErrQueryInvalid = errors.New("query is invalid")

// ErrCommitAmbiguous identifies an abbreviated commit with several matches.
var ErrCommitAmbiguous = errors.New("commit is ambiguous")

// QueryOrder selects chronological result direction.
type QueryOrder string

const (
	// QueryAscending orders older events first.
	QueryAscending QueryOrder = "ascending"
	// QueryDescending orders newer events first.
	QueryDescending QueryOrder = "descending"
)

// Query selects a bounded effective-event view.
type Query struct {
	EventTypes    []EventType `json:"event_types,omitempty"`
	Repository    string      `json:"repository,omitempty"`
	Commit        string      `json:"commit,omitempty"`
	Harness       string      `json:"harness,omitempty"`
	Model         string      `json:"model,omitempty"`
	Text          string      `json:"text,omitempty"`
	RecordedFrom  *time.Time  `json:"recorded_from,omitempty"`
	RecordedTo    *time.Time  `json:"recorded_to,omitempty"`
	AuthoredFrom  *time.Time  `json:"authored_from,omitempty"`
	AuthoredTo    *time.Time  `json:"authored_to,omitempty"`
	CommittedFrom *time.Time  `json:"committed_from,omitempty"`
	CommittedTo   *time.Time  `json:"committed_to,omitempty"`
	Limit         int         `json:"limit"`
	Order         QueryOrder  `json:"order"`
}

// QueryResult contains bounded effective events and truncation metadata.
type QueryResult struct {
	Query         Query   `json:"query"`
	Events        []Event `json:"events"`
	ReturnedCount int     `json:"returned_count"`
	Truncated     bool    `json:"truncated"`
}

// Search returns effective events matching one normalized bounded query.
func (journal *Journal) Search(ctx context.Context, query Query) (QueryResult, error) {
	normalized, err := normalizeQuery(query)
	if err != nil {
		return QueryResult{}, err
	}

	journal.mu.RLock()
	defer journal.mu.RUnlock()

	if !journal.loaded {
		return QueryResult{}, ErrJournalNotLoaded
	}

	matches, err := journal.queryMatches(ctx, normalized)
	if err != nil {
		return QueryResult{}, err
	}

	sortQueryEvents(matches, normalized.Order)

	return boundedQueryResult(normalized, matches), nil
}

func (journal *Journal) queryMatches(ctx context.Context, query Query) ([]Event, error) {
	matches := make([]Event, 0)
	commits := make(map[string]struct{})

	for _, id := range journal.order {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		event := journal.events[id]
		if journal.effective[eventMilestoneKey(event)] != id || !eventMatchesQuery(event, query) {
			continue
		}

		matches = append(matches, event)
		if query.Commit != "" {
			commits[event.Git.Commit] = struct{}{}
		}
	}

	if len(commits) > 1 {
		return nil, fmt.Errorf("%w: %q", ErrCommitAmbiguous, query.Commit)
	}

	return matches, nil
}

func normalizeQuery(query Query) (Query, error) {
	query.Repository = strings.TrimSpace(query.Repository)
	query.Commit = strings.ToLower(strings.TrimSpace(query.Commit))
	query.Harness = normalizeQueryText(query.Harness)
	query.Model = normalizeQueryText(query.Model)
	query.Text = normalizeQueryText(query.Text)

	if query.Limit == 0 {
		query.Limit = defaultQueryLimit
	}

	if query.Limit < 1 || query.Limit > maximumQueryLimit {
		return Query{}, fmt.Errorf("%w: limit must be between 1 and %d", ErrQueryInvalid, maximumQueryLimit)
	}

	if query.Order == "" {
		query.Order = QueryDescending
	}

	if query.Order != QueryAscending && query.Order != QueryDescending {
		return Query{}, fmt.Errorf("%w: unsupported order %q", ErrQueryInvalid, query.Order)
	}

	if err := validateCommitPrefix(query.Commit); err != nil {
		return Query{}, err
	}

	if err := normalizeEventTypes(&query); err != nil {
		return Query{}, err
	}

	normalizeQueryTimes(&query)

	if err := validateQueryRanges(query); err != nil {
		return Query{}, err
	}

	return query, nil
}

func normalizeEventTypes(query *Query) error {
	seen := make(map[EventType]struct{})

	for _, eventType := range query.EventTypes {
		if eventType != EventTypeMilestone && eventType != EventTypeCorrection {
			return fmt.Errorf("%w: unsupported event type %q", ErrQueryInvalid, eventType)
		}

		seen[eventType] = struct{}{}
	}

	query.EventTypes = query.EventTypes[:0]
	for eventType := range seen {
		query.EventTypes = append(query.EventTypes, eventType)
	}

	sort.Slice(query.EventTypes, func(left, right int) bool {
		return query.EventTypes[left] < query.EventTypes[right]
	})

	return nil
}

func normalizeQueryTimes(query *Query) {
	query.RecordedFrom = utcTimePointer(query.RecordedFrom)
	query.RecordedTo = utcTimePointer(query.RecordedTo)
	query.AuthoredFrom = utcTimePointer(query.AuthoredFrom)
	query.AuthoredTo = utcTimePointer(query.AuthoredTo)
	query.CommittedFrom = utcTimePointer(query.CommittedFrom)
	query.CommittedTo = utcTimePointer(query.CommittedTo)
}

func utcTimePointer(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}

	utc := value.UTC()

	return &utc
}

func validateQueryRanges(query Query) error {
	if invalidTimeRange(query.RecordedFrom, query.RecordedTo) {
		return fmt.Errorf("%w: recorded range is empty or reversed", ErrQueryInvalid)
	}

	if invalidTimeRange(query.AuthoredFrom, query.AuthoredTo) {
		return fmt.Errorf("%w: authored range is empty or reversed", ErrQueryInvalid)
	}

	if invalidTimeRange(query.CommittedFrom, query.CommittedTo) {
		return fmt.Errorf("%w: committed range is empty or reversed", ErrQueryInvalid)
	}

	return nil
}

func invalidTimeRange(from, to *time.Time) bool {
	return from != nil && to != nil && !from.Before(*to)
}

func validateCommitPrefix(commit string) error {
	if commit == "" {
		return nil
	}

	if len(commit) > 64 {
		return fmt.Errorf("%w: commit prefix is too long", ErrQueryInvalid)
	}

	for _, character := range commit {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return fmt.Errorf("%w: commit prefix must be lowercase hexadecimal", ErrQueryInvalid)
		}
	}

	return nil
}

func eventMatchesQuery(event Event, query Query) bool {
	return matchesEventType(event, query.EventTypes) &&
		matchesRepository(event, query.Repository) &&
		strings.HasPrefix(event.Git.Commit, query.Commit) &&
		matchesExecution(event, query.Harness, query.Model) &&
		matchesText(event, query.Text) &&
		withinTimeRange(event.RecordedAt, query.RecordedFrom, query.RecordedTo) &&
		withinTimeRange(event.Git.AuthoredAt, query.AuthoredFrom, query.AuthoredTo) &&
		withinTimeRange(event.Git.CommittedAt, query.CommittedFrom, query.CommittedTo)
}

func matchesEventType(event Event, eventTypes []EventType) bool {
	if len(eventTypes) == 0 {
		return true
	}

	for _, eventType := range eventTypes {
		if event.Type == eventType {
			return true
		}
	}

	return false
}

func matchesRepository(event Event, repository string) bool {
	if repository == "" {
		return true
	}

	return strings.EqualFold(event.Repository.Name, repository) || strings.EqualFold(event.Repository.Identity, repository)
}

func matchesExecution(event Event, harness, model string) bool {
	if harness == "" && model == "" {
		return true
	}

	if event.Execution == nil {
		return false
	}

	return (harness == "" || normalizeQueryText(event.Execution.Harness) == harness) &&
		(model == "" || normalizeQueryText(event.Execution.Model) == model)
}

func matchesText(event Event, text string) bool {
	if text == "" {
		return true
	}

	values := []string{
		event.Repository.Name,
		event.Git.Subject,
		event.Note,
		event.Reason,
	}
	if event.Execution != nil {
		values = append(values, event.Execution.Harness, event.Execution.Model)
	}

	for _, value := range values {
		if strings.Contains(normalizeQueryText(value), text) {
			return true
		}
	}

	return false
}

func withinTimeRange(value time.Time, from, to *time.Time) bool {
	return (from == nil || !value.Before(*from)) && (to == nil || value.Before(*to))
}

func normalizeQueryText(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func sortQueryEvents(events []Event, order QueryOrder) {
	sort.Slice(events, func(left, right int) bool {
		if events[left].RecordedAt.Equal(events[right].RecordedAt) {
			if order == QueryAscending {
				return events[left].ID < events[right].ID
			}

			return events[left].ID > events[right].ID
		}

		if order == QueryAscending {
			return events[left].RecordedAt.Before(events[right].RecordedAt)
		}

		return events[left].RecordedAt.After(events[right].RecordedAt)
	})
}

func boundedQueryResult(query Query, matches []Event) QueryResult {
	result := QueryResult{Query: query, Events: make([]Event, 0, min(query.Limit, len(matches)))}

	for index, event := range matches {
		if index >= query.Limit {
			result.Truncated = true

			break
		}

		candidate := result
		candidate.Events = append(append([]Event(nil), result.Events...), event)
		candidate.ReturnedCount = len(candidate.Events)
		candidate.Truncated = index+1 < len(matches)

		if !queryResponseFits(candidate) {
			result.Truncated = true

			break
		}

		result = candidate
	}

	result.ReturnedCount = len(result.Events)

	return result
}

func queryResponseFits(result QueryResult) bool {
	response := SuccessResponse(result)
	encoded, err := json.Marshal(response)

	return err == nil && len(encoded)+1 <= MaxResponseSize
}
