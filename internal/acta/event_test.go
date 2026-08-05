package acta

import (
	"strings"
	"testing"
	"time"
)

func TestEventValidate(t *testing.T) {
	t.Parallel()

	valid := validMilestoneEvent()
	tests := []struct {
		name    string
		change  func(*Event)
		wantErr bool
	}{
		{name: "milestone"},
		{
			name: "correction",
			change: func(event *Event) {
				event.Type = EventTypeCorrection
				event.TargetID = "01K1Z8V7BMMN6W7Y6R8QJHA9ST"
				event.Reason = "Clarify the affected runtime"
			},
		},
		{
			name:    "schema",
			change:  func(event *Event) { event.SchemaVersion = "acta/v2" },
			wantErr: true,
		},
		{
			name:    "event ID",
			change:  func(event *Event) { event.ID = "invalid" },
			wantErr: true,
		},
		{
			name:    "recorded offset",
			change:  func(event *Event) { event.RecordedAt = event.RecordedAt.In(time.FixedZone("test", 3600)) },
			wantErr: true,
		},
		{
			name:    "execution harness",
			change:  func(event *Event) { event.Execution.Harness = "pi\nagent" },
			wantErr: true,
		},
		{
			name:    "repository root",
			change:  func(event *Event) { event.Repository.Root = "relative/repo" },
			wantErr: true,
		},
		{
			name:    "commit",
			change:  func(event *Event) { event.Git.Commit = strings.Repeat("A", 40) },
			wantErr: true,
		},
		{
			name:    "author time",
			change:  func(event *Event) { event.Git.AuthoredAt = time.Time{} },
			wantErr: true,
		},
		{
			name:    "note control",
			change:  func(event *Event) { event.Note = "secret\nredaction" },
			wantErr: true,
		},
		{
			name:    "note length",
			change:  func(event *Event) { event.Note = strings.Repeat("x", maxNoteLength+1) },
			wantErr: true,
		},
		{
			name:    "milestone target",
			change:  func(event *Event) { event.TargetID = "01K1Z8V7BMMN6W7Y6R8QJHA9ST" },
			wantErr: true,
		},
		{
			name: "correction target",
			change: func(event *Event) {
				event.Type = EventTypeCorrection
				event.Reason = "Clarify the affected runtime"
			},
			wantErr: true,
		},
		{
			name: "correction reason",
			change: func(event *Event) {
				event.Type = EventTypeCorrection
				event.TargetID = "01K1Z8V7BMMN6W7Y6R8QJHA9ST"
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			event := valid
			execution := *valid.Execution
			event.Execution = &execution
			if test.change != nil {
				test.change(&event)
			}
			err := event.Validate()
			if test.wantErr && err == nil {
				t.Fatal("expected an error")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("Validate() error = %v", err)
			}
		})
	}
}

func TestEventValidateOptionalFields(t *testing.T) {
	t.Parallel()

	event := validMilestoneEvent()
	event.Execution = nil
	event.Repository.Remote = ""
	event.Git.Branch = ""
	event.Note = ""

	if err := event.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func validMilestoneEvent() Event {
	return Event{
		SchemaVersion: SchemaVersion,
		ID:            "01K1Z8V7BMMN6W7Y6R8QJHA9SS",
		Type:          EventTypeMilestone,
		RecordedAt:    time.Date(2026, 8, 5, 14, 22, 31, 418000000, time.UTC),
		Execution: &ExecutionContext{
			Harness:        "pi",
			HarnessVersion: "0.42.0",
			Model:          "openai/gpt-5.6",
		},
		Repository: RepositoryReference{
			Name:     "example",
			Root:     "/workspace/example",
			Identity: "ssh://git@example.test/team/example.git",
			Remote:   "ssh://git@example.test/team/example.git",
		},
		Git: GitReference{
			Commit:        "0123456789abcdef0123456789abcdef01234567",
			Branch:        "dev",
			Subject:       "feat: add reusable agent protocol",
			AuthoredAt:    time.Date(2026, 8, 5, 14, 18, 3, 0, time.UTC),
			CommittedAt:   time.Date(2026, 8, 5, 14, 20, 47, 0, time.UTC),
			WorktreeClean: true,
		},
		Note: "Reusable protocol behavior across agent-enabled projects",
	}
}
