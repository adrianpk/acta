# Acta

Acta is a local, append-only journal for durable professional activity.

Work does not always produce a Git commit. Configuring Grafana, building a
pipeline, reviewing a pull request, investigating an incident, and preparing a
timesheet can all be meaningful activity. Acta records those outcomes as peer
event types instead of treating Git as the universal shape of work.

The initial event types are:

- `commit` for a completed Git commit;
- `pull_request` for an opened, reviewed, merged, or closed pull request;
- `worklog` for declared work with an explicit duration;
- `note` for user-directed information that does not fit a more specific type.

Each event has an occurrence time, a recording time, a concise summary, and
optional scope, references, and provenance. Published events are immutable.
Corrections append complete replacements and preserve the original history.

## What Acta makes possible

Acta provides one searchable activity index across repositories, projects,
systems, clients, and other named contexts. It can answer questions such as:

- What changed in production observability this week?
- Which pull requests were reviewed or merged for a project?
- How much explicitly recorded work belongs on a timesheet?
- Where is the commit or external page related to a past decision?

Search and reports return compact event projections. `show` returns one
complete canonical record when its properties, references, provenance, or
correction chain are needed.

## Record activity

Record a commit from its Git worktree:

```sh
acta record commit "Established the typed activity model"
```

Record a pull-request action without contacting the forge:

```sh
acta record pull-request \
  --forge https://forge.example.test \
  --repository team/example \
  --number 42 \
  --action merged \
  "Integrated typed activity retrieval"
```

Record work that has no Git artifact:

```sh
acta record worklog \
  --duration 240 \
  --scope-type system \
  --scope-name "Production observability" \
  --scope-identity system:production-observability \
  "Configured Grafana alert routing"
```

Record user-directed fallback information:

```sh
acta record note "The pipeline cutover is scheduled for Friday"
```

`acta record` always requires an explicit type. A commit is not a default or a
prerequisite for another event.

## Retrieve and report

Search across event types and scopes:

```sh
acta search \
  --query "alert routing" \
  --scope system:production-observability \
  --occurred-from 2026-09-01T00:00:00Z
```

Inspect one complete effective event:

```sh
acta show --id <event-id> --json
```

Build a timesheet-oriented report using local occurrence dates:

```sh
acta report \
  --from 2026-09-01 \
  --to 2026-09-08 \
  --timezone Europe/Warsaw \
  --type worklog
```

Only stored `worklog.duration_minutes` values contribute to duration totals.
Acta never derives hours from commits, event counts, timestamps, or elapsed
time.

## Agent integration

Each workspace can define a commit capture policy in `AGENTS.md`. A human can
also instruct an agent to record a `worklog` or `note` for activity outside
Git. User-directed note capture is limited to the information the user asked to
preserve; it does not authorize automatic conversation capture.

The [Acta user guide](docs/tutorials/user-guide/index.md) covers initial setup
and the first typed workflow. [Choose an agent capture
policy](docs/how-to/choose-an-agent-capture-policy.md) provides reusable policy
text.

## Boundaries

Acta stores local facts and optional routes to related artifacts. It does not
monitor applications, infer activity, submit timesheets, mutate referenced
systems, or contact external services for recording and retrieval. It does not
store prompts, conversations, hidden reasoning, diffs, patches, file contents,
environment variables, or terminal history automatically.

Event counts and explicit worklog duration describe recorded data. They are
not measures of effort, difficulty, causality, or productivity.
