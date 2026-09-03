# Typed Activity Journal Plan

Status: approved
Approval: approved by the maintainer on 2026-09-03
Cutover validation: corrected by the maintainer on 2026-09-03 to use the exact
`dev` commit, the complete local repository suite, and deterministic two-run
rehearsal; staging CI is not part of the Acta workflow
Delivery set: Typed Activity Journal
Slice strategy: layered
Tracker: `ops/default/tracker/typed-activity-journal.md`
Specification: `ops/default/spec/stable/typed-activity-journal.md`
Model: `ops/default/spec/stable/typed-activity-journal-model.md`
Source ticket: `ops/default/ticket/closed/20260903094823-adopt-canonical-typed-activity-schema.md`

## Purpose

Replace the Git-milestone-shaped journal with one canonical typed activity
journal. Deliver peer `commit`, `pull_request`, `worklog`, and `note` events,
compact retrieval and timesheet-ready reporting, and deterministic migration of
the installed legacy journal.

This plan authorizes all recorded slices after this plan and its tracker are
committed to `dev` and match `origin/dev`.

## Slice strategy

The delivery uses layered slices because the canonical persisted model,
append-only correction invariants, typed recording adapters, retrieval
projections, and offline migration have distinct correctness boundaries. The
installed legacy runtime remains unchanged until every slice is integrated and
the complete candidate passes rehearsal.

## Architecture constraints

- Keep `main.go` as a thin process entrypoint and keep application behavior in
  `internal/acta/`.
- Use the Go standard library only.
- Keep one canonical normal-runtime decoder after Slice 1. Isolate legacy
  decoding in migration-owned code introduced by Slice 4.
- Keep the service as the sole canonical writer and preserve the bounded
  pending-append transaction.
- Keep monthly JSONL segments canonical and derived indexes disposable.
- Represent commit, pull request, worklog, and note as peer event variants.
- Keep type-specific properties separate from optional references and
  provenance.
- Never infer work duration or activity from other events.
- Keep search and report responses compact. Return complete records only from
  detail operations.
- Do not replace the installed binary, modify the live journal, or stop user
  units during slice implementation.
- Add the repository `make docs-check` gate in Slice 1 so every slice pull
  request can use the mandated validation command set.

## Ordered slices

### Slice 1: Canonical Typed Journal Core

Purpose: Replace the legacy event and persistence model with the strict peer
typed-event envelope while preserving append safety and correction semantics.

Branch: `feat/typed-activity-core`

Expected PR title: `feat: add canonical typed activity core`

Report: `ops/default/report/slices/typed-activity-journal/slice-1-canonical-typed-journal-core.md`

Tasks:

- `T1.1` Add the repository documentation gate and define strict canonical
  activity, scope, properties, references, provenance, and projection values
  with validation tests.
- `T1.2` Replace milestone-key indexing with event chains and type-specific
  natural keys. Preserve monthly append, pending recovery, deterministic
  ordering, and concurrency safety.
- `T1.3` Replace correction and verification invariants with canonical typed
  behavior. Reject legacy and mixed-shape records in normal runtime paths.
- `T1.4` Adapt service boundaries and existing commit recording to the
  canonical envelope with explicit `acta record commit` selection so the
  repository remains executable and fully tested.

Validation:

```sh
make check
make test
make docs-check
make verify
go test -race ./internal/acta
git diff --check
```

Prerequisites: Approved plan and tracker committed to aligned `dev`.

### Slice 2: Peer Activity Recording

Purpose: Deliver explicit recording workflows for every initial peer event
type and optional common context.

Branch: `feat/typed-activity-recording`

Expected PR title: `feat: add peer activity recording`

Report: `ops/default/report/slices/typed-activity-journal/slice-2-peer-activity-recording.md`

Tasks:

- `T2.1` Implement common occurrence time, custom scope, typed references, and
  provenance input contracts with boundary tests.
- `T2.2` Implement pull-request recording with action-specific natural
  identity and no forge network access.
- `T2.3` Implement worklog and user-directed note recording without Git,
  mandatory scope, or mandatory references.
- `T2.4` Complete CLI dispatch, receipts, idempotency, invalid-input behavior,
  and end-to-end service tests for all four explicit record types.

Validation:

```sh
make check
make test
make docs-check
make verify
go test -race ./internal/acta
git diff --check
```

Prerequisites: Slice 1 delivered and closed on `dev`.

### Slice 3: Activity Retrieval and Time Reporting

Purpose: Deliver compact typed search, complete detail lookup, append-only
corrections, and reports with explicit worklog totals.

Branch: `feat/typed-activity-retrieval`

Expected PR title: `feat: add typed activity retrieval`

Report: `ops/default/report/slices/typed-activity-journal/slice-3-activity-retrieval-and-time-reporting.md`

Tasks:

- `T3.1` Implement compact activity projections and typed filters across scope,
  natural identifiers, occurrence time, recording time, duration, provenance,
  and text.
- `T3.2` Implement event-ID detail lookup and type-preserving append-only
  correction with raw-event access where specified.
- `T3.3` Group reports by local occurrence date and scope, and compute event and
  explicit worklog-minute totals without inference.
- `T3.4` Update the README and Diataxis documentation for all typed record,
  retrieval, correction, and timesheet workflows.

Validation:

```sh
make check
make test
make docs-check
make verify
go test -race ./internal/acta
git diff --check
```

Prerequisites: Slice 2 delivered and closed on `dev`.

### Slice 4: Deterministic Journal Migration

Purpose: Deliver the isolated legacy decoder, deterministic offline migrator,
cutover manifest, and complete candidate validation required for live
installation.

Branch: `feat/typed-activity-migration`

Expected PR title: `feat: add deterministic journal migration`

Report: `ops/default/report/slices/typed-activity-journal/slice-4-deterministic-journal-migration.md`

Tasks:

- `T4.1` Implement strict legacy milestone and correction decoding in a
  migration-only boundary and map every preserved field to canonical events.
- `T4.2` Implement offline source verification, non-overwriting destination
  publication, canonical verification, deterministic output, and the
  machine-readable cutover manifest.
- `T4.3` Add malformed-source, lossy-mapping, repeatability, permission,
  partial-publication, and rollback-oriented failure tests.
- `T4.4` Complete end-to-end migration fixtures, command documentation,
  systemd asset verification, and the immutable-candidate rehearsal gate.

Validation:

```sh
make check
make test
make docs-check
make verify
go test -race ./internal/acta
systemd-analyze verify assets/systemd/user/acta.service assets/systemd/user/acta.socket assets/systemd/user/acta-verify.service assets/systemd/user/acta-verify.timer
git diff --check
```

Prerequisites: Slice 3 delivered and closed on `dev`.

## Delivery rules

- Execute slices in numeric order.
- Use the exact branch, pull-request title, report path, and task IDs recorded
  here and in the tracker.
- Use one Conventional Commit per task when practical.
- Target each slice pull request to `dev` and use the global slice pull-request
  body template without modification.
- Keep every slice compiling and passing its focused and repository gates.
- Do not install a slice build or expose it to the legacy live journal.
- After each maintainer-controlled merge, return to `dev`, close the slice
  report and tracker state, and activate the next recorded slice.
- After Slice 4, validate the exact integrated `dev` candidate with the
  complete local repository suite and two byte-identical migrations from one
  verified legacy journal copy.
- Do not create an additional integration or tuning slice without an approved
  plan and tracker update.

## Live cutover gate

Installation and live migration may start only after:

1. all four slices are delivered on `dev`;
2. the exact integrated `dev` candidate passes the complete local repository
   gate;
3. two independent rehearsal migrations from the same verified source copy
   produce byte-identical canonical journals and manifests;
4. the installed binary and live journal have recoverable snapshots;
5. the writer freeze is confirmed before the final migration reads the source.

The live cutover follows the specification exactly. The rollback snapshots are
retained after successful migration until the maintainer explicitly accepts
their removal.
