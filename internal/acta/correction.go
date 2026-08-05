package acta

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

// ErrCorrectionTargetNotFound identifies a missing correction target.
var ErrCorrectionTargetNotFound = errors.New("correction target was not found")

// ErrCorrectionTargetStale identifies a target that is no longer effective.
var ErrCorrectionTargetStale = errors.New("correction target is not effective")

// CorrectParams contains the append-only correction input.
type CorrectParams struct {
	TargetID        EventID `json:"target_id"`
	Reason          string  `json:"reason"`
	ReplacementNote *string `json:"replacement_note,omitempty"`
}

// CorrectResult describes a newly appended correction.
type CorrectResult struct {
	Event Event `json:"event"`
}

// Corrector creates service-owned correction identity and timestamps.
type Corrector struct {
	journal *Journal
	now     func() time.Time
	entropy io.Reader
}

// NewCorrector constructs a corrector without performing IO.
func NewCorrector(journal *Journal) *Corrector {
	return &Corrector{journal: journal, now: time.Now}
}

// Correct appends a self-contained correction to the current effective event.
func (corrector *Corrector) Correct(ctx context.Context, params CorrectParams) (CorrectResult, error) {
	if !corrector.journal.Loaded() {
		return CorrectResult{}, ErrJournalNotLoaded
	}

	target, ok := corrector.journal.Event(params.TargetID)
	if !ok {
		return CorrectResult{}, ErrCorrectionTargetNotFound
	}

	effective, ok := corrector.journal.Effective(eventMilestoneKey(target))
	if !ok || effective.ID != target.ID {
		return CorrectResult{}, ErrCorrectionTargetStale
	}

	recordedAt := corrector.now().UTC()

	id, err := NewEventID(recordedAt, corrector.entropy)
	if err != nil {
		return CorrectResult{}, err
	}

	note := target.Note
	if params.ReplacementNote != nil {
		note = *params.ReplacementNote
	}

	event := Event{
		SchemaVersion: SchemaVersion,
		ID:            id,
		Type:          EventTypeCorrection,
		RecordedAt:    recordedAt,
		Execution:     copyExecution(target.Execution),
		Repository:    target.Repository,
		Git:           target.Git,
		Note:          note,
		TargetID:      target.ID,
		Reason:        params.Reason,
	}
	if err := event.Validate(); err != nil {
		return CorrectResult{}, fmt.Errorf("%w: %w", ErrEventInvalid, err)
	}

	result, err := corrector.journal.Append(ctx, event)
	if err != nil {
		if errors.Is(err, ErrCorrectionTargetStale) {
			return CorrectResult{}, err
		}

		return CorrectResult{}, err
	}

	return CorrectResult{Event: result.Event}, nil
}

func copyExecution(execution *ExecutionContext) *ExecutionContext {
	if execution == nil {
		return nil
	}

	result := *execution

	return &result
}
