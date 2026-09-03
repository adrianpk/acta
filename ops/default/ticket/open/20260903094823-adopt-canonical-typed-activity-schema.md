---
id: TKT-20260903094823
title: Adopt one canonical typed activity schema
status: open
kind: task
severity: medium
priority: normal
scope: persistence
tags: activity-schema, migration, reporting
source: chat
reported_at: 2026-09-03T09:48:23Z
commits: []
---

# Adopt one canonical typed activity schema

## Context

Acta currently persists self-contained Git milestone events and returns the
complete canonical event in machine-readable search and report results. The
journal remains small enough that storage cost is not an operational concern,
but the current shape repeats repository and Git context on every event and
produces verbose retrieval payloads with limited additional information.

The domain also needs to represent professional activity that is not
necessarily a Git commit. Initial peer record types must include commits, pull
requests, human worklogs, and user-directed notes. A commit has no privileged
role in the schema. References can be absent or contain several external
objects, and a worklog or note remains valid without one.

Acta must have one canonical schema. The implementation must not leave old and
new journal schemas active together, expose schema-version selection, or make
normal readers branch between schema generations. Existing history must be
transformed completely before the new runtime resumes writes.

## Expected outcome

- Replace the Git-milestone-shaped event with one typed activity record shared
  by commits, pull requests, worklogs, and notes.
- Keep each event type's data inside a strict type-specific properties object.
  Do not make repository, Git, or evidence fields mandatory for every event.
- Store one stable scope identity without repeating a local root, normalized
  remote, and derived repository name as competing identifiers.
- Distinguish when activity occurred from when Acta recorded it.
- Represent work duration explicitly and never infer it from commits, pull
  requests, event counts, or intervals between events.
- Keep optional typed references and human or agent provenance distinct from
  the activity itself.
- Return purpose-specific compact projections from list and report operations.
  Return the complete canonical record only from an explicit detail operation.
- Deterministically migrate all existing journal history to the canonical
  shape before activating the new runtime.
- Complete deployment and migration through a short, reversible writer-freeze
  window without implementing online migration, dual writes, or a transaction
  engine.

## Canonical record schema

The canonical persisted record has one common envelope. Optional fields are
omitted rather than written as `null` or empty values.

```json
{
  "id": "01K4M6F4M8GJ7Q5AX4Y9R2V3NC",
  "type": "worklog",
  "occurred_at": "2026-09-03T08:00:00Z",
  "recorded_at": "2026-09-03T09:41:12Z",
  "summary": "Configured Grafana alert routing",
  "scope": {
    "type": "system",
    "name": "Production observability",
    "identity": "system:production-observability"
  },
  "properties": {
    "duration_minutes": 240
  },
  "references": [
    {
      "type": "url",
      "url": "https://grafana.example.test/alerting",
      "label": "Grafana alerting"
    }
  ],
  "provenance": {
    "type": "human"
  }
}
```

### Common fields

| Field | Contract |
| --- | --- |
| `id` | Service-assigned, time-sortable, globally unique record identifier. |
| `type` | Typed activity discriminator. Initial values are `commit`, `pull_request`, `worklog`, and `note`. |
| `occurred_at` | UTC time at which the recorded activity occurred. |
| `recorded_at` | UTC publication time assigned by the Acta service. |
| `summary` | Required concise human-readable account of the activity. |
| `scope` | Optional stable project, repository, system, or other work scope. |
| `properties` | Strict object selected by `type`; omitted only when the selected type has no additional properties. |
| `references` | Optional ordered list of typed related objects. |
| `provenance` | Optional description of how the record was supplied. It is not evidence of identity or authorization. |
| `supersedes` | Optional identifier of the effective record replaced by this record. |
| `correction_reason` | Required when `supersedes` is present and absent otherwise. |

There is no `schema_version` field. One release owns one accepted persisted
shape, and the normal runtime rejects legacy or mixed-shape journals with a
clear migration-required error.

### Scope

`scope` uses one stable identity and one display name:

```json
{
  "type": "project",
  "name": "Example project",
  "identity": "project:example"
}
```

- `type`, `name`, and `identity` are required when `scope` is present.
- `identity` is stable across checkout moves and machine changes.
- Absolute worktree paths are not canonical record data.
- A normalized repository remote can be used as the identity when one exists.
- A repository without a remote receives a stable local identity that does not
  embed its checkout path.

### Type-specific properties

`properties` is a discriminated union, not an open generic map. Initial
variants are:

- `commit`: full commit object ID and subject, with optional author time when it
  differs meaningfully from `occurred_at`.
- `pull_request`: stable forge and repository identity, pull request number,
  relevant action, and URL when available.
- `worklog`: explicit `duration_minutes`. Acta never derives this value from
  another event or from elapsed time between records.
- `note`: no additional properties in the initial schema. Its required summary
  is the recorded content.

Each variant is a peer. A commit is not the canonical parent, default evidence
shape, or identity basis for another event type.

### References

References are optional related objects that any activity type may carry. The
initial variants are `git_commit`, `pull_request`, and `url`. A reference can
support retrieval or inspection, but Acta does not require one to accept an
otherwise valid worklog or note.

Commit and pull-request identity and validation come from their own typed
properties. They are not represented by mandatory self-referential entries in
`references`.

### Provenance

`provenance.type` initially supports `human`, `agent`, and `integration`.
Harness, harness version, and model are optional agent details. Provenance is
not displayed in default human output and does not turn a declaration into
independently verifiable evidence.

### Corrections

A correction is a complete replacement activity record with the same activity
`type`, a new `id` and `recorded_at`, plus `supersedes` and
`correction_reason`. Correction is a relation between records, not a separate
activity type. Effective-event lookup follows the append-only supersession
chain.

## Legacy history mapping

The migration preserves record IDs, ordering, UTC segment placement, recording
times, commit object IDs, subjects, notes, correction relationships, and known
execution provenance.

Each legacy milestone maps as follows:

| Legacy value | Canonical value |
| --- | --- |
| `type: milestone` | `type: commit` |
| `recorded_at` | `recorded_at` unchanged |
| `git.committed_at` | `occurred_at` |
| non-empty `note` | `summary` |
| empty `note` | `git.subject` as `summary` |
| repository identity and name | one `repository` scope |
| `git.commit` and `git.subject` | `commit` properties |
| differing `git.authored_at` | optional author time in `commit` properties |
| `execution` | optional `provenance` |
| `target_id` | `supersedes` |
| `reason` | `correction_reason` |

Legacy repository roots, duplicate remotes, branches, and observed worktree
cleanliness are intentionally not copied. They are transient capture context,
not durable activity evidence.

The migrator must fail before publication if it encounters an event that
cannot be mapped without information loss outside these explicit removals. It
must never silently skip or partially reinterpret an event.

## Compact retrieval views

- `record` returns and prints only a receipt containing status, record ID,
  activity type, optional scope name, and an abbreviated type-specific
  identifier when one exists.
- Default `search` output contains recording time, type, scope name, summary,
  and an abbreviated primary reference.
- Default `report` output is summary-first and grouped without repeating scope
  data inside every entry.
- Machine-readable `search` and `report` responses use dedicated projection
  types and do not embed complete canonical records.
- `show --json` is the explicit operation for retrieving one complete
  canonical record.
- Full provenance, properties, and references are omitted from normal updates
  unless the caller explicitly requests record detail.

## Deployment and migration procedure

The implementation must provide a repository-owned, deterministic migration
command. The command may contain a legacy decoder, but the normal service,
recording, query, and reporting paths must contain only the canonical model.

Prepare everything possible before interrupting the writer:

1. Build the candidate runtime and migration command without replacing the
   installed runtime.
2. Run unit and integration tests against representative milestone,
   correction, remote-repository, local-repository, and malformed-journal
   fixtures.
3. Copy the live journal to an isolated rehearsal location, migrate it, and run
   complete canonical verification there.
4. Record the expected source segment count, line count, record IDs, effective
   record count, correction chains, and checksums in a machine-readable cutover
   manifest.

Perform the final cutover during a short controlled write freeze:

1. Stop socket activation and the service, then verify that no Acta writer is
   running.
2. Create a recoverable snapshot of the installed binary and complete journal.
3. Migrate the frozen journal into a sibling staging location. Never rewrite
   the canonical journal in place.
4. Verify every staged segment with the candidate runtime and compare its
   invariants with the cutover manifest.
5. Replace the installed binary and canonical journal from the staged
   artifacts using same-filesystem renames where available.
6. Start socket activation and the service only after both replacements
   succeed.
7. Run read-only smoke checks for verification, search, show, and report before
   allowing normal recording to resume.

The binary and journal cannot be exchanged as one filesystem transaction. The
cutover therefore uses an explicit checkpoint and ordered, reversible renames.
If any step after the writer freeze fails, keep the service stopped, restore
both the prior binary and journal snapshot, verify the restored installation,
and only then restart it.

The legacy snapshot is a temporary rollback artifact, not a supported second
journal. Remove it only after explicit operator acceptance of the migrated
runtime and journal.

## Validation

- Every canonical segment is accepted by the new verifier and contains only
  the canonical record shape.
- The normal runtime rejects legacy and mixed-shape journals instead of
  partially loading them.
- Migration preserves all source IDs and produces the same source and target
  record counts.
- Migration preserves chronological ordering, segment placement, effective
  correction chains, commit properties, summaries, and known provenance.
- Re-running migration from the same frozen source produces byte-identical
  canonical segments and manifest data.
- Search and report semantic fixtures return equivalent activity after
  migration, accounting for the explicit field removals and projection shape.
- Default list and report output does not include absolute paths, full
  canonical records, or repeated grouped scope data.
- Failure injection before and after each cutover rename demonstrates complete
  rollback to the prior runtime and journal.
- Repository quality gates and `git diff --check` pass.

## Non-goals

- Supporting several active journal schema versions.
- Online migration, dual reads, dual writes, or a database transaction layer.
- Inferring work duration or causal relationships from any recorded activity.
- Automatically recording conversations, prompts, hidden reasoning, terminal
  history, diffs, file contents, or indiscriminate telemetry.
- Removing the temporary rollback snapshot before operator acceptance.
