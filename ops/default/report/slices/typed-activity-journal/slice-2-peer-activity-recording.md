# Slice 2: Peer Activity Recording

Status: reviewing
Delivery set: Typed Activity Journal
Plan: `ops/default/plan/typed-activity-journal.md`
Tracker: `ops/default/tracker/typed-activity-journal.md`
Branch: `feat/typed-activity-recording`
PR: pending

## Purpose

Deliver explicit CLI recording workflows for every initial peer activity type
and their optional common context without making Git a requirement for
professional activity records.

## Delivered Behavior

- Dispatches explicit `commit`, `pull-request`, `worklog`, and `note` record
  commands; bare `acta record` and unsupported type names are invalid.
- Accepts explicit occurrence time, all-or-none custom scope, ordered typed
  references, and optional human, agent, or integration provenance.
- Keeps commit recording Git-derived while allowing typed references and
  provenance on the canonical commit event.
- Records pull-request lifecycle actions from supplied forge, repository,
  number, action, optional URL, and summary without contacting the forge.
- Derives pull-request repository scope deterministically and treats each
  lifecycle action as a separate natural identity.
- Records worklogs with explicit bounded minutes and records user-directed
  notes without requiring Git, scope, references, or type-specific properties.
- Defaults omitted non-Git occurrence time to the service-assigned recording
  time and never infers duration.
- Prints compact receipts with status, ID, type, and optional scope and
  type-specific identifier.

## Implementation Notes

Common flags are parsed and validated before any service request. Commit alone
invokes the local Git inspector. The other record types construct canonical
service parameters directly, so they work outside a repository and perform no
network requests.

Pull-request forge URLs are normalized to a credential-free base URL and the
repository scope identity is `<forge>/<owner>/<repository>`. Repeating one
pull-request action returns the existing effective event; worklogs and notes
remain independent declarations and append on every successful request.

## Contracts Added or Changed

- `acta record pull-request` requires forge, `owner/name`, positive number,
  supported action, and summary.
- `acta record worklog` requires an integer duration from 1 through 1,440 and
  a summary.
- `acta record note` requires only a summary.
- `--at`, custom scope fields, typed references, and provenance fields are
  common recording inputs, with commit and pull-request scope derived by their
  adapters.
- Agent provenance requires a harness; human and integration provenance reject
  agent-only fields.
- The canonical pull-request scope must match its properties.

## Files of Interest

- `internal/acta/cli.go`
- `internal/acta/activity.go`
- `internal/acta/record.go`
- `internal/acta/recording_cli_test.go`
- `internal/acta/app_test.go`

## Validation

- `make check` (pass)
- `make test` (pass)
- `make docs-check` (pass)
- `make verify` (pass, 0 lint issues)
- `go test -race ./...` (pass)
- `go test ./internal/acta -run "Test(Parse.*Record|AppRecordsEveryPeerType|ActivityValidatePullRequest|ActivityValidateTypeProperties|Recorder)"` (pass)
- `git diff --check` (pass)

No candidate binary or systemd unit was installed, enabled, started, stopped,
or restarted. The active runtime and live journal remain on the legacy schema.

## Risks and Follow-ups

- Slice 3 must deliver compact search and report projections, the complete
  typed retrieval surface, and user documentation for these commands.
- Legacy journals remain intentionally unreadable by the candidate runtime
  until Slice 4 supplies the isolated deterministic migration path.
