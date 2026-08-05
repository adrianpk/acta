package acta

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestRecorderPublishesIdempotentMilestone(t *testing.T) {
	t.Parallel()

	recorder, journal := testRecorder(t)
	params := validRecordParams()

	first, err := recorder.Record(context.Background(), params)
	if err != nil {
		t.Fatalf("first Record() error = %v", err)
	}

	second, err := recorder.Record(context.Background(), params)
	if err != nil {
		t.Fatalf("second Record() error = %v", err)
	}

	assertRecordResults(t, params, first, second)
	assertRecordedMilestone(t, journal, first.ID, params.Note)
}

func assertRecordResults(t *testing.T, params RecordParams, first, second RecordResult) {
	t.Helper()

	if !first.Recorded || second.Recorded || first.ID != second.ID {
		t.Fatalf("record results = %#v, %#v", first, second)
	}

	if first.Repository != params.Repository.Name || first.Commit != params.Git.Commit {
		t.Fatalf("record result = %#v", first)
	}
}

func assertRecordedMilestone(t *testing.T, journal *Journal, id EventID, note string) {
	t.Helper()

	if journal.Count() != 1 {
		t.Fatalf("journal count = %d, want 1", journal.Count())
	}

	event, ok := journal.Event(id)
	if !ok {
		t.Fatal("recorded event is missing")
	}

	wantRecordedAt := time.Date(2026, 8, 5, 18, 0, 0, 0, time.UTC)
	if !event.RecordedAt.Equal(wantRecordedAt) || event.Note != note {
		t.Fatalf("recorded event = %#v", event)
	}
}

func TestRecorderPreservesOptionalExecution(t *testing.T) {
	t.Parallel()

	recorder, journal := testRecorder(t)
	params := validRecordParams()
	params.Execution = &ExecutionContext{
		Harness:        "codex",
		HarnessVersion: "1.2.3",
		Model:          "openai/example-model",
	}

	result, err := recorder.Record(context.Background(), params)
	if err != nil {
		t.Fatalf("Record() error = %v", err)
	}

	event, ok := journal.Event(result.ID)
	if !ok || event.Execution == nil || *event.Execution != *params.Execution {
		t.Fatalf("recorded execution = %#v", event.Execution)
	}
}

func TestRecorderRejectsInvalidEvidence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		change func(*RecordParams)
	}{
		{
			name: "credential-bearing remote",
			change: func(params *RecordParams) {
				params.Repository.Remote = "https://user:secret@example.test/team/example.git"
				params.Repository.Identity = params.Repository.Remote
			},
		},
		{
			name: "noncanonical identity",
			change: func(params *RecordParams) {
				params.Repository.Identity = "https://other.example.test/team/example.git"
			},
		},
		{
			name: "multiline note",
			change: func(params *RecordParams) {
				params.Note = "important\nsecret"
			},
		},
		{
			name: "invalid execution",
			change: func(params *RecordParams) {
				params.Execution = &ExecutionContext{Harness: "codex\nagent"}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			recorder, journal := testRecorder(t)
			params := validRecordParams()
			test.change(&params)

			if _, err := recorder.Record(context.Background(), params); err == nil {
				t.Fatal("expected validation error")
			}

			if journal.Count() != 0 {
				t.Fatalf("journal count = %d, want 0", journal.Count())
			}
		})
	}
}

func TestRecorderReportsEntropyFailure(t *testing.T) {
	t.Parallel()

	recorder, _ := testRecorder(t)
	recorder.entropy = failingReader{}

	if _, err := recorder.Record(context.Background(), validRecordParams()); err == nil {
		t.Fatal("expected entropy error")
	}
}

func TestRecorderRequiresLoadedJournal(t *testing.T) {
	t.Parallel()

	recorder := NewRecorder(NewJournal(testJournalConfig(t)))
	if _, err := recorder.Record(context.Background(), validRecordParams()); !errors.Is(err, ErrJournalNotLoaded) {
		t.Fatalf("Record() error = %v, want ErrJournalNotLoaded", err)
	}
}

func testRecorder(t *testing.T) (*Recorder, *Journal) {
	t.Helper()

	journal := NewJournal(testJournalConfig(t))
	if err := journal.Load(context.Background()); err != nil {
		t.Fatalf("load journal: %v", err)
	}

	recorder := NewRecorder(journal)
	recorder.now = func() time.Time {
		return time.Date(2026, 8, 5, 20, 0, 0, 0, time.FixedZone("test", 2*60*60))
	}
	recorder.entropy = strings.NewReader(strings.Repeat("\x00", 20))

	return recorder, journal
}

func validRecordParams() RecordParams {
	event := validMilestoneEvent()

	return RecordParams{
		Repository: event.Repository,
		Git:        event.Git,
		Note:       event.Note,
	}
}
