package acta

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	// SchemaVersion is the initial durable event schema.
	SchemaVersion = "acta/v1"
	maxNoteLength = 500
)

// EventType identifies a durable Acta event variant.
type EventType string

const (
	// EventTypeMilestone records one significant repository commit.
	EventTypeMilestone EventType = "milestone"
	// EventTypeCorrection supersedes one previously recorded event.
	EventTypeCorrection EventType = "correction"
)

// Event is one self-contained canonical journal record.
type Event struct {
	SchemaVersion string              `json:"schema_version"`
	ID            EventID             `json:"id"`
	Type          EventType           `json:"type"`
	RecordedAt    time.Time           `json:"recorded_at"`
	Execution     *ExecutionContext   `json:"execution,omitempty"`
	Repository    RepositoryReference `json:"repository"`
	Git           GitReference        `json:"git"`
	Note          string              `json:"note,omitempty"`
	TargetID      EventID             `json:"target_id,omitempty"`
	Reason        string              `json:"reason,omitempty"`
}

// ExecutionContext identifies the agent environment associated with a commit.
type ExecutionContext struct {
	Harness        string `json:"harness"`
	HarnessVersion string `json:"harness_version,omitempty"`
	Model          string `json:"model,omitempty"`
}

// RepositoryReference identifies a repository and its recorded location.
type RepositoryReference struct {
	Name     string `json:"name"`
	Root     string `json:"root"`
	Identity string `json:"identity"`
	Remote   string `json:"remote,omitempty"`
}

// GitReference contains immutable commit evidence and recorded worktree state.
type GitReference struct {
	Commit        string    `json:"commit"`
	Branch        string    `json:"branch,omitempty"`
	Subject       string    `json:"subject"`
	AuthoredAt    time.Time `json:"authored_at"`
	CommittedAt   time.Time `json:"committed_at"`
	WorktreeClean bool      `json:"worktree_clean"`
}

// Validate checks the complete acta/v1 event contract.
func (event Event) Validate() error {
	if err := event.validateCommonFields(); err != nil {
		return err
	}

	return event.validateTypeFields()
}

func (event Event) validateCommonFields() error {
	if event.SchemaVersion != SchemaVersion {
		return fmt.Errorf("validate event: unsupported schema version %q", event.SchemaVersion)
	}

	if err := event.ID.Validate(); err != nil {
		return fmt.Errorf("validate event ID: %w", err)
	}

	if err := validateUTCTime("recorded_at", event.RecordedAt); err != nil {
		return err
	}

	if event.Execution != nil {
		if err := event.Execution.Validate(); err != nil {
			return err
		}
	}

	if err := event.Repository.Validate(); err != nil {
		return err
	}

	if err := event.Git.Validate(); err != nil {
		return err
	}

	if err := validateSignificanceText("note", event.Note, false); err != nil {
		return err
	}

	return nil
}

func (event Event) validateTypeFields() error {
	switch event.Type {
	case EventTypeMilestone:
		if event.TargetID != "" {
			return fmt.Errorf("validate milestone: target_id must be empty")
		}

		if event.Reason != "" {
			return fmt.Errorf("validate milestone: reason must be empty")
		}
	case EventTypeCorrection:
		if err := event.TargetID.Validate(); err != nil {
			return fmt.Errorf("validate correction target: %w", err)
		}

		if err := validateSignificanceText("reason", event.Reason, true); err != nil {
			return err
		}
	default:
		return fmt.Errorf("validate event: unsupported type %q", event.Type)
	}

	return nil
}

// Validate checks optional execution provenance.
func (context ExecutionContext) Validate() error {
	if err := validateIdentifier("execution harness", context.Harness, true); err != nil {
		return err
	}

	if err := validateIdentifier("execution harness version", context.HarnessVersion, false); err != nil {
		return err
	}

	if err := validateIdentifier("execution model", context.Model, false); err != nil {
		return err
	}

	return nil
}

// Validate checks repository identity and location fields.
func (reference RepositoryReference) Validate() error {
	if err := validateIdentifier("repository name", reference.Name, true); err != nil {
		return err
	}

	if !filepath.IsAbs(reference.Root) {
		return fmt.Errorf("validate repository: root must be absolute")
	}

	if filepath.Clean(reference.Root) != reference.Root {
		return fmt.Errorf("validate repository: root must be clean")
	}

	if err := validateIdentifier("repository identity", reference.Identity, true); err != nil {
		return err
	}

	if err := validateIdentifier("repository remote", reference.Remote, false); err != nil {
		return err
	}

	return nil
}

// Validate checks immutable Git evidence.
func (reference GitReference) Validate() error {
	if !validObjectID(reference.Commit) {
		return fmt.Errorf("validate git reference: commit must be a full lowercase object ID")
	}

	if err := validateIdentifier("git branch", reference.Branch, false); err != nil {
		return err
	}

	if containsControl(reference.Subject) {
		return fmt.Errorf("validate git reference: subject contains a control character")
	}

	if err := validateUTCTime("git authored_at", reference.AuthoredAt); err != nil {
		return err
	}

	if err := validateUTCTime("git committed_at", reference.CommittedAt); err != nil {
		return err
	}

	return nil
}

func validateUTCTime(name string, value time.Time) error {
	if value.IsZero() {
		return fmt.Errorf("validate event: %s is required", name)
	}

	_, offset := value.Zone()
	if offset != 0 {
		return fmt.Errorf("validate event: %s must use UTC", name)
	}

	return nil
}

func validateIdentifier(name, value string, required bool) error {
	if value == "" {
		if required {
			return fmt.Errorf("validate event: %s is required", name)
		}

		return nil
	}

	if strings.TrimSpace(value) != value {
		return fmt.Errorf("validate event: %s must be trimmed", name)
	}

	if containsControl(value) {
		return fmt.Errorf("validate event: %s contains a control character", name)
	}

	return nil
}

func validateSignificanceText(name, value string, required bool) error {
	if err := validateIdentifier(name, value, required); err != nil {
		return err
	}

	if utf8.RuneCountInString(value) > maxNoteLength {
		return fmt.Errorf("validate event: %s exceeds %d Unicode code points", name, maxNoteLength)
	}

	return nil
}

func validObjectID(value string) bool {
	if len(value) != 40 && len(value) != 64 {
		return false
	}

	for _, character := range value {
		if character < '0' || character > '9' {
			if character < 'a' || character > 'f' {
				return false
			}
		}
	}

	return true
}

func containsControl(value string) bool {
	for _, character := range value {
		if unicode.IsControl(character) {
			return true
		}
	}

	return false
}
