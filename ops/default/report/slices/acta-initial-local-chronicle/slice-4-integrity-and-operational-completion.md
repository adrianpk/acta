# Slice 4: Integrity and Operational Completion

Status: reviewing
Delivery set: Acta Initial Local Chronicle
Plan: `ops/default/plan/acta-initial-local-chronicle.md`
Tracker: `ops/default/tracker/acta-initial-local-chronicle.md`
Branch: `feat/acta-integrity`
PR: pending

## Purpose

Complete read-only integrity verification, safe scheduled operation, and the
end-to-end acceptance contract for the initial local chronicle.

## Delivered Behavior

- Verifies the current UTC monthly segment without rescanning historical
  segments during normal daily operation.
- Verifies every canonical segment and rebuilds global invariants in full mode.
- Checks UTF-8, LF framing, JSON, schema and field rules, segment placement,
  event IDs, milestone keys, Git timestamps, and correction chains.
- Reports corruption with an exact segment, line, byte offset, and reason.
- Keeps verification read-only for canonical journal and pending files.
- Provides `acta verify`, `acta verify --full`, JSON output, and stable service
  errors.
- Starts the service in a read-only degraded state when journal load finds
  corruption, so full verification remains available while normal operations
  stay blocked.
- Schedules daily active verification at 03:00 local time with persistent
  systemd scheduling.
- Runs one full scheduled verification per local calendar month, including
  after a missed first-day run.

## Implementation Notes

Active verification holds the journal read lock, reads only the current segment,
and uses the loaded historical index to validate cross-segment correction
targets. Full verification reads all segment paths in canonical order and
rebuilds independent indexes.

`IntegrityError` unwraps to `ErrJournalCorrupt` and carries structured location
fields. JSON syntax offsets and invalid UTF-8 bytes identify their exact byte.
Field and global invariant failures identify the start byte of the affected
line.

Scheduled verification records only the last successful full month in
`~/.acta/verification.json`. This disposable state uses mode `0600` and an
atomic synced replacement. It prevents repeated full scans within one month
and ensures a delayed persistent timer performs the missed monthly scan.

## Contracts Added or Changed

- The local protocol adds the `verify` operation with active and full scopes.
- `acta verify [--full] [--json]` is the public verification interface.
- `acta verify --scheduled` is the systemd service entrypoint.
- Successful results contain scope, segment count, and event count.
- Corruption returns `journal_corrupt` with exact diagnostic text.
- A degraded service permits full verification and rejects journal-dependent
  record, retrieval, report, and correction operations.
- `acta-verify.timer` runs at 03:00 local time with `Persistent=true`.

## Files of Interest

- `internal/acta/verify.go`
- `internal/acta/verify_cli.go`
- `internal/acta/schedule.go`
- `internal/acta/service.go`
- `internal/acta/acceptance_test.go`
- `assets/systemd/user/acta-verify.service`
- `assets/systemd/user/acta-verify.timer`

## Validation

- `go fmt ./...` (pass)
- `make lint-strict` (pass, 0 issues)
- `go test ./internal/acta` (pass)
- `go test -race ./internal/acta` (pass)
- `go test ./...` (pass)
- `systemd-analyze verify assets/systemd/user/acta.service assets/systemd/user/acta.socket assets/systemd/user/acta-verify.service assets/systemd/user/acta-verify.timer` (pass)

Acceptance coverage includes interrupted append recovery, concurrent writes,
request and response limits, cross-segment corrections, execution provenance
filters, full verification, service restart, and index rebuild behavior.

No live user service or timer was installed, enabled, started, stopped, or
restarted.

## Risks and Follow-ups

- Installation and activation of the binary and user units still require an
  explicit user action.
- Historical checksums and automatic general repair remain outside the initial
  delivery scope.
