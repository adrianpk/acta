package acta

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestCorrectorAppendsSelfContainedCorrection(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)
	corrector := deterministicCorrector(fixture.journal)
	note := "Replacement significance note"

	result, err := corrector.Correct(context.Background(), CorrectParams{
		TargetID:        fixture.correction.ID,
		Reason:          "Clarify the final behavior",
		ReplacementNote: &note,
	})
	if err != nil {
		t.Fatalf("Correct() error = %v", err)
	}

	event := result.Event
	assertSelfContainedCorrection(t, event, fixture.correction, note)

	effective, ok := fixture.journal.Effective(eventMilestoneKey(event))
	if !ok || effective.ID != event.ID {
		t.Fatalf("Effective() = %#v, %v", effective, ok)
	}
}

func assertSelfContainedCorrection(t *testing.T, event, target Event, note string) {
	t.Helper()

	if event.Type != EventTypeCorrection || event.TargetID != target.ID || event.Note != note {
		t.Fatalf("Correct() event = %#v", event)
	}

	if event.Repository != target.Repository || event.Git != target.Git {
		t.Fatalf("correction evidence = %#v", event)
	}

	if event.Execution == target.Execution || *event.Execution != *target.Execution {
		t.Fatalf("correction execution = %#v", event.Execution)
	}
}

func TestCorrectorCopiesNoteWhenReplacementIsAbsent(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)
	corrector := deterministicCorrector(fixture.journal)

	result, err := corrector.Correct(context.Background(), CorrectParams{
		TargetID: fixture.correction.ID,
		Reason:   "Preserve the effective note",
	})
	if err != nil {
		t.Fatalf("Correct() error = %v", err)
	}

	if result.Event.Note != fixture.correction.Note {
		t.Fatalf("correction note = %q, want %q", result.Event.Note, fixture.correction.Note)
	}
}

func TestCorrectorRejectsMissingStaleAndInvalidTargets(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)

	tests := []struct {
		name   string
		params CorrectParams
		want   error
	}{
		{
			name:   "missing",
			params: CorrectParams{TargetID: "01K1Z8V7BMMN6W7Y6R8QJHA9ZZ", Reason: "Missing target"},
			want:   ErrCorrectionTargetNotFound,
		},
		{
			name:   "stale",
			params: CorrectParams{TargetID: fixture.milestoneA.ID, Reason: "Stale target"},
			want:   ErrCorrectionTargetStale,
		},
		{
			name:   "invalid reason",
			params: CorrectParams{TargetID: fixture.correction.ID, Reason: ""},
			want:   ErrEventInvalid,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			corrector := deterministicCorrector(fixture.journal)
			if _, err := corrector.Correct(context.Background(), test.params); !errors.Is(err, test.want) {
				t.Fatalf("Correct() error = %v, want %v", err, test.want)
			}
		})
	}
}

func deterministicCorrector(journal *Journal) *Corrector {
	corrector := NewCorrector(journal)
	corrector.now = func() time.Time {
		return time.Date(2026, 8, 8, 12, 0, 0, 0, time.UTC)
	}
	corrector.entropy = strings.NewReader(strings.Repeat("\x01", 10))

	return corrector
}
