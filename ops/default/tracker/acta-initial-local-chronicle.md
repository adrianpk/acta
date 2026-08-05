# Acta Initial Local Chronicle Tracker

Status: active
Approval: approved
Delivery set: Acta Initial Local Chronicle
Slice strategy: layered
Plan: `ops/default/plan/acta-initial-local-chronicle.md`
Specification: `ops/default/spec/drafts/acta-system-specification-draft.md`
Model: `ops/default/spec/drafts/acta-system-specification-draft-model.md`
Progress: 3 of 4 slices completed
Active slice: Slice 4, Integrity and Operational Completion
Active slice status: reviewing

## Execution Gate

The maintainer approved the complete delivery map. Slice 1 may start after this
approval state is committed to `dev`.

## Slice Map

| Slice | Status | Branch | Expected PR | Report |
|---|---|---|---|---|
| Slice 1: Durable Journal Core | delivered | `feat/acta-journal-core` | `feat: add durable monthly journal core` | `ops/default/report/slices/acta-initial-local-chronicle/slice-1-durable-journal-core.md` |
| Slice 2: Local Service and Recording | delivered | `feat/acta-record-service` | `feat: add local service recording workflow` | `ops/default/report/slices/acta-initial-local-chronicle/slice-2-local-service-and-recording.md` |
| Slice 3: Bounded Retrieval and Corrections | delivered | `feat/acta-retrieval` | `feat: add bounded chronicle retrieval` | `ops/default/report/slices/acta-initial-local-chronicle/slice-3-bounded-retrieval-and-corrections.md` |
| Slice 4: Integrity and Operational Completion | reviewing | `feat/acta-integrity` | `feat: add journal integrity verification` | `ops/default/report/slices/acta-initial-local-chronicle/slice-4-integrity-and-operational-completion.md` |

## Slice 1 Tasks

| ID | Task | Status | Commit | Validation |
|---|---|---|---|---|
| `T1.1` | Define self-contained event values, optional execution provenance, validation, ULID generation, repository identity normalization, and tests. | completed | `0104004` | Focused tests and race tests pass. |
| `T1.2` | Implement monthly segment selection, journal loading, global indexes, and idempotent append with tests. | completed | `4f531f1` | Focused tests and race tests pass. |
| `T1.3` | Implement pending-operation recovery, corruption blocking, and concurrent writer tests. | completed | `ade10d9` | Focused tests, race tests, and strict lint pass. |

## Slice 2 Tasks

| ID | Task | Status | Commit | Validation |
|---|---|---|---|---|
| `T2.1` | Implement bounded JSON envelopes, Unix socket transport, limits, errors, and tests. | completed | `4d9ff97` | Focused tests, race tests, and strict lint pass. |
| `T2.2` | Implement Git metadata and timestamps, remote sanitization, optional execution provenance, recording, idempotency, and failure tests. | completed | `bdfa6cf` | Focused tests, race tests, and strict lint pass. |
| `T2.3` | Add root `main.go`, `record` and `serve`, lifecycle behavior, and end-to-end tests. | completed | `b710718` | Full suite and service-restart integration tests pass. |
| `T2.4` | Add and validate per-user service and socket activation assets without changing the live service. | completed | `8258ef3` | Socket adoption tests and systemd unit verification pass. |

## Slice 3 Tasks

| ID | Task | Status | Commit | Validation |
|---|---|---|---|---|
| `T3.1` | Implement filtering, ordering, text matching, self-contained correction chains, execution and Git-time filters, and bounded results with tests. | completed | `311dd6a` | Focused search, effective-view, ordering, filter, and truncation tests pass. |
| `T3.2` | Implement `show`, `search`, and `correct` service and CLI operations with boundary tests. | completed | `dcd5dd5` | Service, CLI, not-found, stale-target, count-limit, byte-limit, and end-to-end tests pass. |
| `T3.3` | Implement interval reports, timezone handling, grouped and JSON output, and cross-month tests. | completed | `0e00531` | Local/IANA timezone, interval, grouping, cross-month, JSON, and truncation tests pass. |

## Slice 4 Tasks

| ID | Task | Status | Commit | Validation |
|---|---|---|---|---|
| `T4.1` | Implement active and full verification with exact corruption diagnostics. | completed | `7095d38` | UTF-8, framing, schema, field, placement, identity, correction, line, and byte-offset tests pass. |
| `T4.2` | Add `verify` operations and the daily per-user timer with monthly full-verification behavior. | completed | `e574068`, `0514785`, `cfe6751` | CLI, service, degraded startup, persistent monthly state, timer asset, and systemd unit tests pass. |
| `T4.3` | Complete end-to-end recovery, concurrency, IPC, self-contained correction, execution filtering, verification, and restart tests. | completed | `01994dc` | Integrated lifecycle acceptance tests and race tests pass. |

## Validation Log

- 2026-08-05: `go fmt ./...` completed successfully.
- 2026-08-05: `make lint-strict` passed with 0 issues.
- 2026-08-05: `go test ./internal/acta` passed.
- 2026-08-05: `go test -race ./internal/acta` passed.
- 2026-08-05: `go test ./...` passed.
- 2026-08-05: Slice 2 `make lint-strict` passed with 0 issues.
- 2026-08-05: Slice 2 `go test ./internal/acta` passed.
- 2026-08-05: Slice 2 `go test -race ./internal/acta` passed.
- 2026-08-05: Slice 2 `go test ./...` passed.
- 2026-08-05: `systemd-analyze verify assets/systemd/user/acta.service assets/systemd/user/acta.socket` passed.
- 2026-08-05: Slice 3 `make lint-strict` passed with 0 issues.
- 2026-08-05: Slice 3 `go test ./internal/acta` passed.
- 2026-08-05: Slice 3 `go test -race ./internal/acta` passed.
- 2026-08-05: Slice 3 `go test ./...` passed.
- 2026-08-05: Slice 3 `systemd-analyze verify assets/systemd/user/acta.service assets/systemd/user/acta.socket` passed.
- 2026-08-05: Slice 4 `make lint-strict` passed with 0 issues.
- 2026-08-05: Slice 4 `go test ./internal/acta` passed.
- 2026-08-05: Slice 4 `go test -race ./internal/acta` passed.
- 2026-08-05: Slice 4 `go test ./...` passed.
- 2026-08-05: Slice 4 `systemd-analyze verify assets/systemd/user/acta.service assets/systemd/user/acta.socket assets/systemd/user/acta-verify.service assets/systemd/user/acta-verify.timer` passed.

## Blockers

None recorded.

## History

- 2026-08-05: Proposed the four-slice delivery map and created the execution
  tracker on `dev`.
- 2026-08-05: Aligned proposed tasks with self-contained corrections, Git
  timestamps, and optional execution provenance.
- 2026-08-05: Maintainer approved the delivery map and activated Slice 1.
- 2026-08-05: Completed T1.1, T1.2, and T1.3 and moved Slice 1 to review.
- 2026-08-05: Opened Slice 1 pull request `#1` against `dev`.
- 2026-08-05: Verified pull request `#1` merged into `dev`, delivered Slice 1,
  and activated Slice 2.
- 2026-08-05: Completed T2.1 through T2.4 and moved Slice 2 to review.
- 2026-08-05: Opened Slice 2 pull request `#2` against `dev`.
- 2026-08-05: Verified pull request `#2` merged into `dev`, delivered Slice 2,
  and activated Slice 3.
- 2026-08-05: Completed T3.1, T3.2, and T3.3 and moved Slice 3 to review.
- 2026-08-05: Opened Slice 3 pull request `#3` against `dev`.
- 2026-08-05: Verified pull request `#3` merged into `dev`, delivered Slice 3,
  and activated Slice 4.
- 2026-08-05: Completed T4.1, T4.2, and T4.3 and moved Slice 4 to review.
- 2026-08-05: Opened Slice 4 pull request `#4` against `dev`.
