# Slice 1: Canonical Typed Journal Core

Status: reviewing
Delivery set: Typed Activity Journal
Plan: `ops/default/plan/typed-activity-journal.md`
Tracker: `ops/default/tracker/typed-activity-journal.md`
Branch: `feat/typed-activity-core`
PR: pending

## Purpose

Replace the milestone-centered journal foundation with one strict canonical
activity envelope in which commit, pull request, worklog, and note are peer
types, while preserving the existing durable local service boundary.

## Delivered Behavior

- Defines strict canonical values and validation for all four peer activity
  types, scopes, references, provenance, properties, and corrections.
- Persists only canonical typed activities and rejects unknown, legacy, or
  mixed event fields during normal journal loading and pending recovery.
- Indexes raw events, correction roots, effective events, natural identities,
  and append order without treating commits as the parent record type.
- Keeps commit and pull-request natural identities idempotent while allowing
  independently recorded worklogs and notes.
- Enforces event-ID uniqueness, ULID timestamp agreement with `recorded_at`,
  monthly segment placement, and type-preserving correction chains.
- Uses the first Git root commit as the stable identity for repositories that
  have no usable remote, without persisting a checkout path.
- Requires explicit `acta record commit`; bare `acta record` is invalid.
- Retains the owner-only Unix service, pending-write recovery, bounded local
  protocol, and read-only verification behavior.

## Implementation Notes

The disposable journal index maps every event ID to its correction root and
maps each root to its current effective event. Commit and pull-request natural
keys point to correction roots, so a repeated record request returns the
current effective event without appending.

Normal decoding uses a single canonical activity decoder with unknown-field
rejection. Legacy decoding is deliberately absent from the runtime path and
remains reserved for the migration slice.

The retrieval and report internals were adapted to consume the canonical
envelope so the service remains coherent. Their final compact projections,
complete filters, and user-facing documentation remain Slice 3 work.

## Contracts Added or Changed

- The persisted event has no `schema_version` and uses the root fields defined
  by the stable typed activity model.
- `commit`, `pull_request`, `worklog`, and `note` are the only activity types.
- Corrections are complete same-type replacements linked by `supersedes` and
  `correction_reason`; correction is not an activity type.
- `occurred_at` and `recorded_at` are distinct UTC timestamps, and the ULID
  timestamp equals `recorded_at` at millisecond precision.
- A local repository identity uses `git-root:<root-object-id>`.
- Commit recording is explicitly selected with `acta record commit`.
- Markdown validation is available through `make docs-check`.

## Files of Interest

- `internal/acta/activity.go`
- `internal/acta/journal.go`
- `internal/acta/correction.go`
- `internal/acta/record.go`
- `internal/acta/verify.go`
- `internal/acta/git.go`
- `internal/acta/cli.go`
- `Makefile`

## Validation

- `make check` (pass)
- `make test` (pass)
- `make docs-check` (pass)
- `make verify` (pass, 0 lint issues)
- `go test -race ./...` (pass)
- `go test ./internal/acta -run 'Test(Activity|Journal|Corrector|Recorder|Verifier|GitInspector)'` (pass)
- `git diff --check` (pass)

No candidate binary or systemd unit was installed, enabled, started, stopped,
or restarted. The active runtime and live journal remain on the legacy schema.

## Risks and Follow-ups

- Slice 2 must add the complete common recording flags and pull-request,
  worklog, and note CLI commands.
- Slice 3 must replace full search and report records with compact projections
  and complete the final retrieval interface and documentation.
- Legacy journals remain intentionally unreadable by the candidate runtime
  until Slice 4 supplies the isolated deterministic migration path.
