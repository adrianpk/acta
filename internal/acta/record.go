package acta

import (
	"context"
	"io"
	"time"
)

// RecordParams contains client-derived evidence and optional provenance.
type RecordParams struct {
	Repository RepositoryReference `json:"repository"`
	Git        GitReference        `json:"git"`
	Note       string              `json:"note,omitempty"`
	Execution  *ExecutionContext   `json:"execution,omitempty"`
}

// RecordResult describes one newly recorded or existing milestone.
type RecordResult struct {
	ID         EventID `json:"id"`
	Recorded   bool    `json:"recorded"`
	Repository string  `json:"repository"`
	Commit     string  `json:"commit"`
}

// Recorder creates service-owned event identity and publication timestamps.
type Recorder struct {
	journal *Journal
	now     func() time.Time
	entropy io.Reader
}

// NewRecorder constructs a recorder without performing IO.
func NewRecorder(journal *Journal) *Recorder {
	return &Recorder{journal: journal, now: time.Now}
}

// Record validates client evidence and appends one idempotent milestone.
func (recorder *Recorder) Record(ctx context.Context, params RecordParams) (RecordResult, error) {
	recordedAt := recorder.now().UTC()

	id, err := NewEventID(recordedAt, recorder.entropy)
	if err != nil {
		return RecordResult{}, err
	}

	event := Event{
		SchemaVersion: SchemaVersion,
		ID:            id,
		Type:          EventTypeMilestone,
		RecordedAt:    recordedAt,
		Execution:     params.Execution,
		Repository:    params.Repository,
		Git:           params.Git,
		Note:          params.Note,
	}

	result, err := recorder.journal.Append(ctx, event)
	if err != nil {
		return RecordResult{}, err
	}

	return RecordResult{
		ID:         result.Event.ID,
		Recorded:   result.Recorded,
		Repository: result.Event.Repository.Name,
		Commit:     result.Event.Git.Commit,
	}, nil
}
