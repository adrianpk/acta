package acta

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

// RepositoryEvidence contains the repository and Git facts derived by a client.
type RepositoryEvidence struct {
	Repository RepositoryReference
	Git        GitReference
}

// GitInspector reads repository evidence without changing Git state.
type GitInspector struct{}

// NewGitInspector constructs a read-only Git metadata inspector.
func NewGitInspector() *GitInspector {
	return &GitInspector{}
}

// Inspect reads canonical repository and HEAD metadata from a worktree.
func (inspector *GitInspector) Inspect(ctx context.Context, directory string) (RepositoryEvidence, error) {
	root, err := inspectRepositoryRoot(ctx, directory)
	if err != nil {
		return RepositoryEvidence{}, err
	}

	remote, err := inspectOptionalGitValue(ctx, root, "config", "--get", "remote.origin.url")
	if err != nil {
		return RepositoryEvidence{}, err
	}

	repository := repositoryReferenceForGit(root, remote)

	commit, err := inspectCommit(ctx, root)
	if err != nil {
		return RepositoryEvidence{}, err
	}

	branch, err := inspectOptionalGitValue(ctx, root, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil {
		return RepositoryEvidence{}, err
	}

	status, err := runGit(ctx, root, "status", "--porcelain=v1", "--untracked-files=normal")
	if err != nil {
		return RepositoryEvidence{}, err
	}

	commit.Branch = branch
	commit.WorktreeClean = len(status) == 0

	if err := repository.Validate(); err != nil {
		return RepositoryEvidence{}, err
	}

	if err := commit.Validate(); err != nil {
		return RepositoryEvidence{}, err
	}

	return RepositoryEvidence{Repository: repository, Git: commit}, nil
}

func inspectRepositoryRoot(ctx context.Context, directory string) (string, error) {
	if directory == "" {
		return "", fmt.Errorf("inspect Git repository: directory is required")
	}

	output, err := runGit(ctx, directory, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("inspect Git repository root: %w", err)
	}

	root := string(output)

	root, err = filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve Git repository root: %w", err)
	}

	canonical, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("resolve canonical Git repository root: %w", err)
	}

	return filepath.Clean(canonical), nil
}

func inspectCommit(ctx context.Context, root string) (GitReference, error) {
	output, err := runGit(ctx, root, "show", "-s", "--format=%H%x00%s%x00%aI%x00%cI", "HEAD")
	if err != nil {
		return GitReference{}, fmt.Errorf("inspect Git HEAD: %w", err)
	}

	fields := bytes.Split(output, []byte{0})
	if len(fields) != 4 {
		return GitReference{}, fmt.Errorf("inspect Git HEAD: unexpected metadata format")
	}

	authoredAt, err := parseGitTime("author", string(fields[2]))
	if err != nil {
		return GitReference{}, err
	}

	committedAt, err := parseGitTime("committer", string(fields[3]))
	if err != nil {
		return GitReference{}, err
	}

	return GitReference{
		Commit:      string(fields[0]),
		Subject:     string(fields[1]),
		AuthoredAt:  authoredAt,
		CommittedAt: committedAt,
	}, nil
}

func parseGitTime(name, value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("inspect Git %s time: %w", name, err)
	}

	return parsed.UTC(), nil
}

func repositoryReferenceForGit(root, remote string) RepositoryReference {
	reference, err := NewRepositoryReference(root, remote)
	if err == nil {
		return reference
	}

	reference, localErr := NewRepositoryReference(root, "")
	if localErr != nil {
		return RepositoryReference{Name: filepath.Base(root), Root: root}
	}

	return reference
}

func inspectOptionalGitValue(ctx context.Context, root string, arguments ...string) (string, error) {
	output, err := runGitCommand(ctx, root, arguments...)
	if err == nil {
		return string(trimGitOutput(output)), nil
	}

	var exitError *exec.ExitError
	if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
		return "", nil
	}

	return "", fmt.Errorf("run Git %s: %w", arguments[0], err)
}

func runGit(ctx context.Context, root string, arguments ...string) ([]byte, error) {
	output, err := runGitCommand(ctx, root, arguments...)
	if err != nil {
		return nil, fmt.Errorf("run Git %s: %w", arguments[0], err)
	}

	return trimGitOutput(output), nil
}

func runGitCommand(ctx context.Context, root string, arguments ...string) ([]byte, error) {
	commandArguments := append([]string{"-C", root}, arguments...)
	command := exec.CommandContext(ctx, "git", commandArguments...)

	return command.Output()
}

func trimGitOutput(output []byte) []byte {
	return bytes.TrimSuffix(output, []byte("\n"))
}
