# Typed Activity Journal Model

Status: accepted
Approval: approved
Specification: `ops/default/spec/stable/typed-activity-journal.md`
Source ticket: `ops/default/ticket/closed/20260903094823-adopt-canonical-typed-activity-schema.md`

## 1. Canonical event

Every canonical JSONL line contains one complete event. Optional fields are
omitted instead of being encoded as `null` or empty values.

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
      "url": "https://grafana.example.test/alerting"
    }
  ],
  "provenance": {
    "type": "human"
  }
}
```

The canonical field order is:

```text
id
type
occurred_at
recorded_at
summary
scope
properties
references
provenance
supersedes
correction_reason
```

The persisted shape has no `schema_version` field. Canonical decoding rejects
unknown fields and any legacy field such as `schema_version`, `repository`,
`git`, `note`, `target_id`, or `reason`.

## 2. Common values

### 2.1 Event ID

`id` is a service-assigned ULID. It is globally unique and time-sortable. The
ULID timestamp must equal `recorded_at` at millisecond precision.

### 2.2 Activity type

`type` is one of:

```text
commit
pull_request
worklog
note
```

No type is the parent, default, or evidence form of another type.

### 2.3 Timestamps

`occurred_at` and `recorded_at` are required RFC 3339 timestamps normalized to
UTC. `recorded_at` is assigned by the service and selects the monthly segment.
`occurred_at` can precede or equal `recorded_at`; it must not be later than the
service recording time.

### 2.4 Summary

`summary` is required trimmed UTF-8 plain text. It contains no newline or
control character and is at most 500 Unicode code points.

## 3. Scope

Scope is optional and has this shape:

```json
{
  "type": "project",
  "name": "Example project",
  "identity": "project:example"
}
```

When present, all fields are required, trimmed, single-line strings. `type` and
`identity` use normalized exact matching. `name` is a display value.

Initial scope types are `repository`, `project`, `system`, `client`, and
`other`. Scope identity is opaque to generic journal logic. Type-specific
recorders own deterministic identity construction.

A commit scope uses a normalized credential-free remote when one exists.
Without a remote, it uses `git-root:<root-commit-object-id>`. When `HEAD` has
several root commits, the lexically first full object ID is used. A pull-request
scope uses `<normalized-forge-url>/<owner>/<repository>`.

Absolute checkout paths are not canonical event data.

## 4. Type-specific properties

`properties` is a strict discriminated object selected by the event type. The
canonical decoder validates the selected variant and rejects fields owned by a
different variant.

### 4.1 Commit properties

```json
{
  "object_id": "0123456789abcdef0123456789abcdef01234567",
  "subject": "feat: add alert routing",
  "authored_at": "2026-09-03T07:58:00Z"
}
```

- `object_id` is required, lowercase hexadecimal, and 40 or 64 characters.
- `subject` is required, single-line text with no control characters.
- `authored_at` is optional and stored only when it differs from
  `occurred_at`.
- The event requires a `repository` scope.
- `occurred_at` is the Git committer time.

The event does not store checkout root, remote duplication, branch, or observed
worktree cleanliness.

### 4.2 Pull-request properties

```json
{
  "forge": "https://forge.example.test",
  "repository": "team/example",
  "number": 42,
  "action": "merged",
  "url": "https://forge.example.test/team/example/pulls/42"
}
```

- `forge`, `repository`, `number`, and `action` are required.
- `forge` is a normalized absolute HTTP or HTTPS URL without user information,
  query, or fragment.
- `repository` is a normalized `owner/name` identifier.
- `number` is a positive integer.
- `action` is `opened`, `reviewed`, `merged`, or `closed`.
- `url` is optional and must be an absolute HTTP or HTTPS URL without user
  information.
- The event requires the repository scope derived from `forge` and
  `repository`.

### 4.3 Worklog properties

```json
{
  "duration_minutes": 240
}
```

`duration_minutes` is required and is an integer from 1 through 1,440. It is a
human or integration-supplied fact. Acta never derives or adjusts it.

### 4.4 Note properties

A note has no `properties` field in the initial canonical shape. Its summary is
the complete recorded content.

## 5. References

`references` is an optional ordered array with at most 32 items. References do
not define the containing event's type or identity. Exact duplicate references
inside one event are invalid.

### 5.1 Git commit reference

```json
{
  "type": "git_commit",
  "scope_identity": "https://forge.example.test/team/example.git",
  "object_id": "0123456789abcdef0123456789abcdef01234567"
}
```

### 5.2 Pull-request reference

```json
{
  "type": "pull_request",
  "url": "https://forge.example.test/team/example/pulls/42"
}
```

### 5.3 URL reference

```json
{
  "type": "url",
  "url": "https://grafana.example.test/alerting"
}
```

Reference URLs are absolute HTTP or HTTPS URLs without user information or
control characters. Acta does not fetch or verify them.

## 6. Provenance

Provenance is optional:

```json
{
  "type": "agent",
  "harness": "example-harness",
  "harness_version": "1.2.3",
  "model": "provider/example-model"
}
```

`type` is `human`, `agent`, or `integration`. Agent provenance requires
`harness`; harness version and model require agent provenance and a harness.
All values are trimmed single-line identifiers. Provenance does not participate
in event identity or authorization.

When an agent records content explicitly supplied by the user, the semantic
source is `human`. Harness details can be omitted because provenance describes
the source of the record, not the process that transported it.

## 7. Corrections and effective events

A correction is a complete canonical event with the same activity type and
business identity as its current effective target. It adds:

```json
{
  "supersedes": "01K4M6F4M8GJ7Q5AX4Y9R2V3NC",
  "correction_reason": "Correct the reported duration"
}
```

`supersedes` must identify the current effective event in one chain.
`correction_reason` follows the summary text constraints and is required only
when `supersedes` exists.

The replacement can change occurrence time, summary, non-identity properties,
references, or provenance while preserving the same type, scope identity, and
natural identity. Worklog and note corrections can change scope because those
types have no natural key. The original event remains unchanged. Default
search, report, and show results use the last effective event. Raw verification
and migration retain all events.

## 8. Identity and deduplication

Event IDs are globally unique. Natural keys provide idempotency where an
external object has stable identity:

- `commit`: `(type, scope.identity, properties.object_id)`;
- `pull_request`: `(type, properties.forge, properties.repository,
  properties.number, properties.action)`.

At most one correction chain can own a natural key. A repeated record request
returns the existing effective event and `recorded: false`.

`worklog` and `note` have no inferred natural key. Two separately submitted
records remain separate events even when their content matches.

## 9. Storage and append lifecycle

The ordered monthly JSONL segments remain the sole canonical event store:

```text
$HOME/.acta/events/<year>/<month>.jsonl
```

Segment placement uses the UTC year and month of `recorded_at`. Every line is a
compact JSON object followed by `LF` and is at most 8 MiB.

The existing pending operation continues to contain the target segment,
original segment length, and complete canonical line. Recovery can complete an
absent append, accept an already complete append, or replace a matching partial
terminal prefix. Any other mismatch is corruption.

Derived indexes contain:

- event ID to raw event;
- chain member to chain root;
- chain root to effective event;
- commit natural key to chain root;
- pull-request natural key to chain root;
- chronological event order.

All indexes are disposable and rebuilt from canonical segments.

## 10. Query model

The normalized query can contain:

```text
types
scope
commit
pull_request
text
occurred_from
occurred_to
recorded_from
recorded_to
duration_min
duration_max
source
harness
model
limit
order
```

Time ranges are inclusive at `from` and exclusive at `to`. Default order is
descending by `occurred_at`, then `recorded_at`, then `id`. The maximum limit is
1,000.

Text search matches summary, scope name and identity, type-specific display
fields, reference targets, correction reason, and provenance identifiers.

### 10.1 Activity projection

Search and report use this compact projection:

```json
{
  "id": "01K4M6F4M8GJ7Q5AX4Y9R2V3NC",
  "type": "worklog",
  "occurred_at": "2026-09-03T08:00:00Z",
  "recorded_at": "2026-09-03T09:41:12Z",
  "scope_name": "Production observability",
  "summary": "Configured Grafana alert routing",
  "duration_minutes": 240
}
```

Optional fields are omitted. Complete properties, references, provenance, and
correction data appear only in detail output.

### 10.2 Report projection

Reports group projections by local occurrence date and scope. Each group
contains `event_count` and `worklog_minutes`. The report contains top-level
totals with the same two values. Only effective worklog duration contributes to
`worklog_minutes`.

## 11. Legacy mapping

Migration maps each legacy event without changing ID, `recorded_at`, monthly
segment, or append position:

| Legacy field | Canonical field |
| --- | --- |
| `type: milestone` | `type: commit` |
| `git.committed_at` | `occurred_at` |
| non-empty `note` | `summary` |
| empty `note` | `git.subject` copied to `summary` |
| repository name and identity | repository `scope` |
| `git.commit` | `properties.object_id` |
| `git.subject` | `properties.subject` |
| differing `git.authored_at` | `properties.authored_at` |
| `execution` | `provenance` with `type: agent` |
| `target_id` | `supersedes` |
| `reason` | `correction_reason` |

Legacy root paths, duplicate remotes, branches, and worktree cleanliness are
intentionally removed. The manifest records these explicit removals.

Legacy correction events become `commit` events in the same supersession chain.
The migrator validates that their copied Git and repository values match the
target chain before mapping them.

## 12. Migration manifest

The deterministic manifest contains:

```text
format
source_segment_count
source_line_count
source_event_ids
source_effective_count
source_correction_chains
source_checksums
target_segment_count
target_line_count
target_event_ids
target_effective_count
target_correction_chains
target_checksums
removed_legacy_fields
```

Paths in the manifest are relative to their journal roots. Arrays and map keys
use deterministic lexical ordering. The manifest does not contain the absolute
source or destination path.

## 13. Verification invariants

Full verification rejects:

- invalid UTF-8 or line framing;
- unknown, missing, null, or wrong-type fields;
- unknown activity, properties, reference, provenance, or action values;
- invalid event IDs or timestamp relationships;
- segment placement that disagrees with `recorded_at`;
- invalid summaries, scopes, properties, references, or provenance;
- duplicate event IDs or natural keys;
- invalid, stale, cyclic, cross-type, or identity-changing supersession;
- legacy fields or a mixture of legacy and canonical records;
- non-deterministic event ordering.

Active verification applies the same per-line and cross-segment invariants
using indexes rebuilt from earlier immutable segments.
