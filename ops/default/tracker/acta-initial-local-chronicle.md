# Acta Initial Local Chronicle Tracker

Status: proposed
Approval: pending
Delivery set: Acta Initial Local Chronicle
Slice strategy: layered
Plan: `ops/default/plan/acta-initial-local-chronicle.md`
Specification: `ops/default/spec/drafts/acta-system-specification-draft.md`
Model: `ops/default/spec/drafts/acta-system-specification-draft-model.md`
Progress: 0 of 4 slices completed
Active slice: Slice 1, Durable Journal Core
Active slice status: planned

## Execution Gate

Implementation is blocked until the maintainer approves the plan and tracker,
both artifacts have status `approved`, and that state is committed to `dev`.

## Slice Map

| Slice | Status | Branch | Expected PR | Report |
|---|---|---|---|---|
| Slice 1: Durable Journal Core | planned | `feat/acta-journal-core` | `feat: add durable monthly journal core` | `ops/default/report/slices/acta-initial-local-chronicle/slice-1-durable-journal-core.md` |
| Slice 2: Local Service and Recording | pending | `feat/acta-record-service` | `feat: add local service recording workflow` | `ops/default/report/slices/acta-initial-local-chronicle/slice-2-local-service-and-recording.md` |
| Slice 3: Bounded Retrieval and Corrections | pending | `feat/acta-retrieval` | `feat: add bounded chronicle retrieval` | `ops/default/report/slices/acta-initial-local-chronicle/slice-3-bounded-retrieval-and-corrections.md` |
| Slice 4: Integrity and Operational Completion | pending | `feat/acta-integrity` | `feat: add journal integrity verification` | `ops/default/report/slices/acta-initial-local-chronicle/slice-4-integrity-and-operational-completion.md` |

## Slice 1 Tasks

| ID | Task | Status | Commit | Validation |
|---|---|---|---|---|
| `T1.1` | Define configuration, event values, validation, ULID generation, and repository identity normalization with tests. | pending | pending | pending |
| `T1.2` | Implement monthly segment selection, journal loading, global indexes, and idempotent append with tests. | pending | pending | pending |
| `T1.3` | Implement pending-operation recovery, corruption blocking, and concurrent writer tests. | pending | pending | pending |

## Slice 2 Tasks

| ID | Task | Status | Commit | Validation |
|---|---|---|---|---|
| `T2.1` | Implement bounded JSON envelopes, Unix socket transport, limits, errors, and tests. | pending | pending | pending |
| `T2.2` | Implement Git metadata collection, remote sanitization, recording, idempotency, and failure tests. | pending | pending | pending |
| `T2.3` | Add root `main.go`, `record` and `serve`, lifecycle behavior, and end-to-end tests. | pending | pending | pending |
| `T2.4` | Add and validate per-user service and socket activation assets without changing the live service. | pending | pending | pending |

## Slice 3 Tasks

| ID | Task | Status | Commit | Validation |
|---|---|---|---|---|
| `T3.1` | Implement filtering, ordering, text matching, correction application, and bounded results with tests. | pending | pending | pending |
| `T3.2` | Implement `show`, `search`, and `correct` service and CLI operations with boundary tests. | pending | pending | pending |
| `T3.3` | Implement interval reports, timezone handling, grouped and JSON output, and cross-month tests. | pending | pending | pending |

## Slice 4 Tasks

| ID | Task | Status | Commit | Validation |
|---|---|---|---|---|
| `T4.1` | Implement active and full verification with exact corruption diagnostics. | pending | pending | pending |
| `T4.2` | Add `verify` operations and the daily per-user timer with monthly full-verification behavior. | pending | pending | pending |
| `T4.3` | Complete end-to-end recovery, concurrency, IPC, correction, verification, and restart tests. | pending | pending | pending |

## Validation Log

No implementation validation has run. Planning validation is recorded in the
plan commit handoff.

## Blockers

- Maintainer approval of the proposed delivery map is pending.

## History

- 2026-08-05: Proposed the four-slice delivery map and created the execution
  tracker on `dev`.
