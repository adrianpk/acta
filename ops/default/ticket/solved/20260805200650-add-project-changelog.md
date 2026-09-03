---
id: TKT-20260805200650
title: Add the project changelog
status: solved
kind: follow_up
severity: low
priority: normal
scope: docs
tags: changelog, releases
source: chat
reported_at: 2026-08-05T20:06:50Z
started_at: 2026-09-03T13:53:59Z
reviewed_at: 2026-09-03T13:54:22Z
closed_at: 2026-09-03T13:54:22Z
resolution: fixed
commits:
  - 534f528ff668cb8fb5858614a531a748da4328d2
---

# Add the project changelog

## Context

Acta needs a changelog that records notable project changes without coupling
release history to the user-guide structure.

## Expected outcome

- Add a changelog consistent with the conventions used by established sibling
  projects.
- Keep changelog maintenance independent from the Diataxis documentation
  structure.

## Validation

- `make docs-check` passed.
- `git diff --cached --check` passed.
