# Acta

Acta preserves durable technical landmarks for fast, agent-driven development.

Before AI-assisted development, a developer often remembered most of the work
completed during a week because they handled each implementation detail. An
agent can now work quickly across many projects while the human acts more like
a lead: setting direction, reviewing outcomes, and making decisions. With
dozens of agents operating in different repositories, useful technical history
is produced faster than the human can retain or rediscover it.

Acta closes that gap. After a significant Git commit, the working agent records
a concise milestone in a local append-only journal. A later agent can locate
when and where the work happened, then inspect the recorded commit to recover
transitively how it was done.

This creates one searchable activity index across repositories without moving
the underlying code or documentation out of Git. It gives the human an
evidence-backed account of which significant milestones landed, where, and
when. It gives agents a short route from each milestone to the exact commit and
its implementation context.

Acta is an evidence index, not a transcript store or an automatic log of every
commit.

## What Acta makes possible

An agent integration can translate natural-language questions into bounded
Acta queries and use the returned repository and commit as evidence.

### Locate a previous implementation

> When did we implement PII redaction for stored secrets, and in which project?

Acta can locate a matching milestone, including its repository, timestamp,
commit, subject, and significance note. The agent can then inspect that commit
instead of searching every repository or conversation.

### Reuse a pattern from another project

> Last week, project X introduced a spec, plan, and tracker workflow. Reuse that
> pattern here.

Specs, plans, architecture decisions, and documentation can be milestones when
they are committed to Git. Acta finds the relevant granular commit. The agent
loads the original artifacts from Git and adapts the pattern to the current
project without copying project-specific assumptions.

### Trace an observed delivery cycle

> Show the progression from the initial versioned idea through specification,
> planning, implementation, and integration.

When each significant phase was recorded, Acta can order those milestone
commits and expose the observed interval between them. It cannot determine when
an unrecorded idea began, how many hours the work required, or whether two
commits are causally related when their evidence does not say so.

### Understand relative activity

> Which repositories received the most significant work this week, and how
> concentrated was it?

An agent can derive milestone counts, active days, distribution by repository,
work type, harness, or model, and concentrations of recorded activity. These
are transparent measures of recorded milestones. They are not estimates of
hours, effort, difficulty, or individual productivity.

Other uses include finding a prior technical decision, comparing similar
problems across harnesses or models, identifying projects resumed after a quiet
period, and producing bounded weekly summaries across repositories.

## Current capabilities

- Records significant Git milestones with repository, commit, timestamp, note,
  and optional agent provenance.
- Stores self-contained events in a local append-only journal.
- Searches by text, repository, commit, time, harness, or model.
- Shows the effective record for a commit and produces bounded activity
  reports.
- Corrects records without rewriting history.
- Detects journal corruption and known interrupted writes.
- Avoids duplicate records for the same repository and commit.

## How it fits into the workflow

```text
human direction
    -> agent implements and commits a significant milestone
    -> agent records the commit in Acta
    -> a later agent runs a bounded query
    -> Acta returns repository and commit evidence
    -> the agent inspects Git for the implementation details
```

Acta stores the landmark and the route back to the source. Git remains the
source for code, specs, plans, diffs, and detailed history.

## Agent integration

Each workspace defines its own capture policy. One workspace can record only
important milestones, while another can record every completed intentional
change. The [Acta user guide](docs/tutorials/user-guide/index.md) covers the
first setup. [Choose an agent capture
policy](docs/how-to/choose-an-agent-capture-policy.md) provides concise
`AGENTS.md` rules for both approaches.

## Evidence boundaries

Acta records milestone metadata. It does not store prompts, conversations,
hidden reasoning, diffs, patches, file contents, environment variables, or
terminal history. It does not contact an external service for search or
reporting.

Acta reports stored facts and transparent derivations. It does not invent
missing activity, elapsed effort, working hours, causal relationships, or
productivity measurements.
