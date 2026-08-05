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

const reportDateLayout = "2006-01-02"

// ErrReportInvalid identifies an invalid report interval or timezone.
var ErrReportInvalid = errors.New("report is invalid")

// ReportParams selects an explicit local-calendar recording interval.
type ReportParams struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Timezone string `json:"timezone,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

// ReportInterval describes the effective local-calendar interval.
type ReportInterval struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Timezone string `json:"timezone"`
}

// ReportRepository groups effective events for one repository.
type ReportRepository struct {
	Name     string  `json:"name"`
	Identity string  `json:"identity"`
	Events   []Event `json:"events"`
}

// ReportDay groups repositories by local recording date.
type ReportDay struct {
	Date         string             `json:"date"`
	Repositories []ReportRepository `json:"repositories"`
}

// ReportResult contains one bounded grouped chronicle report.
type ReportResult struct {
	Interval      ReportInterval `json:"interval"`
	Days          []ReportDay    `json:"days"`
	ReturnedCount int            `json:"returned_count"`
	Truncated     bool           `json:"truncated"`
}

type normalizedReport struct {
	interval ReportInterval
	location *time.Location
	from     time.Time
	to       time.Time
	limit    int
}

// Report returns effective events grouped by local date and repository.
func (journal *Journal) Report(ctx context.Context, params ReportParams) (ReportResult, error) {
	normalized, err := normalizeReport(params)
	if err != nil {
		return ReportResult{}, err
	}

	queryResult, err := journal.Search(ctx, Query{
		RecordedFrom: &normalized.from,
		RecordedTo:   &normalized.to,
		Limit:        normalized.limit,
		Order:        QueryAscending,
	})
	if err != nil {
		return ReportResult{}, err
	}

	return boundedReportResult(normalized, queryResult.Events, queryResult.Truncated), nil
}

func normalizeReport(params ReportParams) (normalizedReport, error) {
	location, err := reportLocation(strings.TrimSpace(params.Timezone))
	if err != nil {
		return normalizedReport{}, err
	}

	from, err := parseReportDate("from", params.From, location)
	if err != nil {
		return normalizedReport{}, err
	}

	to, err := parseReportDate("to", params.To, location)
	if err != nil {
		return normalizedReport{}, err
	}

	if !from.Before(to) {
		return normalizedReport{}, fmt.Errorf("%w: from must be before to", ErrReportInvalid)
	}

	limit := params.Limit
	if limit == 0 {
		limit = maximumQueryLimit
	}

	if limit < 1 || limit > maximumQueryLimit {
		return normalizedReport{}, fmt.Errorf("%w: limit must be between 1 and %d", ErrReportInvalid, maximumQueryLimit)
	}

	return normalizedReport{
		interval: ReportInterval{
			From:     from.In(location).Format(reportDateLayout),
			To:       to.In(location).Format(reportDateLayout),
			Timezone: location.String(),
		},
		location: location,
		from:     from.UTC(),
		to:       to.UTC(),
		limit:    limit,
	}, nil
}

func reportLocation(name string) (*time.Location, error) {
	if name == "" {
		return time.Local, nil
	}

	location, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("%w: load timezone %q: %w", ErrReportInvalid, name, err)
	}

	return location, nil
}

func parseReportDate(name, value string, location *time.Location) (time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, fmt.Errorf("%w: %s date is required", ErrReportInvalid, name)
	}

	parsed, err := time.ParseInLocation(reportDateLayout, trimmed, location)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: parse %s date: %w", ErrReportInvalid, name, err)
	}

	return parsed, nil
}

func boundedReportResult(report normalizedReport, events []Event, sourceTruncated bool) ReportResult {
	result := buildReportResult(report, events, sourceTruncated)
	if reportResponseFits(result) {
		return result
	}

	low := 0
	high := len(events)

	for low < high {
		middle := low + (high-low+1)/2

		candidate := buildReportResult(report, events[:middle], true)
		if reportResponseFits(candidate) {
			low = middle
		} else {
			high = middle - 1
		}
	}

	return buildReportResult(report, events[:low], true)
}

func buildReportResult(report normalizedReport, events []Event, truncated bool) ReportResult {
	result := ReportResult{
		Interval:      report.interval,
		Days:          groupReportEvents(events, report.location),
		ReturnedCount: len(events),
		Truncated:     truncated,
	}

	return result
}

func groupReportEvents(events []Event, location *time.Location) []ReportDay {
	days := make([]ReportDay, 0)
	dayIndexes := make(map[string]int)
	repositoryIndexes := make(map[string]map[string]int)

	for _, event := range events {
		date := event.RecordedAt.In(location).Format(reportDateLayout)

		dayIndex, ok := dayIndexes[date]
		if !ok {
			dayIndex = len(days)
			dayIndexes[date] = dayIndex
			repositoryIndexes[date] = make(map[string]int)
			days = append(days, ReportDay{Date: date, Repositories: make([]ReportRepository, 0)})
		}

		repositoryKey := event.Repository.Name + "\x00" + event.Repository.Identity

		repositoryIndex, ok := repositoryIndexes[date][repositoryKey]
		if !ok {
			repositoryIndex = len(days[dayIndex].Repositories)
			repositoryIndexes[date][repositoryKey] = repositoryIndex

			days[dayIndex].Repositories = append(days[dayIndex].Repositories, ReportRepository{
				Name:     event.Repository.Name,
				Identity: event.Repository.Identity,
				Events:   make([]Event, 0),
			})
		}

		repository := &days[dayIndex].Repositories[repositoryIndex]
		repository.Events = append(repository.Events, event)
	}

	for index := range days {
		sort.Slice(days[index].Repositories, func(left, right int) bool {
			leftRepository := days[index].Repositories[left]
			rightRepository := days[index].Repositories[right]

			if leftRepository.Name == rightRepository.Name {
				return leftRepository.Identity < rightRepository.Identity
			}

			return leftRepository.Name < rightRepository.Name
		})
	}

	return days
}

func reportResponseFits(result ReportResult) bool {
	response := SuccessResponse(result)
	encoded, err := json.Marshal(response)

	return err == nil && len(encoded)+1 <= MaxResponseSize
}
