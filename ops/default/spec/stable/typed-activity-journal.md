# Typed Activity Journal Specification

Status: accepted
Approval: approved
Source ticket: `ops/default/ticket/closed/20260903094823-adopt-canonical-typed-activity-schema.md`
Model: `ops/default/spec/stable/typed-activity-journal-model.md`
Supersedes: the product and event-model contract in
`ops/default/spec/drafts/acta-system-specification-draft.md`

## 1. Purpose

Acta records durable, typed professional activity for later retrieval and
reporting. A Git commit is one activity type. It is not the parent model, the
default event, or privileged evidence for other activity.

The initial canonical activity types are:

1. `commit` for a completed Git commit;
2. `pull_request` for a meaningful pull-request lifecycle action;
3. `worklog` for explicitly declared work and duration;
4. `note` for user-directed information that does not fit a more specific
   type.

Acta remains a deliberate journal. It does not infer activity from background
observation and does not automatically capture every Git or agent action.

## 2. Core outcomes

The product must support two complementary uses:

1. Humans record work that may not produce a repository artifact and use
   explicit worklog duration to prepare timesheets.
2. Humans and agents locate typed activity across scopes and follow optional
   references to related systems or artifacts.

Acta reports recorded facts and transparent totals. It must not infer effort,
duration, causality, or productivity from commits, pull requests, event counts,
or time between events.

## 3. Product principles

### 3.1 Peer event types

Every canonical line is one activity event with a shared envelope. The event
`type` selects a strict properties variant. Each variant is a peer.

Type-specific identity belongs to the selected properties object. A commit
object ID belongs to `commit` properties. A pull-request number belongs to
`pull_request` properties. Neither is stored as mandatory self-referential
evidence.

### 3.2 Activity, scope, references, and provenance

The model separates four concerns:

1. Activity states what occurred through `type`, `occurred_at`, `summary`, and
   type-specific `properties`.
2. Scope identifies the repository, project, system, client, or other work
   context.
3. References identify optional related objects that can help later inspection.
4. Provenance describes how the record was supplied. Provenance is not
   authorization or independently verified identity.

An otherwise valid `worklog` or `note` does not require a scope or reference.

### 3.3 Deliberate recording

Agents may record commits under a workspace capture policy. Agents must record
a `note` only when the user explicitly asks Acta to preserve the information.
This permission applies to the requested information only. It does not permit
automatic conversation capture.

### 3.4 One canonical persisted shape

After migration, the installed runtime accepts one canonical journal shape.
It does not expose schema selection, dual reads, or dual writes. The persisted
record does not contain a `schema_version` field. Legacy or mixed-shape
journals fail with a migration-required error before normal writes resume.

## 4. Actors and boundaries

- The user owns the journal and supplies or approves professional activity.
- A working agent can submit records under workspace policy or direct user
  instruction.
- The CLI validates local inputs and derives type-specific facts when defined.
- The per-user Acta service is the sole canonical journal writer.
- Retrieval clients issue explicit, bounded queries.
- The migration command transforms an isolated legacy journal without using
  the running service.

Acta does not replace Git, a forge, a timesheet system, Grafana, a CI system, or
another source system. It records activity and optional references to them.

## 5. Recording interface

The CLI provides explicit type selection:

```sh
acta record commit [flags] [summary]
acta record pull-request --forge <url> --repository <owner/name> --number <n> --action <action> [flags] <summary>
acta record worklog --duration <minutes> [flags] <summary>
acta record note [flags] <summary>
```

`acta record` without a type is invalid. This keeps `commit` from becoming an
implicit default.

### 5.1 Common inputs

All record commands accept:

- `--at <RFC3339>` for `occurred_at`;
- `--scope-type <type>`, `--scope-name <name>`, and
  `--scope-identity <identity>` as an all-or-none custom scope;
- repeatable `--reference-url <url>`;
- repeatable `--reference-commit <scope-identity>@<object-id>`;
- repeatable `--reference-pull-request <url>`;
- `--source <human|agent|integration>`;
- `--harness`, `--harness-version`, and `--model` when agent provenance is
  known reliably.

The service assigns `recorded_at` and the event ID. If `--at` is omitted for a
non-Git record, `occurred_at` equals the service-assigned `recorded_at`.

Custom scope fields must either all be absent or all be present. A type that
derives a repository scope rejects conflicting custom scope fields.

### 5.2 Commit recording

`commit` runs inside a Git worktree and records `HEAD`. The CLI derives the
commit object ID, subject, commit time, optional differing author time, and a
stable repository scope without changing Git state.

The normalized remote is the repository scope identity when available. A
repository without a remote uses a deterministic identity derived from its
first root commit. The identity does not contain an absolute checkout path.

`occurred_at` is the commit time and `--at` is invalid. The optional positional
summary overrides the commit subject as the event summary. Without it, the
commit subject becomes the summary.

Recording the same repository scope and commit object ID is idempotent.

### 5.3 Pull-request recording

`pull_request` requires a normalized forge URL, repository identity, positive
pull-request number, lifecycle action, and summary. Initial actions are
`opened`, `reviewed`, `merged`, and `closed`.

The CLI derives a repository scope from the forge and repository values. The
URL is optional because a stable forge, repository, and number identify the
pull request. Each action is a separate recordable event. Repeating the same
forge, repository, number, and action is idempotent.

The command does not contact the forge or claim that supplied state was
externally verified.

### 5.4 Worklog recording

`worklog` requires a summary and `--duration` as an integer number of minutes.
Duration is explicit stored data. It is never derived from another event or
from elapsed time between records.

A worklog can describe configuration, investigation, coordination, research,
maintenance, or other professional work without a commit or pull request.
Scope and references are optional.

### 5.5 Note recording

`note` requires a summary and has no type-specific properties in the initial
model. Scope and references are optional.

The note type is the deliberate fallback when no more specific type applies.
It must not be named or presented as `misc`, `manual`, or `ad_hoc` in the
canonical model.

### 5.6 Recording result

On success, the CLI prints a compact receipt containing whether a new event was
recorded, its ID, type, optional scope name, and an abbreviated type-specific
identifier when one exists.

Recording fails without appending when required fields are absent, fields are
invalid, the service is unavailable, the journal is corrupt, or publication
cannot complete durably.

## 6. Corrections

Published events are immutable. `acta correct` appends a complete replacement
event and sets `supersedes` to the current effective event. The replacement
keeps the same activity type. `correction_reason` is required.

A correction is a relation between events, not an activity type. Queries and
reports return the last effective event in each supersession chain by default.
Raw verification and migration retain every event in append order.

## 7. Retrieval interface

### 7.1 Show

```sh
acta show --id <event-id> [--json]
```

`show` is the operation that returns the complete canonical effective event.
An ID for a superseded event resolves to that exact raw event only when
`--raw` is supplied.

### 7.2 Search

```sh
acta search [filters] [--json]
```

Search supports bounded filtering by activity type, scope name or identity,
commit object ID, pull-request identity, text, occurrence time, recording time,
duration, source, harness, and model.

Default text and JSON search output use a compact projection containing event
ID, occurrence time, recording time, type, optional scope name, summary,
optional duration, and abbreviated type-specific identifier. Search does not
embed complete properties, references, or provenance.

### 7.3 Report

```sh
acta report --from <date> --to <date> [--timezone <name>] [--type <types>] [--json]
```

Report intervals are local-calendar, inclusive at `from` and exclusive at
`to`. Reports group effective activity by local occurrence date and scope.
Unscoped activity is grouped under `Unscoped`.

Reports include event counts and sum `duration_minutes` from matching worklogs.
No other event type contributes inferred duration. Compact report entries use
the same summary-first projection as search.

Search and report return at most 1,000 records and stop before the bounded
service response would exceed 8 MiB. They report truncation explicitly.

## 8. Service and storage

Acta remains a per-user service activated through the owner-only Unix socket at
`$XDG_RUNTIME_DIR/acta/acta.sock`. The runtime directory uses mode `0700` and
the socket uses mode `0600`.

Canonical events remain compact UTF-8 JSON objects stored one per line under:

```text
$HOME/.acta/events/<UTC-year>/<UTC-month>.jsonl
```

Segment selection uses `recorded_at`. Completed earlier monthly segments are
immutable. All derived indexes are disposable and rebuildable.

The existing bounded pending-append transaction remains the only automatic
write recovery. Verification remains read-only and reports exact segment,
line, and byte-offset failures when available.

## 9. Migration interface

The candidate executable provides an offline command:

```sh
acta migrate --source <legacy-journal-root> --destination <canonical-journal-root> --manifest <path>
```

The source and destination must be distinct. The destination must not contain
journal segments. Migration never rewrites the source. It decodes only the
known legacy milestone/correction shape and emits only the canonical typed
shape.

The command verifies the complete source, migrates deterministically, verifies
the complete destination, and writes a machine-readable manifest containing
segment counts, line counts, event IDs, effective-event counts, correction
chains, and source and target SHA-256 checksums.

Re-running migration from the same frozen source produces byte-identical
segments and manifest content. Any unsupported, corrupt, or lossy mapping
fails before the destination is accepted.

## 10. Installation and cutover

Implementation and migration rehearsal occur while the installed runtime keeps
serving the legacy journal. Final cutover uses this order:

1. Build the candidate executable without replacing the installed binary.
2. Copy the live journal to an isolated rehearsal location.
3. Migrate and fully verify the rehearsal copy twice for determinism.
4. Record the frozen-source invariants and checksums.
5. Stop socket activation, the service, and the verification timer.
6. Confirm that no Acta writer is active.
7. Create recoverable snapshots of the installed executable and complete
   journal.
8. Migrate the frozen journal into a sibling staging location.
9. Verify the staging journal and compare it with the rehearsal manifest.
10. Replace the executable and journal using ordered same-filesystem renames.
11. Restore and start the user units.
12. Run full verification and smoke checks for record, search, show, report,
    and correction behavior.

If any cutover step fails, the service remains stopped until both the prior
binary and prior journal snapshot are restored and verified. Rollback snapshots
remain until explicit operator acceptance.

## 11. Security and privacy

- The journal remains private local user data.
- Summaries and references must not contain credentials, tokens, personal data,
  hidden reasoning, prompts, or conversation transcripts.
- Generic URL references reject user information and control characters.
- Git remotes and forge URLs are normalized without credentials.
- Recording never modifies a repository or external service.
- Search, report, show, verification, and migration do not contact external
  services.

## 12. Accepted decisions

- The canonical entity is a typed activity event.
- `commit`, `pull_request`, `worklog`, and `note` are initial peer types.
- Type-specific properties are strict and separate from optional references.
- `note` is the explicit user-directed fallback type.
- Worklog duration is stored explicitly in minutes.
- `occurred_at` and `recorded_at` are distinct.
- Scope is optional and uses one stable identity plus a display name.
- Corrections are append-only supersession relations, not an event type.
- Normal search and report results are compact projections.
- `show` owns complete record retrieval.
- One canonical persisted shape replaces the legacy schema after deterministic
  offline migration.

## 13. Boundaries

Acta does not infer time, activity, significance, causality, or productivity.
It does not monitor applications, terminals, repositories, agent sessions, or
external systems. It does not submit timesheets or mutate referenced systems.

## 14. Acceptance criteria

The specification is satisfied when the implementation demonstrates that:

1. All four peer event types can be recorded and retrieved through the service.
2. `acta record` without an explicit type fails.
3. Commit and pull-request properties are type-owned, not mandatory references.
4. A worklog and note can be recorded without Git, scope, or references.
5. A worklog requires explicit bounded duration and reports contribute exactly
   that duration to totals.
6. An agent-facing note policy requires explicit user direction.
7. Occurrence time and recording time remain independently queryable.
8. Commit and pull-request natural duplicates are idempotent.
9. Corrections preserve append-only history and effective-event queries.
10. Search and report return compact bounded projections.
11. Show returns one complete canonical record by event ID.
12. Active and full verification enforce the canonical shape and reject legacy
    or mixed journals.
13. Migration preserves IDs, append order, segment placement, commit identity,
    summaries, known provenance, and correction chains.
14. Migration is deterministic and never rewrites its source.
15. Rehearsal, cutover, rollback, installation, and post-install smoke checks
    satisfy the procedure in this specification.
