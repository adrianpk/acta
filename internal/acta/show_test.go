package acta

import (
	"context"
	"errors"
	"testing"
)

func TestJournalShowReturnsEffectiveCorrection(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)

	result, err := fixture.journal.Show(context.Background(), ShowParams{
		Commit: fixture.milestoneA.Git.Commit[:12],
	})
	if err != nil {
		t.Fatalf("Show() error = %v", err)
	}

	if result.Event.ID != fixture.correction.ID {
		t.Fatalf("Show() event = %#v", result.Event)
	}
}

func TestJournalShowRejectsMissingAndAmbiguousCommit(t *testing.T) {
	t.Parallel()

	fixture := newQueryFixture(t)
	duplicateCommit := distinctMilestoneEvent(t, 4)
	duplicateCommit.Git.Commit = fixture.milestoneA.Git.Commit
	duplicateCommit.Repository = RepositoryReference{
		Name:     "duplicate",
		Root:     "/workspace/duplicate",
		Identity: "https://example.test/team/duplicate.git",
		Remote:   "https://example.test/team/duplicate.git",
	}
	appendJournalEvent(t, fixture.journal, duplicateCommit)

	tests := []struct {
		name   string
		params ShowParams
		want   error
	}{
		{
			name:   "missing",
			params: ShowParams{Commit: "bbbbbbbbbbbb"},
			want:   ErrEventNotFound,
		},
		{
			name:   "ambiguous repository",
			params: ShowParams{Commit: fixture.milestoneA.Git.Commit},
			want:   ErrCommitAmbiguous,
		},
		{
			name:   "missing commit",
			params: ShowParams{},
			want:   ErrQueryInvalid,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if _, err := fixture.journal.Show(context.Background(), test.params); !errors.Is(err, test.want) {
				t.Fatalf("Show() error = %v, want %v", err, test.want)
			}
		})
	}

	result, err := fixture.journal.Show(context.Background(), ShowParams{
		Commit:     fixture.milestoneA.Git.Commit,
		Repository: fixture.milestoneA.Repository.Name,
	})
	if err != nil || result.Event.ID != fixture.correction.ID {
		t.Fatalf("qualified Show() = %#v, %v", result, err)
	}
}
