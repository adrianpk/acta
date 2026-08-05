---
id: TKT-20260805181408
title: Write a clear project README
status: in_progress
kind: follow_up
severity: low
priority: normal
scope: docs
tags: readme, onboarding, agent-integration
source: chat
reported_at: 2026-08-05T18:14:08Z
ready_at: 2026-08-05T20:00:15Z
started_at: 2026-08-05T20:00:15Z
commits: []
---

# Write a clear project README

## Context

Acta needs a concise entry point that explains the product, its purpose, and
its normal usage without requiring readers to inspect implementation files.

AI-assisted development moves the human developer toward directing and
reviewing fast work across several projects. The README must explain how Acta
preserves technical landmarks when the human no longer handles every low-level
implementation detail.

## Expected outcome

- Explain what Acta records and how it fits into a development workflow.
- Explain how Acta locates when and where significant work happened and uses
  the recorded commit to recover how it was done from Git.
- Present current capabilities without becoming a complete user manual.
- Include at least four concrete uses: locating prior implementation, reusing
  a versioned pattern, summarizing relative activity intensity, and tracing an
  observed idea-to-integration cycle through recorded milestone commits.
- Distinguish stored evidence, transparent derivations, and unavailable facts.
- Include a short agent-integration example that shows how an `AGENTS.md` file
  can invoke Acta, which values it should pass, and which events it should omit.
- Use a trivial punctuation-only commit as an example of an event that should
  not be recorded.
- Link to the detailed user guide when that documentation exists.
