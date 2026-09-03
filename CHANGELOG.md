# Changelog

All notable changes to Acta are documented in this file.

This project uses Semantic Versioning.

## Unreleased

## 0.1.0 - 2026-09-03

### Added

- Record completed activity as peer `commit`, `pull_request`, `worklog`, and
  `note` events, with optional scope, references, occurrence time, and
  provenance.
- Search activity across repositories and other named contexts by type, scope,
  time, duration, source, harness, model, commit, pull request, or text.
- Inspect complete records and append corrections without changing the
  original journal history.
- Build occurrence-date reports with event counts and explicit worklog-minute
  totals for timesheet preparation.
- Run Acta as a socket-activated per-user service backed by a local,
  append-only monthly JSONL journal.
- Verify the active journal segment, the full journal, or the scheduled
  verification scope with precise corruption diagnostics.
- Migrate a frozen legacy journal offline into the canonical typed format with
  deterministic output and a machine-readable cutover manifest.
- Connect workspace agents through repository capture policies while keeping
  human-directed notes and non-Git work under explicit user control.

### Changed

- `acta record` now requires an explicit event type. A Git commit is one type
  of activity, not a prerequisite or default for other records.
- Acta now uses one strict typed event format for normal operation. Existing
  legacy journals must be migrated before use with this release.
