# Acta

Acta is a cross-session agent chronicle manager. It records concise evidence
about significant Git commits in a local, append-only journal that independent
agent sessions can query later.

Acta is not a transcript store, an automatic commit logger, or shared live
context. The working agent decides whether a completed commit is significant.

## Current workflow

Start the per-user service:

```sh
acta serve
```

After creating a significant commit in a Git worktree, record it:

```sh
acta record "Reusable protocol behavior across agent-enabled projects"
```

When the commit subject is sufficient:

```sh
acta record
```

An agent integration can attach known execution provenance:

```sh
acta record \
  --harness codex \
  --harness-version 1.2.3 \
  --model provider/model \
  "Resolved a non-obvious service boundary issue"
```

Acta derives the repository root, sanitized repository identity, branch, full
commit hash, subject, Git author and committer times, and worktree state. The
service assigns the event ID and recording time.

## Agent instruction

A workspace can include this abbreviated rule in `AGENTS.md`:

```md
After creating a significant commit, run `acta record` to add it to the shared activity journal.
A commit is significant when it completes a milestone, records a reusable decision, or resolves a non-obvious problem.
Do not record routine, mechanical, formatting, punctuation-only, or work-in-progress commits.
```

Do not record a commit that only changes a comma, spelling, formatting, generated
output, or routine metadata unless that change represents a distinct reusable
milestone.

## Local data

- Configuration: `$HOME/.config/acta/`
- Canonical journal: `$HOME/.acta/events/YYYY/MM.jsonl`
- Recovery operation: `$HOME/.acta/pending.json`
- Service socket: `$XDG_RUNTIME_DIR/acta/acta.sock`

The service is the sole journal writer. The CLI does not edit journal files
when the service is unavailable.

## Development

```sh
make verify
```

The implementation uses the Go standard library. Runtime logic lives in
`internal/acta/`; root `main.go` only assembles process dependencies and runs
the application.
