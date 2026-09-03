# Slice 3: Activity Retrieval and Time Reporting

Status: reviewing
Delivery set: Typed Activity Journal
Plan: `ops/default/plan/typed-activity-journal.md`
Tracker: `ops/default/tracker/typed-activity-journal.md`
Branch: `feat/typed-activity-retrieval`
PR: pending

## Purpose

Deliver compact typed discovery, complete event detail, append-only corrections,
and occurrence-based reports with explicit worklog totals, then document the
complete typed activity workflow.

## Delivered Behavior

- Returns compact search projections containing event ID, type, occurrence and
  recording times, summary, optional scope name, optional worklog duration, and
  optional typed identifier.
- Filters effective events by peer activity type, scope, commit, pull-request
  identity, occurrence time, recording time, worklog duration, provenance, and
  normalized text across canonical fields.
- Resolves event IDs to complete effective events by default and returns exact
  stored chain members through raw detail lookup.
- Appends type-preserving complete corrections while allowing summary and
  occurrence changes plus explicit duration changes for worklogs.
- Groups reports by occurrence date in the selected local timezone and by
  scope, with unscoped activity represented explicitly.
- Reports event counts and sums only stored worklog duration at scope, day, and
  top-level boundaries.
- Documents peer recording, typed retrieval, corrections, reports, storage,
  capture policy, privacy, and interpretation boundaries in the README and
  Diataxis documentation.

## Implementation Notes

Search and reports share one effective-event selection path. Full canonical
events remain internal while search and report responses use
`ActivityProjection`; only show and correction results expose complete event
data. Both retrieval operations retain the 1,000-event and 8 MiB response
bounds.

Corrections clone the effective target before applying supplied replacement
fields. The service assigns a new ID and recording time, links the old event
through `supersedes`, validates the resulting type-specific shape, and appends
without rewriting history.

Report interval conversion occurs in the requested IANA timezone and is then
applied to `occurred_at`. Duration totals use only
`worklog.properties.duration_minutes`; commits, pull requests, and notes
contribute zero minutes.

## Contracts Added or Changed

- Search and report JSON event arrays now contain compact projections instead
  of complete canonical events.
- Search accepts `--source` in addition to all previously implemented typed
  filters.
- Show requires `--id`; `--raw` selects the exact stored event instead of the
  effective chain member.
- Correction uses `--summary` and additionally accepts `--at` and worklog-only
  `--duration` replacement values.
- Report accepts `--type`, groups by local occurrence date and scope, and emits
  event and explicit worklog-minute totals.

## Files of Interest

- `internal/acta/query.go`
- `internal/acta/report.go`
- `internal/acta/show.go`
- `internal/acta/correction.go`
- `internal/acta/retrieval_cli.go`
- `internal/acta/report_cli.go`
- `README.md`
- `docs/`

## Validation

- `make check` (pass)
- `make test` (pass)
- `make docs-check` (pass)
- `make verify` (pass, 0 lint issues)
- `go test -race ./...` (pass)
- `go test ./internal/acta -run 'Test(JournalSearch|JournalReport|JournalShow|Corrector|ParseSearch|ParseReport|AppRetrieval)'` (pass)
- `git diff --check` (pass)

No candidate binary or systemd unit was installed, enabled, started, stopped,
or restarted. The active runtime and live journal remain on the legacy schema.

## Risks and Follow-ups

- Legacy journals remain intentionally unreadable by the candidate runtime
  until Slice 4 supplies the isolated deterministic migration path.
- Installation and live migration remain gated on all slices, exact-candidate
  staging CI, recoverable snapshots, rehearsal, and a confirmed writer freeze.
