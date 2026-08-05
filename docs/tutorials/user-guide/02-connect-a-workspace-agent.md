# Connect a Workspace Agent

Add an Acta capture rule to one Git workspace. The rule tells the agent when to
record a completed commit and what evidence it may supply.

Add this selective policy to the workspace `AGENTS.md`:

```md
After creating a significant commit, immediately run `acta record [flags] [note]` while that commit is still `HEAD`.
Record completed milestones, reusable decisions, accepted specifications, and non-obvious fixes.
Use the note to state why the commit matters when the subject does not make that clear.
Supply `--harness`, `--harness-version`, and `--model` only when those values are known reliably.
Do not record routine, generated, formatting-only, punctuation-only, or work-in-progress commits.
Never include credentials, tokens, personal data, or hidden reasoning in the note.
```

Acta derives the repository, branch, commit, subject, Git timestamps, and
worktree state. The agent should not duplicate those values in the note.

Commit the workspace rule so later agent sessions receive the same capture
policy.

For a broader policy and guidance about choosing granularity, see [Choose an
agent capture policy](../../how-to/choose-an-agent-capture-policy.md).

Continue with [Record and retrieve a
milestone](03-record-and-retrieve-a-milestone.md).
