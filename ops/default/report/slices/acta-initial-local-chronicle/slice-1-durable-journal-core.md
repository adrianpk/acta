# Slice 1: Durable Journal Core

Status: reviewing
Delivery set: Acta Initial Local Chronicle
Plan: `ops/default/plan/acta-initial-local-chronicle.md`
Tracker: `ops/default/tracker/acta-initial-local-chronicle.md`
Branch: `feat/acta-journal-core`
PR: pending

## Purpose

Define the canonical Acta event model and deliver a durable monthly journal
that can recover a narrowly interrupted append before later service work uses
it.

## Delivered Behavior

- Defines validated `acta/v1` milestone and correction events with optional
  execution provenance and self-contained correction snapshots.
- Generates and validates ULIDs and normalizes credential-free repository
  identities.
- Stores compact JSONL events in UTC monthly segments with owner-only file and
  directory permissions.
- Rebuilds global event, milestone, and effective-event indexes from canonical
  segments.
- Returns an existing milestone for duplicate repository and commit keys.
- Recovers missing, complete, and matching partial appends from one durable
  pending operation. Divergent states block journal loading and remain intact.
- Serializes concurrent writers and preserves valid journal framing under race
  testing.

## Implementation Notes

The journal holds one in-process write lock. It writes and syncs
`pending.json` before the segment append. Recovery compares the current segment
length and terminal bytes with the stored original length and complete event
line. Derived indexes update only after the durable append succeeds.

The repository also gained its strict Go linting scaffold. The configured
quality gate includes `nlreturn`, `wsl_v5`, error checks, complexity checks,
formatters, and the standard repository Make targets.

## Contracts Added or Changed

- Canonical events use schema version `acta/v1`.
- Journal segments use `events/YYYY/MM.jsonl` under the configured files root.
- A milestone key is `(repository.identity, git.commit)`.
- Corrections must target the current effective event and repeat its execution,
  repository, and Git snapshots.
- The recoverable operation uses schema version `acta/pending/v1` and records a
  relative segment path, original byte length, and complete canonical event
  line.
- Files use mode `0600`; Acta-owned directories use mode `0700`.

## Files of Interest

- `internal/acta/event.go`
- `internal/acta/identity.go`
- `internal/acta/journal.go`
- `internal/acta/pending.go`
- `.golangci.yml`

## Validation

- `go fmt ./...` (pass)
- `make lint-strict` (pass, 0 issues)
- `go test ./internal/acta` (pass)
- `go test -race ./internal/acta` (pass)
- `go test ./...` (pass)

## Risks and Follow-ups

- Cross-process writer ownership is not part of this core layer. Slice 2
  provides the sole-writer service boundary.
- README, Diataxis documentation, changelog, and agent-skill examples remain in
  separate follow-up tickets.
