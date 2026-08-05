package acta

import (
	"context"
	"errors"
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

type queryFilterCase struct {
	name  string
	query Query
	want  []EventID
}

func TestJournalSearchReturnsEffectiveEvents(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)

	result, err := fixture.journal.Search(context.Background(), Query{Limit: 10})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if result.ReturnedCount != 3 || result.Truncated {
		t.Fatalf("Search() result = %#v", result)
	}

	wantIDs := []EventID{fixture.correction.ID, fixture.milestoneC.ID, fixture.milestoneB.ID}
	assertEventIDs(t, result.Events, wantIDs)

	for _, event := range result.Events {
		if event.ID == fixture.milestoneA.ID {
			t.Fatal("superseded milestone was returned")
		}
	}
}

func TestJournalSearchFilters(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)
	recordedFrom := time.Date(2026, 8, 6, 0, 0, 0, 0, time.UTC)
	recordedTo := time.Date(2026, 8, 6, 15, 0, 0, 0, time.UTC)
	authoredFrom := fixture.milestoneB.Git.AuthoredAt.Add(-time.Minute)
	authoredTo := fixture.milestoneB.Git.AuthoredAt.Add(time.Minute)

	tests := append(identityQueryFilterCases(fixture), queryFilterCase{
		name:  "recording interval",
		query: Query{RecordedFrom: &recordedFrom, RecordedTo: &recordedTo},
		want:  []EventID{fixture.milestoneB.ID},
	}, queryFilterCase{
		name:  "author interval",
		query: Query{AuthoredFrom: &authoredFrom, AuthoredTo: &authoredTo},
		want:  []EventID{fixture.milestoneB.ID},
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			test.query.Limit = 10

			result, err := fixture.journal.Search(context.Background(), test.query)
			if err != nil {
				t.Fatalf("Search() error = %v", err)
			}

			assertEventIDs(t, result.Events, test.want)
		})
	}
}

func identityQueryFilterCases(fixture queryFixture) []queryFilterCase {
	return []queryFilterCase{
		{
			name:  "event type",
			query: Query{EventTypes: []EventType{EventTypeCorrection}},
			want:  []EventID{fixture.correction.ID},
		},
		{
			name:  "repository name",
			query: Query{Repository: "EXAMPLE"},
			want:  []EventID{fixture.correction.ID},
		},
		{
			name:  "repository identity",
			query: Query{Repository: fixture.milestoneB.Repository.Identity},
			want:  []EventID{fixture.milestoneB.ID},
		},
		{
			name:  "commit",
			query: Query{Commit: fixture.milestoneA.Git.Commit[:12]},
			want:  []EventID{fixture.correction.ID},
		},
		{
			name:  "execution",
			query: Query{Harness: "PI", Model: "PROVIDER/MODEL"},
			want:  []EventID{fixture.correction.ID},
		},
		{
			name:  "text reason",
			query: Query{Text: "affected runtime"},
			want:  []EventID{fixture.correction.ID},
		},
	}
}

func TestJournalSearchOrdersAndTruncates(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)

	ascending, err := fixture.journal.Search(context.Background(), Query{Limit: 10, Order: QueryAscending})
	if err != nil {
		t.Fatalf("ascending Search() error = %v", err)
	}

	assertEventIDs(t, ascending.Events, []EventID{
		fixture.milestoneB.ID,
		fixture.milestoneC.ID,
		fixture.correction.ID,
	})

	limited, err := fixture.journal.Search(context.Background(), Query{Limit: 1})
	if err != nil {
		t.Fatalf("limited Search() error = %v", err)
	}

	if limited.ReturnedCount != 1 || !limited.Truncated || limited.Events[0].ID != fixture.correction.ID {
		t.Fatalf("limited result = %#v", limited)
	}
}

func TestJournalSearchRejectsAmbiguousCommit(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)

	_, err := fixture.journal.Search(context.Background(), Query{Commit: "aaaa", Limit: 10})
	if !errors.Is(err, ErrCommitAmbiguous) {
		t.Fatalf("Search() error = %v, want ErrCommitAmbiguous", err)
	}
}

func TestJournalSearchNormalizesAndValidatesQuery(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)

	result, err := fixture.journal.Search(context.Background(), Query{
		Harness: "  PI ",
		Text:    " Runtime ",
	})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	normalized := result.Query.Limit == defaultQueryLimit &&
		result.Query.Order == QueryDescending &&
		result.Query.Harness == "pi" &&
		result.Query.Text == "runtime"

	if !normalized {
		t.Fatalf("normalized query = %#v", result.Query)
	}

	now := time.Now()
	tests := []Query{
		{Limit: -1},
		{Limit: maximumQueryLimit + 1},
		{Order: "sideways"},
		{Commit: "not-hex"},
		{EventTypes: []EventType{"unknown"}},
		{RecordedFrom: &now, RecordedTo: &now},
	}

	for _, query := range tests {
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
	milestoneA.Execution = &ExecutionContext{Harness: "pi", Model: "provider/model"}
	milestoneA.RecordedAt = time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)

	milestoneB := distinctMilestoneEvent(t, 2)
	milestoneB.RecordedAt = time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC)
	milestoneB.Repository = RepositoryReference{
		Name:     "other",
		Root:     "/workspace/other",
		Identity: "https://example.test/team/other.git",
		Remote:   "https://example.test/team/other.git",
	}
	milestoneB.Git.Commit = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1"
	milestoneB.Git.Subject = "feat: add transport behavior"
	milestoneB.Git.AuthoredAt = time.Date(2026, 8, 6, 10, 0, 0, 0, time.UTC)
	milestoneB.Git.CommittedAt = time.Date(2026, 8, 6, 11, 0, 0, 0, time.UTC)

	milestoneC := distinctMilestoneEvent(t, 3)
	milestoneC.RecordedAt = time.Date(2026, 8, 6, 18, 0, 0, 0, time.UTC)
	milestoneC.Repository = RepositoryReference{
		Name:     "third",
		Root:     "/workspace/third",
		Identity: "https://example.test/team/third.git",
		Remote:   "https://example.test/team/third.git",
	}
	milestoneC.Git.Commit = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa2"
	milestoneC.Execution = nil

	for _, event := range []Event{milestoneA, milestoneB, milestoneC} {
		appendJournalEvent(t, journal, event)
	}

	correction := milestoneA
	correction.ID = "01K1Z8V7BMMN6W7Y6R8QJHA9ST"
	correction.Type = EventTypeCorrection
	correction.RecordedAt = time.Date(2026, 8, 7, 12, 0, 0, 0, time.UTC)
	correction.Note = "Clarified reusable runtime behavior"
	correction.TargetID = milestoneA.ID
	correction.Reason = "Clarify the affected runtime"
	appendJournalEvent(t, journal, correction)

	return queryFixture{
		journal:    journal,
		milestoneA: milestoneA,
		correction: correction,
		milestoneB: milestoneB,
		milestoneC: milestoneC,
	}
}

func assertEventIDs(t *testing.T, events []Event, want []EventID) {
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
