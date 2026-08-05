# Slice 3: Bounded Retrieval and Corrections

Status: reviewing
Delivery set: Acta Initial Local Chronicle
Plan: `ops/default/plan/acta-initial-local-chronicle.md`
Tracker: `ops/default/tracker/acta-initial-local-chronicle.md`
Branch: `feat/acta-retrieval`
PR: `#3`

## Purpose

Deliver explicit bounded retrieval, append-only corrections, and calendar-based
reports through the existing local service.

## Delivered Behavior

- Returns one effective event per milestone correction chain by default.
- Filters by event type, repository, commit, execution harness, model, text,
  recording time, Git author time, and Git committer time.
- Provides deterministic ascending or descending order with a safe default
  limit, a maximum of 1,000 events, and an 8 MiB encoded response limit.
- Resolves one effective event by a full or unambiguous abbreviated commit.
- Appends a self-contained correction that copies the target's repository, Git,
  execution, and effective note snapshot.
- Rejects missing and stale correction targets without changing journal history.
- Produces reports for an explicit local-calendar interval in the invoking
  process timezone or a supplied IANA timezone.
- Groups reports by local date and repository and supports concise human output
  and machine-readable JSON.

## Implementation Notes

Search reads the disposable global journal indexes and excludes superseded
events. Commit prefixes must use lowercase hexadecimal characters. A prefix
that matches different commits is ambiguous. `show` also rejects a commit that
exists in several repositories unless the query identifies one repository.

Correction identity and `recorded_at` are service-owned. The optional
replacement note distinguishes an omitted value, which copies the effective
note, from an explicit empty value, which clears it. Concurrent correction
attempts are validated again while the journal write lock is held.

Report `from` dates are inclusive and `to` dates are exclusive. Both represent
local midnight in the effective timezone. Stored timestamps remain UTC. Report
group order and event order are deterministic.

## Contracts Added or Changed

- The local protocol adds `search`, `show`, `correct`, and `report` operations.
- Search results contain the normalized query, effective events, returned
  count, and truncation state.
- `acta show --commit <commit> [--repository <repository>] [--json]` resolves
  one effective event.
- `acta search` exposes repository, commit, execution, event type, text, time,
  order, count, and JSON controls.
- `acta correct --target <event-id> --reason <reason> [--note <note>] [--json]`
  appends one correction.
- `acta report --from <date> --to <date> [--timezone <name>] [--limit <count>]
  [--json]` returns grouped effective events.
- Retrieval failures use stable `invalid_request`, `not_found`,
  `invalid_event`, `journal_corrupt`, and `journal_unavailable` categories.

## Files of Interest

- `internal/acta/query.go`
- `internal/acta/show.go`
- `internal/acta/correction.go`
- `internal/acta/report.go`
- `internal/acta/retrieval_cli.go`
- `internal/acta/report_cli.go`
- `internal/acta/service.go`

## Validation

- `go fmt ./...` (pass)
- `make lint-strict` (pass, 0 issues)
- `go test ./internal/acta` (pass)
- `go test -race ./internal/acta` (pass)
- `go test ./...` (pass)
- `systemd-analyze verify assets/systemd/user/acta.service assets/systemd/user/acta.socket` (pass)

No live user service was installed, enabled, started, stopped, or restarted.

## Risks and Follow-ups

- Integrity verification and its scheduled user timer remain in Slice 4.
- Raw export remains outside the initial delivery scope. Search and report only
  expose effective bounded views.
