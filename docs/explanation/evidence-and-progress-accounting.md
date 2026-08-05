# Evidence and Progress Accounting

Agent-assisted development changes who retains implementation context. A human
can direct and review many concurrent agents without handling every low-level
change. Across many repositories, completed work can accumulate faster than the
human can remember it or a later agent can rediscover it.

Acta records a compact landmark after a commit. The landmark answers where and
when meaningful work was recorded. Its commit identifies the immutable Git
evidence that a later agent can inspect to recover how the result was produced.

## Two consumers

For an agent, Acta is a retrieval index. A natural-language request can become
a bounded search. The result provides the repository and commit needed to load
the original diff, specification, plan, or implementation from Git.

For a human, Acta is an evidence-backed account of progress across
repositories. Reports show which recorded milestones landed, where, and when.
They can help review a week, locate concentrated activity, or see that a quiet
project became active again.

## Evidence levels

Acta supports four distinct levels of interpretation:

1. Stored facts include the recording time, repository, commit, subject, note,
   Git timestamps, worktree state, and supplied execution provenance.
2. Transparent derivations include milestone counts, active days, repository
   distributions, and intervals between recorded events.
3. Approximate inferences require context outside Acta. For example, an agent
   can inspect several commits and infer that they form one delivery cycle.
4. Unavailable facts include unrecorded work, hours, effort, difficulty,
   productivity, hidden reasoning, and the start time of an unversioned idea.

An idea-to-integration cycle is observable only when its meaningful stages were
committed and recorded. Event order alone does not prove causal dependence.

## Capture policy affects interpretation

A selective workspace records a sparse set of major landmarks. A broad
workspace records most completed intentional commits. Both are valid, but
their event counts do not mean the same thing. Comparisons should state the
capture policy and use counts as activity indicators, not performance scores.
