# Slice 2: Local Service and Recording

Status: reviewing
Delivery set: Acta Initial Local Chronicle
Plan: `ops/default/plan/acta-initial-local-chronicle.md`
Tracker: `ops/default/tracker/acta-initial-local-chronicle.md`
Branch: `feat/acta-record-service`
PR: pending

## Purpose

Deliver the complete `acta record` path through a protected per-user service
without allowing clients to write canonical storage directly.

## Delivered Behavior

- Exchanges one bounded `acta/v1` JSON request and response per Unix stream
  connection.
- Enforces the 64 KiB request limit, 8 MiB response limit, deadlines, stable
  error codes, and socket permissions.
- Derives repository root, sanitized identity, full `HEAD` object ID, branch,
  subject, author time, committer time, and complete worktree state without
  modifying Git.
- Accepts an optional note plus reliably supplied harness, harness version, and
  model provenance.
- Assigns event IDs and UTC recording times in the service and preserves
  idempotent milestone recording across service restarts.
- Provides `record` and `serve` commands through a thin root `main.go`.
- Adopts one systemd-provided Unix listener or creates the configured socket for
  direct service execution.
- Provides per-user service and socket units without installing or activating
  them.

## Implementation Notes

The client closes its write side after one request and reads until the service
closes the response side. The service handles concurrent connections while the
journal remains the serialized canonical writer.

Git inspection uses read-only Git commands. An unusable `origin` falls back to
a canonical local repository identity instead of persisting an unsafe remote.
The service validates the complete repository and Git evidence again before it
publishes an event.

The systemd service receives exactly one descriptor starting at file descriptor
3. Tests adopt a duplicated temporary Unix listener and exercise it through the
normal client and server.

## Contracts Added or Changed

- IPC requests and responses use protocol `acta/v1`.
- `record` returns `recorded: false` with the existing event for an idempotent
  duplicate.
- `acta record [flags] [note]` supports `--harness`, `--harness-version`, and
  `--model`; version or model values require a harness.
- `$XDG_RUNTIME_DIR/acta` uses mode `0700`; `acta.sock` uses mode `0600`.
- The socket unit uses `%t/acta/acta.sock` with `Accept=no` and descriptor name
  `acta`.
- The service unit uses an explicit PATH that includes `%h/.local/bin`.

## Files of Interest

- `main.go`
- `internal/acta/transport.go`
- `internal/acta/git.go`
- `internal/acta/record.go`
- `internal/acta/cli.go`
- `internal/acta/service.go`
- `internal/acta/activation.go`
- `assets/systemd/user/acta.service`
- `assets/systemd/user/acta.socket`

## Validation

- `go fmt ./...` (pass)
- `make lint-strict` (pass, 0 issues)
- `go test ./internal/acta` (pass)
- `go test -race ./internal/acta` (pass)
- `go test ./...` (pass)
- `systemd-analyze verify assets/systemd/user/acta.service assets/systemd/user/acta.socket` (pass)

No live user service was installed, enabled, started, stopped, or restarted.

## Risks and Follow-ups

- The binary and user units still require an explicit installation step before
  normal socket activation is available.
- Retrieval, correction, and report operations remain in Slice 3.
- The initial README exists, but its detailed user-guide links remain tracked as
  separate documentation work.
