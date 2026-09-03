# Typed Activity Journal Tracker

Status: approved
Approval: approved by the maintainer on 2026-09-03
Delivery set: Typed Activity Journal
Slice strategy: layered
Plan: `ops/default/plan/typed-activity-journal.md`
Specification: `ops/default/spec/stable/typed-activity-journal.md`
Model: `ops/default/spec/stable/typed-activity-journal-model.md`
Progress: 1 of 4 slices completed
Active slice: Slice 2: Peer Activity Recording
Active slice status: in_progress

## Execution gate

The maintainer approved specification, planning, implementation, installation,
and deterministic live migration on 2026-09-03. Slice 1 may start after this
plan and tracker are committed to `dev` and match `origin/dev`.

## Slice map

| Slice | Status | Branch | Expected PR | Report |
| --- | --- | --- | --- | --- |
| Slice 1: Canonical Typed Journal Core | delivered | `feat/typed-activity-core` | `feat: add canonical typed activity core` | `ops/default/report/slices/typed-activity-journal/slice-1-canonical-typed-journal-core.md` |
| Slice 2: Peer Activity Recording | in_progress | `feat/typed-activity-recording` | `feat: add peer activity recording` | `ops/default/report/slices/typed-activity-journal/slice-2-peer-activity-recording.md` |
| Slice 3: Activity Retrieval and Time Reporting | pending | `feat/typed-activity-retrieval` | `feat: add typed activity retrieval` | `ops/default/report/slices/typed-activity-journal/slice-3-activity-retrieval-and-time-reporting.md` |
| Slice 4: Deterministic Journal Migration | pending | `feat/typed-activity-migration` | `feat: add deterministic journal migration` | `ops/default/report/slices/typed-activity-journal/slice-4-deterministic-journal-migration.md` |

## Slice 1 tasks

| ID | Task | Status | Commit | Validation |
| --- | --- | --- | --- | --- |
| `T1.1` | Add documentation gate and strict canonical typed values and validation. | completed | `be160e6` | `make docs-check`; `make lint-strict`; focused activity tests |
| `T1.2` | Replace journal indexes, correction chains, append behavior, and ordering. | completed | `b18efdd` | full repository gates; race tests; focused journal tests |
| `T1.3` | Replace correction and verification invariants and reject legacy or mixed journals. | completed | `b18efdd` | full repository gates; race tests; focused correction and verifier tests |
| `T1.4` | Adapt service and explicit commit recording to the canonical envelope. | completed | `b18efdd` | full repository gates; focused recorder and Git inspector tests |

## Slice 2 tasks

| ID | Task | Status | Commit | Validation |
| --- | --- | --- | --- | --- |
| `T2.1` | Implement occurrence time, scope, references, and provenance inputs. | pending | pending | pending |
| `T2.2` | Implement pull-request recording and action-specific idempotency. | pending | pending | pending |
| `T2.3` | Implement Git-independent worklog and user-directed note recording. | pending | pending | pending |
| `T2.4` | Complete typed CLI dispatch, receipts, boundaries, and end-to-end tests. | pending | pending | pending |

## Slice 3 tasks

| ID | Task | Status | Commit | Validation |
| --- | --- | --- | --- | --- |
| `T3.1` | Implement compact projections and typed activity filters. | pending | pending | pending |
| `T3.2` | Implement event-ID detail and type-preserving append-only correction. | pending | pending | pending |
| `T3.3` | Implement occurrence-date reports and explicit worklog totals. | pending | pending | pending |
| `T3.4` | Update README and Diataxis guidance for the complete typed workflow. | pending | pending | pending |

## Slice 4 tasks

| ID | Task | Status | Commit | Validation |
| --- | --- | --- | --- | --- |
| `T4.1` | Implement strict migration-only legacy decoding and canonical mapping. | pending | pending | pending |
| `T4.2` | Implement deterministic offline migration and cutover manifest. | pending | pending | pending |
| `T4.3` | Add malformed, lossy, repeatability, permission, publication, and rollback tests. | pending | pending | pending |
| `T4.4` | Complete end-to-end migration, command docs, asset checks, and rehearsal gate. | pending | pending | pending |

## Validation log

- `make check` (pass)
- `make test` (pass)
- `make docs-check` (pass)
- `make verify` (pass, 0 lint issues)
- `go test -race ./...` (pass)
- `go test ./internal/acta -run 'Test(Activity|Journal|Corrector|Recorder|Verifier|GitInspector)'` (pass)
- `git diff --check` (pass)

## Deployment state

- Installed runtime: legacy schema, active.
- Live journal: legacy schema, verified before planning.
- Candidate installation: pending.
- Rehearsal migration: pending.
- Writer freeze: pending.
- Live migration: pending.
- Rollback snapshots: pending.

## Blockers

None recorded.

## History

- 2026-09-03: The maintainer approved the typed activity redesign from
  specification through installation and live journal migration.
- 2026-09-03: Created the approved four-slice layered delivery map and
  activated Slice 1 pending the planning commit and origin alignment.
- 2026-09-03: Started Slice 1 on `feat/typed-activity-core` from aligned
  `dev`.
- 2026-09-03: Completed Slice 1 implementation and validation; moved it to
  review with the live legacy runtime and journal unchanged.
- 2026-09-03: Opened pull request `#7` from `feat/typed-activity-core` to
  `dev`.
- 2026-09-03: Verified pull request `#7` merged into `dev`, delivered Slice 1,
  and activated Slice 2.
