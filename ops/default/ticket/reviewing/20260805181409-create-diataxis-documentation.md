---
id: TKT-20260805181409
title: Create the Diataxis documentation structure
status: reviewing
kind: follow_up
severity: low
priority: normal
scope: docs
tags: diataxis, user-guide, agent-integration
source: chat
reported_at: 2026-08-05T18:14:09Z
ready_at: 2026-08-05T20:06:50Z
started_at: 2026-08-05T20:06:50Z
reviewed_at: 2026-08-05T20:15:44Z
branch: docs/ticket-20260805181409-diataxis-user-guide
pr: 5
commits:
  - 9b45046fd76aadb5d64b5924ae1756ae10907b18
---

# Create the Diataxis documentation structure

## Context

Acta needs durable user documentation that follows the Diataxis structure and
the documentation conventions used by established sibling projects.

## Expected outcome

- Create one `docs/` directory for each Diataxis quadrant: tutorials, how-to
  guides, reference, and explanation.
- Use nested topic directories when they reduce noise.
- Make the user guide the primary documentation path. Give it an index that
  links each file or topic in the progression of using Acta.
- Document how a workspace can instruct an agent when to record an event and
  which values to supply.
- Include concise `AGENTS.md` examples for a selective milestone policy and a
  broad intentional-change policy.
- Use established sibling user guides as structural references without naming
  or copying project-specific content.

## Validation

- All relative Markdown links resolve to files in the branch.
- `git diff --check` passed.
- A current-tree scan found no personal home paths, named sibling repositories,
  private hosts or ports, private-key markers, or recognized token patterns.
- Command, storage, limit, and systemd behavior was checked against the current
  Go implementation and supplied user units.
