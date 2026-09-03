# Slice 4: Deterministic Journal Migration

Status: reviewing
Delivery set: Typed Activity Journal
Plan: `ops/default/plan/typed-activity-journal.md`
Tracker: `ops/default/tracker/typed-activity-journal.md`
Branch: `feat/typed-activity-migration`
PR: `#10`

## Purpose

Deliver a strict offline conversion from the installed legacy journal to the
canonical typed journal, including deterministic evidence and non-overwriting
publication behavior required before live cutover.

## Delivered Behavior

- Decodes only exact `acta/v1` milestone and correction records inside the
  migration boundary while normal runtime decoding remains canonical-only.
- Validates legacy event, repository, Git snapshot, execution, and correction
  invariants before mapping either legacy type to a canonical `commit` event.
- Preserves event IDs, recording times, Git occurrence times, repository
  identity, Git object evidence, agent provenance, and correction chains.
- Rejects unknown fields, malformed framing, invalid snapshots, lossy summary
  mappings, unexpected source files, and source mutation during migration.
- Writes and validates a private sibling staging tree before publishing a new
  destination without overwriting an existing destination or manifest.
- Emits a deterministic machine-readable manifest with journal and segment
  digests, ordered event IDs, effective counts, and correction chains.
- Removes partial output when destination or manifest publication fails.
- Provides an offline `acta migrate` command that does not require the service
  socket or runtime directory.
- Documents command behavior, storage evidence, and a two-run rehearsal that
  must precede the final frozen-journal cutover.

## Implementation Notes

The migrator accepts distinct absolute source, destination, and manifest paths.
The source must be a real directory containing only monthly `YYYY/MM.jsonl`
segments. The destination and manifest must not exist, and all three paths must
be non-nested so publication cannot modify or obscure the source.

Each legacy line is strictly decoded, validated against the legacy journal
index, mapped to the canonical envelope, validated again, and written in source
segment order. The completed staging journal is then read through the normal
canonical decoder and compared with the source summary before publication.
Source file digests are checked again at the publication boundary.

Publication atomically reserves the destination, moves verified segment
directories into it, and links the staged manifest into place. Failures after
reservation remove the new output while leaving the source unchanged. The
manifest excludes timestamps and machine-specific absolute paths so equivalent
inputs produce byte-identical journals and manifests.

## Contracts Added or Changed

- Adds `acta migrate --source <legacy-journal-root> --destination
  <canonical-journal-root> --manifest <path>` as an offline-only command.
- Requires every migration path explicitly; no live journal or service path is
  inferred.
- Maps legacy milestones and corrections to canonical `commit` events, using
  `git.committed_at` as `occurred_at` and retaining correction links through
  `supersedes` and `correction_reason`.
- Preserves legacy repository identity as repository scope while intentionally
  dropping the former machine-local root and duplicated remote field.
- Publishes destination files with mode `0600` and directories with mode
  `0700`.

## Files of Interest

- `internal/acta/migration.go`
- `internal/acta/migration_cli.go`
- `internal/acta/migration_test.go`
- `internal/acta/migration_cli_test.go`
- `docs/how-to/rehearse-a-legacy-journal-migration.md`
- `docs/reference/command-line.md`
- `docs/reference/storage-and-event-data.md`

## Validation

- `make check` (pass)
- `make test` (pass)
- `make docs-check` (pass)
- `make verify` (pass, 0 lint issues)
- `go test -race ./...` (pass)
- `go test -count=1 ./internal/acta -run 'Test(Migrator|ParseMigration|AppMigrates|Legacy)'` (pass)
- `systemd-analyze verify assets/systemd/user/acta.service assets/systemd/user/acta.socket assets/systemd/user/acta-verify.service assets/systemd/user/acta-verify.timer` (pass)
- `git diff --check` (pass)

No candidate binary or systemd unit was installed, enabled, started, stopped,
or restarted. The active runtime and live journal remain on the legacy schema.

## Risks and Follow-ups

- The actual rehearsal requires the exact integrated candidate after this slice
  is merged into `dev`; the repository guide defines that gate but this slice
  does not execute live-state operations.
- Installation and live migration remain gated on all slices being delivered,
  exact-candidate staging CI, recoverable snapshots, successful rehearsal, and
  a confirmed writer freeze.
- Rollback snapshots must remain after cutover until the maintainer explicitly
  accepts their removal.
