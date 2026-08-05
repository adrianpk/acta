package acta

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// ErrEventNotFound identifies a retrieval with no effective event.
var ErrEventNotFound = errors.New("event was not found")

// ShowParams selects one effective event by commit.
type ShowParams struct {
	Commit     string `json:"commit"`
	Repository string `json:"repository,omitempty"`
}

// ShowResult contains one effective event.
type ShowResult struct {
	Event Event `json:"event"`
}

// Show returns one unambiguous effective event by commit.
func (journal *Journal) Show(ctx context.Context, params ShowParams) (ShowResult, error) {
	if strings.TrimSpace(params.Commit) == "" {
		return ShowResult{}, fmt.Errorf("%w: commit is required", ErrQueryInvalid)
	}

	result, err := journal.Search(ctx, Query{
		Repository: params.Repository,
		Commit:     params.Commit,
		Limit:      2,
	})
	if err != nil {
		return ShowResult{}, err
	}

	switch len(result.Events) {
	case 0:
		return ShowResult{}, fmt.Errorf("%w: commit %q", ErrEventNotFound, params.Commit)
	case 1:
		return ShowResult{Event: result.Events[0]}, nil
	default:
		return ShowResult{}, fmt.Errorf("%w: commit %q matches several repositories", ErrCommitAmbiguous, params.Commit)
	}
}
