package acta

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

type queryFixture struct {
	journal    *Journal
	milestoneA Event
	correction Event
	milestoneB Event
	milestoneC Event
}

func TestJournalSearchReturnsOnlyEffectiveActivities(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)

	result, err := fixture.journal.Search(context.Background(), Query{Limit: 10})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	assertEventIDs(t, result.Events, []EventID{fixture.milestoneC.ID, fixture.milestoneB.ID, fixture.correction.ID})

	if result.ReturnedCount != 3 || result.Truncated {
		t.Fatalf("Search() result = %#v", result)
	}
}

func TestJournalSearchFiltersCanonicalFields(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)
	from := time.Date(2026, 8, 6, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 8, 7, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		query Query
		want  []EventID
	}{
		{name: "type", query: Query{Types: []ActivityType{ActivityTypeCommit}}, want: []EventID{fixture.milestoneC.ID, fixture.milestoneB.ID, fixture.correction.ID}},
		{name: "scope name", query: Query{Scope: "EXAMPLE"}, want: []EventID{fixture.correction.ID}},
		{name: "scope identity", query: Query{Scope: fixture.milestoneB.Scope.Identity}, want: []EventID{fixture.milestoneB.ID}},
		{name: "commit", query: Query{Commit: fixture.milestoneA.Properties.ObjectID[:12]}, want: []EventID{fixture.correction.ID}},
		{name: "provenance", query: Query{Harness: "CODEX", Model: "PROVIDER/MODEL"}, want: []EventID{fixture.correction.ID}},
		{name: "text", query: Query{Text: "affected runtime"}, want: []EventID{fixture.correction.ID}},
		{name: "occurrence interval", query: Query{OccurredFrom: &from, OccurredTo: &to}, want: []EventID{fixture.milestoneC.ID, fixture.milestoneB.ID}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.query.Limit = 10

			result, err := fixture.journal.Search(context.Background(), test.query)
			if err != nil {
				t.Fatalf("Search() error = %v", err)
			}

			assertEventIDs(t, result.Events, test.want)
		})
	}
}

func TestJournalSearchOrdersAndTruncates(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)

	ascending, err := fixture.journal.Search(context.Background(), Query{Limit: 10, Order: QueryAscending})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	assertEventIDs(t, ascending.Events, []EventID{fixture.correction.ID, fixture.milestoneB.ID, fixture.milestoneC.ID})

	limited, err := fixture.journal.Search(context.Background(), Query{Limit: 1})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if limited.ReturnedCount != 1 || !limited.Truncated || limited.Events[0].ID != fixture.milestoneC.ID {
		t.Fatalf("limited result = %#v", limited)
	}
}

func TestQueryResultStopsBeforeResponseByteLimit(t *testing.T) {
	t.Parallel()

	const matchCount = 10

	event := validMilestoneEvent()
	event.Summary = strings.Repeat("x", 1<<20)

	matches := make([]Activity, matchCount)
	for index := range matches {
		matches[index] = event
	}

	result := boundedQueryResult(Query{Limit: matchCount}, matches)
	if !result.Truncated || result.ReturnedCount >= matchCount {
		t.Fatalf("bounded result = %#v", result)
	}

	encoded, err := json.Marshal(SuccessResponse(result))
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}

	if len(encoded)+1 > MaxResponseSize {
		t.Fatalf("response size = %d, maximum = %d", len(encoded)+1, MaxResponseSize)
	}
}

func TestJournalSearchRejectsAmbiguousCommitAndInvalidQueries(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)
	if _, err := fixture.journal.Search(context.Background(), Query{Commit: "aaaa", Limit: 10}); !errors.Is(err, ErrCommitAmbiguous) {
		t.Fatalf("Search() error = %v, want ErrCommitAmbiguous", err)
	}

	now := time.Now()
	for _, query := range []Query{
		{Limit: -1},
		{Limit: maximumQueryLimit + 1},
		{Order: "sideways"},
		{Commit: "not-hex"},
		{Types: []ActivityType{"unknown"}},
		{RecordedFrom: &now, RecordedTo: &now},
		{DurationMin: 20, DurationMax: 10},
	} {
		if _, err := fixture.journal.Search(context.Background(), query); !errors.Is(err, ErrQueryInvalid) {
			t.Fatalf("Search(%#v) error = %v, want ErrQueryInvalid", query, err)
		}
	}
}

func newQueryFixture(t *testing.T) queryFixture {
	t.Helper()

	journal := NewJournal(testJournalConfig(t))
	if err := journal.Load(context.Background()); err != nil {
		t.Fatalf("load journal: %v", err)
	}

	milestoneA := validMilestoneEvent()
	milestoneA.OccurredAt = time.Date(2026, 8, 5, 11, 0, 0, 0, time.UTC)
	milestoneA.RecordedAt = time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)
	milestoneA.ID = nextEventID(t, milestoneA.RecordedAt, 1)
	milestoneA.Provenance = &ActivityProvenance{Type: ProvenanceTypeAgent, Harness: "codex", Model: "provider/model"}

	milestoneB := distinctMilestoneEvent(t, 2)
	milestoneB.OccurredAt = time.Date(2026, 8, 6, 11, 0, 0, 0, time.UTC)
	milestoneB.RecordedAt = time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC)
	milestoneB.ID = nextEventID(t, milestoneB.RecordedAt, 2)
	milestoneB.Scope = &ActivityScope{Type: ScopeTypeRepository, Name: "other", Identity: "https://example.test/team/other.git"}
	milestoneB.Properties.ObjectID = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1"
	milestoneB.Properties.Subject = "feat: add transport behavior"

	milestoneC := distinctMilestoneEvent(t, 3)
	milestoneC.OccurredAt = time.Date(2026, 8, 6, 17, 0, 0, 0, time.UTC)
	milestoneC.RecordedAt = time.Date(2026, 8, 6, 18, 0, 0, 0, time.UTC)
	milestoneC.ID = nextEventID(t, milestoneC.RecordedAt, 3)
	milestoneC.Scope = &ActivityScope{Type: ScopeTypeRepository, Name: "third", Identity: "https://example.test/team/third.git"}
	milestoneC.Properties.ObjectID = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa2"
	milestoneC.Provenance = nil

	for _, event := range []Event{milestoneA, milestoneB, milestoneC} {
		appendJournalEvent(t, journal, event)
	}

	correction := cloneActivity(milestoneA)
	correction.RecordedAt = time.Date(2026, 8, 7, 12, 0, 0, 0, time.UTC)
	correction.ID = nextEventID(t, correction.RecordedAt, 4)
	correction.Summary = "Clarified reusable runtime behavior"
	correction.Supersedes = milestoneA.ID
	correction.CorrectionReason = "Clarify the affected runtime"
	appendJournalEvent(t, journal, correction)

	return queryFixture{journal: journal, milestoneA: milestoneA, correction: correction, milestoneB: milestoneB, milestoneC: milestoneC}
}

func assertEventIDs(t *testing.T, events []Activity, want []EventID) {
	t.Helper()

	if len(events) != len(want) {
		t.Fatalf("event count = %d, want %d: %#v", len(events), len(want), events)
	}

	for index := range want {
		if events[index].ID != want[index] {
			t.Fatalf("event %d ID = %q, want %q", index, events[index].ID, want[index])
		}
	}
}
