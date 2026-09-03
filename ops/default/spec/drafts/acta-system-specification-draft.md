# Acta System Specification

Status: superseded

Superseded by: `ops/default/spec/stable/typed-activity-journal.md`

Scope: Product, service, and behavioral contract for a low-ceremony chronicle
manager shared by independent agent sessions. This specification does not
authorize implementation.

Product description: Cross-session agent chronicle manager.

Executable: `acta`

Service unit: `acta.service`

Go module: `forge.adrianpk.com/adrianpk/acta`

Model companion: [Acta System Model](acta-system-specification-draft-model.md)

Implementation language: Go

## 1. Purpose

Acta manages small, durable records of significant work completed with coding
agents across repositories and agent runtimes.

An agent records a milestone after creating a significant commit. The journal
links a short statement of significance to exact Git evidence. Later, any agent
can locate related work in another project, and the user can produce daily or
weekly activity summaries.

Acta is a shared chronicle manager, not shared live context. Ordinary agent
sessions must not load its journal by default.

## 2. Core outcome

The system must make this workflow inexpensive:

1. An agent completes a significant commit.
2. The agent decides whether the commit meets the recording policy.
3. The agent runs one short command.
4. The CLI derives repository and Git metadata and submits one record request.
5. The Acta service validates the request and appends one durable event.
6. Work continues without loading unrelated journal history.

The normal recording command is:

```sh
acta record "Reusable protocol normalization across agent-enabled projects"
```

When the commit subject already explains its significance, the agent may run:

```sh
acta record
```

## 3. Actors

- **User:** owns the journal and asks for records, searches, or summaries.
- **Working agent:** decides whether a completed commit is significant and
  invokes Acta. Its integration may supply optional harness, harness version,
  and model provenance without user input.
- **CLI client:** derives local repository facts and sends bounded requests to
  the Acta service.
- **Acta service:** validates requests, owns canonical storage, serializes
  writes, and executes queries.
- **Retrieval skill:** translates a user request into bounded CLI queries and
  presents the result.

Acta must remain independent of Codex, Pi, Herdr, terminal emulators, and other
agent hosts.

### 3.1 Service boundary

Acta is installed as a per-user local service named `acta.service`. A companion
`acta.socket` unit provides socket activation. The `acta` executable is the
command-line client and also contains the service entrypoint.

The service is the sole canonical journal writer. Clients must not fall back to
editing canonical journal segments directly when the service is unavailable.

Client-to-service transport uses a Unix domain stream socket at:

```text
$XDG_RUNTIME_DIR/acta/acta.sock
```

The runtime directory uses mode `0700`. The socket uses mode `0600`.

Each connection contains one bounded JSON request and receives one JSON
response before it closes. Requests identify the `acta/v1` protocol, one
operation, and its parameters. Responses contain either a result or a stable
error code and message. Normal requests are limited to 64 KiB. Normal responses
are limited to 8 MiB. Queries stop at 1,000 events or the response byte limit,
whichever occurs first, and report truncation. A future full export must use a
separate streaming operation.

Request and response envelopes have these shapes:

```json
{"protocol":"acta/v1","operation":"record","params":{}}
```

```json
{"protocol":"acta/v1","ok":true,"result":{}}
```

```json
{"protocol":"acta/v1","ok":false,"error":{"code":"invalid_request","message":"request exceeds the allowed size"}}
```

The protocol does not use request files, TCP, HTTP, gRPC, or an external
network service. Connecting to `acta.socket` starts `acta.service` when it is
not already running. The service may remain active after the request completes.

The initial service error codes are:

```text
invalid_request
unsupported_protocol
unknown_operation
not_found
invalid_event
journal_corrupt
journal_unavailable
internal_error
```

An idempotent duplicate record is a successful response with the existing
event ID and `recorded` set to `false`. It is not an error response.

### 3.2 Chronicle boundary

A chronicle is a logical, queryable view of milestone events. Acta may present
chronicles by project, agent session, topic, or time interval without requiring
one physical file per view.

The initial product does not persist stable agent-session identity or one stored
book per session. Session-oriented chronicles are deferred query projections
and require reliable provenance in a later schema version. Cross-session
project and time queries do not depend on session identity.

## 4. Significance policy

A commit is significant when at least one of these conditions is true:

- it completes a recognizable milestone;
- it records a decision or pattern likely to be reused;
- it resolves a difficult or non-obvious problem;
- it materially changes behavior, architecture, or an operating workflow;
- it publishes or integrates a relevant delivery.

The following work is not significant by default:

- formatting, spelling, or cosmetic cleanup;
- routine dependency or metadata maintenance;
- mechanical regeneration;
- work-in-progress checkpoints;
- trivial refactors with no reusable conclusion;
- ordinary merge commits that add no distinct milestone information.

The working agent owns this semantic decision. The CLI must not attempt to
score every commit or record commits automatically.

## 5. Agent instruction contract

A workspace may include the following instruction in `AGENTS.md`:

```md
After creating a significant commit, run `acta record` to add it to the shared activity journal.
A commit is significant when it completes a milestone, records a reusable decision, or resolves a non-obvious problem.
Do not record routine, mechanical, formatting, or work-in-progress commits.
```

A workspace may add one sentence requiring a short note when its commit
subjects do not explain significance. It must not add a longer recording
ceremony or require the journal to be read before ordinary work.

## 6. Recording contract

### 6.1 Preconditions

`record` must run inside a Git worktree with at least one commit. By default it
records `HEAD` after the commit exists.

The CLI must not create, amend, stage, commit, reset, or otherwise modify Git
state.

### 6.2 Inputs

The only normal user-supplied input is an optional single-line note explaining
why the commit matters.

Repository name, root path, remote identity, branch, commit hash, commit
subject, author time, committer time, and worktree state must be derived
automatically. The service assigns the separate journal recording time.

An agent integration may automatically supply an optional execution harness,
harness version, and model identifier when it knows those values reliably. A
manual `acta record` does not require them. Acta must not infer execution
provenance by inspecting parent processes.

### 6.3 Output

On success, the service appends one milestone event. The CLI prints a concise
confirmation containing the event ID, repository, and abbreviated commit.

Recording the same repository identity and commit again must be idempotent. It
must not create a second milestone event. The CLI must return the existing
event ID and report that the commit was already recorded.

### 6.4 Failure behavior

The CLI must fail without changing either the repository or the existing
journal when:

- the current directory is not inside a Git worktree;
- the worktree has no commit;
- required Git metadata cannot be read;
- the note is invalid;
- the Acta service is unavailable or rejects the request;
- the journal cannot be locked, validated, or appended safely.

A recording failure must not invalidate or reverse the completed commit. The
agent reports the failure briefly as part of its normal task handoff.

## 7. Journal boundaries

The journal stores milestone metadata only. It must not store:

- prompts or conversation transcripts;
- hidden agent reasoning;
- diffs, patches, or complete file contents;
- environment variables or credentials;
- terminal scrollback;
- a general event stream of every agent action;
- an automatic copy of every Git commit.

The journal is logically append-only. Published events must not be edited or
deleted through the CLI. A correction is a new event that references an older
event and repeats its effective execution, repository, Git, and note context so
the correction remains independently searchable and understandable.

## 8. Storage contract

Acta configuration belongs under:

```text
~/.config/acta/
```

Acta persistent and transactional files belong under:

```text
~/.acta/
```

The canonical journal root is:

```text
~/.acta/events/
```

These locations are outside `~/.codex` and `~/.pi` because Acta belongs to the
user and is shared by multiple agent runtimes.

The journal is segmented by the UTC year and month of `recorded_at`:

```text
events/
    2026/
        08.jsonl
        09.jsonl
    2027/
        01.jsonl
```

The ordered set of monthly segments is the canonical journal. Each event is one
complete compact JSON object followed by `LF`. The segment selected by the
service recording timestamp is the only segment open for append. After the
service begins writing a later month, earlier monthly segments are immutable.

The initial product is single-user and single-host. Cross-host replication,
backup, and conflict resolution are outside the initial scope.

The service must serialize concurrent clients, validate the active segment
tail, append a complete UTF-8 JSON object followed by one newline, and sync the
append before confirming success. A global derived index supports event lookup
and milestone deduplication across all segments. The index is not canonical and
must be rebuildable from the monthly JSONL files.

### 8.1 Interrupted append recovery

Before appending, the service writes one pending operation to:

```text
~/.acta/pending.json
```

The pending operation contains the target segment, its original byte length,
and the complete event line. The service syncs the pending operation before it
appends and syncs the segment before it removes the pending operation.
The state directory and pending file are accessible only to the owning user.

After an interrupted append, the service may recover only these known states:

- append the pending event when the segment still has its original length;
- clear the pending operation when the complete event is already present;
- remove a partial terminal prefix that matches the pending event, restore the
  original length, and append the complete event again.

Any other mismatch is corruption. The service must preserve the files, reject
new writes to the affected active segment, and report the exact segment, line,
and byte offset when available. It must not perform general automatic repair.

### 8.2 Integrity verification

The CLI provides:

```sh
acta verify
acta verify --full
```

`acta verify` validates the active monthly segment. `acta verify --full`
validates every canonical segment and the global invariants across them.

An `acta-verify.timer` per-user timer runs the active-segment verification each
day at 03:00 in the host's local timezone. The timer uses persistent scheduling
so a missed run occurs after the host becomes available. At the start of a new
month, the scheduled verification also runs a full-journal validation.

Verification checks UTF-8, line framing, JSON decoding, schema versions, field
rules, event ID uniqueness, milestone key uniqueness, correction targets,
event order, and the correspondence between `recorded_at` and the segment path.
It is read-only and does not repair corruption.

The initial product does not calculate or periodically scan historical file
checksums. Full semantic verification is available explicitly and at the
monthly boundary without adding a separate checksum lifecycle.

## 9. Retrieval contract

Reading is explicit and bounded. The journal must not be injected into an
agent's context at session start or before every task.

The CLI must support these conceptual operations:

```sh
acta show --commit <commit> --json
acta search --query <text> --limit <count> --json
acta search --query <text> --harness <name> --model <identifier> --limit <count> --json
acta report --from <date> --to <date> --json
```

Queries may filter by repository, recording time, Git author time, Git
committer time, event type, execution harness, model, or text. Every search must
use an explicit or safe default result limit. Machine-readable output must
contain only the matching records and query metadata. Effective correction
results repeat their repository, Git, execution, and note context and require no
event join for ordinary inspection.

The initial text search may be deterministic rather than semantic. A retrieval
skill may use the bounded result set to rank, explain, or summarize matches.

## 10. Retrieval skill contract

A dedicated skill may provide user-facing workflows such as:

- find where a similar problem was solved;
- locate the exact repository and commit for a prior decision;
- summarize activity for a day, week, or explicit interval;
- group milestones by project or topic;
- compare related milestones across projects;
- compare similar issues across agent harnesses or models;
- distinguish Git author time, commit time, and Acta recording time.

The skill must call the CLI with a bounded query before reading event content.
It must not load the entire journal unless the user explicitly requests a full
export or full-journal audit.

Reports are generated views. They do not become journal events unless the user
explicitly records or saves them through a separate workflow.

## 11. Context-budget invariants

- Recording must require no journal read beyond the service's integrity and
  duplicate checks.
- Ordinary work must not receive journal contents automatically.
- Retrieval must select by query, repository, or time range before agent
  interpretation.
- CLI output intended for agents must support JSON and explicit limits.
- A journal failure must not cause the agent to re-read unrelated history.

## 12. Security and privacy

- Notes must be short plain text and must not contain secrets.
- Execution provenance must contain identifier strings and must not contain
  prompts, environment values, credentials, or session data.
- Git remotes must be normalized without embedding credentials.
- URL user information, tokens, and credential-bearing query parameters must be
  removed before persistence.
- The service must create journal files with permissions limited to the owning
  user.
- Search and report commands must not contact an external network service in
  the initial product.

## 13. Non-goals

The initial product does not:

- share live conversational context between agents;
- coordinate tasks or allocate work;
- replace Git history, repository documentation, specifications, or tickets;
- record every file edit, command, or commit;
- infer significance without an agent decision;
- monitor agent processes or terminal sessions;
- provide a web interface or remote network API;
- provide semantic embeddings or a vector database;
- synchronize journals across machines;
- modify project files while recording an event.

## 14. Accepted decisions

- The product and executable name is Acta.
- The concise product description is `Cross-session agent chronicle manager.`
- Acta is implemented in Go as module `forge.adrianpk.com/adrianpk/acta`.
- Acta includes a persistent local service named `acta.service` and a
  command-line client named `acta`.
- Acta runs as a per-user service with `acta.socket` activation.
- The client and service exchange bounded JSON requests and responses over a
  protected Unix domain socket.
- The service is the sole canonical journal writer.
- Recording happens at the natural boundary after a significant commit.
- The agent decides significance using a short shared policy.
- One optional note is the only normal manual input.
- Git and repository facts are derived by the CLI.
- Git author and committer times are stored separately from the Acta recording
  time.
- Optional execution provenance distinguishes the agent harness, harness
  version, and model associated with the work.
- Acta service name and version are not repeated in each event.
- Corrections repeat the effective milestone context and remain independently
  searchable while preserving their target link.
- Configuration uses `~/.config/acta/`. Persistent files use `~/.acta/`.
- Storage is global to the user, append-only, segmented by UTC month, and
  outside agent-specific homes.
- One pending operation provides bounded recovery from an interrupted append.
- Daily verification covers the active segment. Monthly and explicit full
  verification cover all canonical segments.
- The initial product does not maintain historical file checksums.
- Stable session identity and per-session stored books are not part of the
  initial product.
- Retrieval is explicit, filtered, and bounded.
- A separate skill provides searches and periodic summaries.
- The journal stores evidence pointers and significance, not complete work
  context.

## 15. Open questions

- Whether a later version should support explicit backfill of a non-`HEAD`
  commit.
- How a future multi-host journal would synchronize while preserving append-only
  semantics.
- Whether generated reports need a separate durable artifact format.

These questions do not block a planning session. A plan must resolve each
question before implementing the affected service or persistence boundary.

## 16. Acceptance criteria

The specification is satisfied when an implementation can demonstrate that:

1. Connecting to the protected per-user Unix socket activates the Acta service
   and the journal remains available across a service restart.
2. An agent records a significant `HEAD` commit with one CLI command.
3. The CLI derives the repository, remote identity, branch, full commit hash,
   commit subject, author time, committer time, and worktree state.
4. An optional single-line note is stored without requiring other parameters.
5. Repeating `record` for the same repository and commit creates no duplicate.
6. Concurrent clients cannot corrupt or interleave monthly JSONL events.
7. The CLI does not edit canonical storage when the service is unavailable.
8. Published events cannot be modified through ordinary CLI operations.
9. A correction is represented by a new linked, self-contained event that can
   be understood and searched without reading its target.
10. Recording does not modify Git state or project files.
11. Search and report operations return bounded machine-readable results.
12. A weekly report can be produced from an explicit date range without loading
    unrelated journal history.
13. Codex and Pi can use the same executable and journal without runtime-specific
    storage.
14. The journal contains no prompts, transcripts, diffs, credentials, or hidden
    reasoning.
15. An interrupted append is recovered idempotently from the pending operation
    without rewriting a previously published event.
16. Daily verification detects invalid active-segment framing or event data and
    reports the exact location when available.
17. Full verification detects invalid event data and global identity or
    correction violations across monthly segments.
18. A query can filter comparable milestones by execution harness or model and
    can distinguish author, committer, and recording times.
