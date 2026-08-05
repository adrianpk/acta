# Acta Initial Local Chronicle Plan

Status: approved
Approval: approved
Delivery set: Acta Initial Local Chronicle
Tracker: `ops/default/tracker/acta-initial-local-chronicle.md`
Specification: `ops/default/spec/drafts/acta-system-specification-draft.md`
Model: `ops/default/spec/drafts/acta-system-specification-draft-model.md`

## Purpose

Implement the initial Acta product as a small Go application. The delivery must
record significant Git commits through a local per-user service, preserve them
in recoverable monthly JSONL segments, provide bounded retrieval, and verify
journal integrity.

This plan authorizes the recorded slices after its approved state and the
matching tracker state are committed to `dev`.

## Slice Strategy

Slice strategy: layered

Reason: The durable journal transaction, recovery rules, and global invariants
must be stable before the CLI and service expose them. Each later slice adds a
complete adapter or behavior layer while keeping one small `internal/acta`
package and a passing repository.

## Architecture Constraints

- Keep `main.go` at the repository root.
- Limit `main.go` to `main`, dependency assembly, process IO, and exit handling.
- Keep application logic in `internal/acta/`.
- Use one Go file for each important concern. Do not create packages for small
  internal details.
- Use the Go standard library only unless a later approved plan change permits
  an external dependency.
- Keep constructors free of failable IO. Use explicit run, setup, verify, and
  shutdown paths.
- Keep the service as the sole canonical journal writer.
- Keep monthly JSONL segments canonical and all indexes disposable.
- Do not start, enable, or modify live user services during implementation
  validation without separate authorization.

The expected concern map is:

```text
main.go
internal/acta/
    app.go
    cli.go
    config.go
    event.go
    git.go
    identity.go
    journal.go
    pending.go
    query.go
    record.go
    report.go
    service.go
    transport.go
    ulid.go
    verify.go
```

Files may be combined when two concerns remain small. New concern files require
a clear responsibility and must remain in `internal/acta/` unless they are
systemd assets or test fixtures.

## Delivery Scope

The delivery set includes:

- `acta/v1` milestone and correction models;
- optional execution harness, harness version, and model provenance;
- repository identity and Git `HEAD` metadata collection;
- distinct Git author, Git committer, and Acta recording timestamps;
- self-contained correction events that preserve an append-only audit link;
- UTC monthly JSONL storage and global derived indexes;
- idempotent milestone recording;
- interrupted append recovery through one pending operation;
- bounded JSON IPC over a protected Unix domain socket;
- per-user systemd service and socket activation assets;
- `record`, `show`, `search`, `report`, `correct`, and `verify` operations;
- daily active-segment and monthly full-journal verification assets;
- focused, integration, concurrency, and race tests.

## Out of Scope

- non-`HEAD` backfill;
- durable agent-session identity or stored session books;
- cross-host replication or conflict resolution;
- semantic search, embeddings, or external network services;
- historical segment checksums;
- durable generated report artifacts;
- installation into or activation of the live user system during normal tests;
- `dev -> main` alignment or mirror publication.

## Ordered Slices

### Slice 1: Durable Journal Core

Purpose: Define and verify the canonical event and persistence layer without a
transport dependency.

Branch: `feat/acta-journal-core`

Expected PR title: `feat: add durable monthly journal core`

Report: `ops/default/report/slices/acta-initial-local-chronicle/slice-1-durable-journal-core.md`

Tasks:

- `T1.1` Define configuration, self-contained milestone and correction values,
  optional execution provenance, validation, ULID generation, repository
  identity normalization, and table-driven tests.
- `T1.2` Implement monthly segment selection, append-only journal loading,
  global event and milestone indexes, and idempotent append with tests.
- `T1.3` Implement the pending-operation transaction, interrupted append
  recovery, corruption blocking, and concurrent writer tests.

Validation:

```sh
go fmt ./...
go test ./internal/acta
go test -race ./internal/acta
go test ./...
```

Prerequisites: Approved plan and tracker committed to current `dev`.

### Slice 2: Local Service and Recording

Purpose: Deliver the complete `acta record` path through the local service.

Branch: `feat/acta-record-service`

Expected PR title: `feat: add local service recording workflow`

Report: `ops/default/report/slices/acta-initial-local-chronicle/slice-2-local-service-and-recording.md`

Tasks:

- `T2.1` Implement the bounded `acta/v1` JSON envelopes, Unix socket client and
  server, deadlines, limits, stable errors, and transport tests.
- `T2.2` Implement read-only Git metadata and timestamps, remote sanitization,
  optional agent-supplied execution provenance, and the record workflow with
  idempotency and failure tests.
- `T2.3` Add root `main.go`, CLI dispatch for `record` and `serve`, graceful
  service lifecycle, and end-to-end client-service tests.
- `T2.4` Add per-user `acta.service` and `acta.socket` assets and validate socket
  activation without enabling or changing the live user service.

Validation:

```sh
go fmt ./...
go test ./internal/acta
go test -race ./internal/acta
go test ./...
systemd-analyze verify assets/systemd/user/acta.service assets/systemd/user/acta.socket
```

Prerequisites: Slice 1 delivered and closed on `dev`. `systemd-analyze` may be
reported as unavailable when the host does not provide it.

### Slice 3: Bounded Retrieval and Corrections

Purpose: Deliver explicit retrieval, reporting, and append-only correction
behavior over the same service boundary.

Branch: `feat/acta-retrieval`

Expected PR title: `feat: add bounded chronicle retrieval`

Report: `ops/default/report/slices/acta-initial-local-chronicle/slice-3-bounded-retrieval-and-corrections.md`

Tasks:

- `T3.1` Implement cross-segment filtering, deterministic ordering, text
  matching, self-contained correction chains, harness and model filters, Git
  time filters, and bounded query results with tests.
- `T3.2` Implement `show`, `search`, and `correct` service and CLI operations,
  including not-found, invalid-target, count-limit, and byte-limit behavior.
- `T3.3` Implement `report` with explicit intervals, local or supplied
  timezones, grouped output, JSON output, and cross-month tests.

Validation:

```sh
go fmt ./...
go test ./internal/acta
go test -race ./internal/acta
go test ./...
```

Prerequisites: Slice 2 delivered and closed on `dev`.

### Slice 4: Integrity and Operational Completion

Purpose: Complete integrity verification and the safe scheduled operating
contract.

Branch: `feat/acta-integrity`

Expected PR title: `feat: add journal integrity verification`

Report: `ops/default/report/slices/acta-initial-local-chronicle/slice-4-integrity-and-operational-completion.md`

Tasks:

- `T4.1` Implement active and full verification with exact segment, line, and
  byte-offset diagnostics for all specified invariants.
- `T4.2` Add `verify` CLI and service operations plus the per-user daily timer
  with monthly full-verification behavior.
- `T4.3` Add end-to-end acceptance coverage for recovery, concurrency, IPC
  limits, self-contained cross-segment correction, execution filtering,
  verification, and restart-safe indexing.

Validation:

```sh
go fmt ./...
go test ./internal/acta
go test -race ./internal/acta
go test ./...
systemd-analyze verify assets/systemd/user/acta.service assets/systemd/user/acta.socket assets/systemd/user/acta-verify.service assets/systemd/user/acta-verify.timer
```

Prerequisites: Slice 3 delivered and closed on `dev`. `systemd-analyze` may be
reported as unavailable when the host does not provide it.

## Delivery Rules

- Execute slices in numeric order.
- Use the exact branch, PR title, and report path recorded for each slice.
- Use one Conventional Commit per task when practical.
- Target every slice PR to `dev`.
- Keep each slice compiling and passing its required checks.
- After a slice merge, close its tracker and report state on `dev` before
  activating the next slice.
- Do not create an additional integration or tuning slice without an approved
  plan and tracker update.
- The final slice report records the complete acceptance run. No separate cycle
  report is required for this delivery set.

## Approval Gate

Implementation may start only after:

1. The maintainer approves this complete delivery map.
2. This plan and its tracker have status `approved`.
3. The approved artifacts are committed to `dev` and match `origin/dev`.
4. Slice 1 is the active slice in the tracker.
