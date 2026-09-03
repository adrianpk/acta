# Typed Activity Journal Tracker

Status: completed
Approval: approved by the maintainer on 2026-09-03
Cutover validation: exact `dev` commit, complete local repository suite, and
deterministic two-run rehearsal; staging CI is not part of the Acta workflow
Delivery set: Typed Activity Journal
Slice strategy: layered
Plan: `ops/default/plan/typed-activity-journal.md`
Specification: `ops/default/spec/stable/typed-activity-journal.md`
Model: `ops/default/spec/stable/typed-activity-journal-model.md`
Progress: 4 of 4 slices completed
Active slice: None
Active slice status: none

## Execution gate

The maintainer approved specification, planning, implementation, installation,
and deterministic live migration on 2026-09-03. Slice 1 may start after this
plan and tracker are committed to `dev` and match `origin/dev`.

## Slice map

| Slice | Status | Branch | Expected PR | Report |
| --- | --- | --- | --- | --- |
| Slice 1: Canonical Typed Journal Core | delivered | `feat/typed-activity-core` | `feat: add canonical typed activity core` | `ops/default/report/slices/typed-activity-journal/slice-1-canonical-typed-journal-core.md` |
| Slice 2: Peer Activity Recording | delivered | `feat/typed-activity-recording` | `feat: add peer activity recording` | `ops/default/report/slices/typed-activity-journal/slice-2-peer-activity-recording.md` |
| Slice 3: Activity Retrieval and Time Reporting | delivered | `feat/typed-activity-retrieval` | `feat: add typed activity retrieval` | `ops/default/report/slices/typed-activity-journal/slice-3-activity-retrieval-and-time-reporting.md` |
| Slice 4: Deterministic Journal Migration | delivered | `feat/typed-activity-migration` | `feat: add deterministic journal migration` | `ops/default/report/slices/typed-activity-journal/slice-4-deterministic-journal-migration.md` |

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
| `T2.1` | Implement occurrence time, scope, references, and provenance inputs. | completed | `f10633a` | full repository gates; common-context parser tests |
| `T2.2` | Implement pull-request recording and action-specific idempotency. | completed | `f10633a` | peer-type end-to-end and natural-identity tests |
| `T2.3` | Implement Git-independent worklog and user-directed note recording. | completed | `f10633a` | peer-type end-to-end and occurrence-time tests |
| `T2.4` | Complete typed CLI dispatch, receipts, boundaries, and end-to-end tests. | completed | `f10633a` | full repository gates; focused record CLI tests |

## Slice 3 tasks

| ID | Task | Status | Commit | Validation |
| --- | --- | --- | --- | --- |
| `T3.1` | Implement compact projections and typed activity filters. | completed | `88a21c3` | full repository gates; focused search and parsing tests |
| `T3.2` | Implement event-ID detail and type-preserving append-only correction. | completed | `88a21c3` | full repository gates; focused show and corrector tests |
| `T3.3` | Implement occurrence-date reports and explicit worklog totals. | completed | `88a21c3` | full repository gates; focused report tests |
| `T3.4` | Update README and Diataxis guidance for the complete typed workflow. | completed | `5703c1d` | `make docs-check`; `git diff --check` |

## Slice 4 tasks

| ID | Task | Status | Commit | Validation |
| --- | --- | --- | --- | --- |
| `T4.1` | Implement strict migration-only legacy decoding and canonical mapping. | completed | `952f191` | full repository gates; focused migration tests |
| `T4.2` | Implement deterministic offline migration and cutover manifest. | completed | `952f191` | repeatability, manifest, and publication tests |
| `T4.3` | Add malformed, lossy, repeatability, permission, publication, and rollback tests. | completed | `952f191` | focused migration tests; race tests |
| `T4.4` | Complete end-to-end migration, command docs, asset checks, and rehearsal gate. | completed | `952f191`, `3a35f3e` | CLI tests; `make docs-check`; systemd asset verification |

## Validation log

- `make check` (pass)
- `make test` (pass)
- `make docs-check` (pass)
- `make verify` (pass, 0 lint issues)
- `go test -race ./...` (pass)
- `go test ./internal/acta -run 'Test(Activity|Journal|Corrector|Recorder|Verifier|GitInspector)'` (pass)
- `git diff --check` (pass)
- Slice 2: `make check` (pass)
- Slice 2: `make test` (pass)
- Slice 2: `make docs-check` (pass)
- Slice 2: `make verify` (pass, 0 lint issues)
- Slice 2: `go test -race ./...` (pass)
- Slice 2: `go test ./internal/acta -run "Test(Parse.*Record|AppRecordsEveryPeerType|ActivityValidatePullRequest|ActivityValidateTypeProperties|Recorder)"` (pass)
- Slice 2: `git diff --check` (pass)
- Slice 3: `make check` (pass)
- Slice 3: `make test` (pass)
- Slice 3: `make docs-check` (pass)
- Slice 3: `make verify` (pass, 0 lint issues)
- Slice 3: `go test -race ./...` (pass)
- Slice 3: `go test ./internal/acta -run 'Test(JournalSearch|JournalReport|JournalShow|Corrector|ParseSearch|ParseReport|AppRetrieval)'` (pass)
- Slice 3: `git diff --check` (pass)
- Slice 4: `make check` (pass)
- Slice 4: `make test` (pass)
- Slice 4: `make docs-check` (pass)
- Slice 4: `make verify` (pass, 0 lint issues)
- Slice 4: `go test -race ./...` (pass)
- Slice 4: `go test -count=1 ./internal/acta -run 'Test(Migrator|ParseMigration|AppMigrates|Legacy)'` (pass)
- Slice 4: `systemd-analyze verify assets/systemd/user/acta.service assets/systemd/user/acta.socket assets/systemd/user/acta-verify.service assets/systemd/user/acta-verify.timer` (pass)
- Slice 4: `git diff --check` (pass)
- Exact candidate `b2ca90b`: `make check` (pass)
- Exact candidate `b2ca90b`: `make test` (pass)
- Exact candidate `b2ca90b`: `make docs-check` (pass)
- Exact candidate `b2ca90b`: `make verify` (pass, 0 lint issues)
- Exact candidate `b2ca90b`: `go test -race ./...` (pass)
- Exact candidate `b2ca90b`: `go test -count=1 ./internal/acta -run 'Test(Migrator|ParseMigration|AppMigrates|Legacy)'` (pass)
- Exact candidate `b2ca90b`: `systemd-analyze verify assets/systemd/user/acta.service assets/systemd/user/acta.socket assets/systemd/user/acta-verify.service assets/systemd/user/acta-verify.timer` (pass)
- Exact candidate `b2ca90b`: `git diff --check` (pass)
- Rehearsal: two independent 2-segment, 2,372-event migrations produced
  byte-identical journal trees and manifests (pass)
- Frozen cutover: final migration manifest matched the rehearsal manifest
  byte for byte (pass)
- Post-install: `acta verify --full` verified 2 segments and 2,374 events after
  the write and correction smoke checks (pass)
- Post-install: note recording, commit idempotency, search, effective and raw
  show, occurrence-date report, and append-only correction smoke checks (pass)

## Deployment state

- Installed runtime: canonical schema from `b2ca90b`, active.
- Live journal: canonical schema, full verification passed after smoke checks.
- Candidate installation: completed on 2026-09-03; installed binary SHA-256 is
  `53e13c75cf6820d848f58bf08726bb448ef48f64c6e9ce2d86dc2e3eaaf3525c`.
- Rehearsal migration: completed twice from one verified 2,372-event source
  copy; journal trees and manifests were byte-identical.
- Writer freeze: completed with socket, service, timer, and verification unit
  inactive and no Acta process present before snapshots and final migration.
- Live migration: completed from 2 legacy segments and 2,372 stored events to
  2 canonical segments and 2,372 stored events, with 2,371 effective events
  before post-install smoke writes.
- Migration manifest: retained at
  `$HOME/.acta/migration-manifest-20260903T154026CEST.json`, SHA-256
  `9006e5c75af981f708051e188044a6b88fddc75ec417c131f8d4f482702bb157`.
- Rollback snapshots: retained at
  `$HOME/.local/bin/acta.legacy-20260903T154026CEST`,
  `$HOME/.acta/events.legacy-snapshot-20260903T154026CEST`, and
  `$HOME/.acta/events.legacy-live-20260903T154026CEST` pending explicit
  maintainer acceptance.

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
- 2026-09-03: Completed Slice 2 implementation and validation; moved it to
  review with the installed runtime and live legacy journal unchanged.
- 2026-09-03: Opened pull request `#8` from
  `feat/typed-activity-recording` to `dev`.
- 2026-09-03: Verified pull request `#8` merged into `dev`, delivered Slice 2,
  and activated Slice 3.
- 2026-09-03: Completed Slice 3 implementation, documentation, and validation;
  moved it to review with the installed runtime and live legacy journal
  unchanged.
- 2026-09-03: Opened pull request `#9` from
  `feat/typed-activity-retrieval` to `dev`.
- 2026-09-03: Verified pull request `#9` merged into `dev`, delivered Slice 3,
  and activated Slice 4.
- 2026-09-03: Completed Slice 4 implementation, documentation, and validation;
  moved it to review with installation, rehearsal, writer freeze, and live
  migration still pending.
- 2026-09-03: Opened pull request `#10` from
  `feat/typed-activity-migration` to `dev`.
- 2026-09-03: Verified pull request `#10` merged into `dev`, delivered Slice 4,
  and completed the four-slice delivery set.
- 2026-09-03: The maintainer corrected the inherited cutover gate: Acta uses
  exact-commit local validation and deterministic rehearsal, not staging CI.
- 2026-09-03: Validated exact candidate `b2ca90b`, completed two byte-identical
  rehearsal migrations, and confirmed the frozen live source matched the
  rehearsal source.
- 2026-09-03: Snapshotted the legacy binary and journal, migrated and installed
  the canonical runtime and journal, then restored socket activation and
  scheduled verification.
- 2026-09-03: Completed full verification and record, idempotency, search,
  show, report, and append-only correction smoke checks; retained all rollback
  evidence pending explicit acceptance.
