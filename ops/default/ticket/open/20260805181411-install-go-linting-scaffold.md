---
id: TKT-20260805181411
title: Install the strict Go linting scaffold
status: open
kind: follow_up
severity: low
priority: normal
scope: infra
tags: go, linting, quality-gate
source: chat
reported_at: 2026-08-05T18:14:11Z
commits: []
---

# Install the strict Go linting scaffold

## Context

Acta needs the standard strict Go quality gate used by established sibling
projects. In particular, the gate should detect dense control-flow blocks such
as adjacent `if` statements that reduce readability.

## Expected outcome

- Add the standard `golangci-lint` v2 configuration for this module.
- Enable the standard formatting, whitespace, control-flow, complexity, and
  correctness checks, including `nlreturn` and `wsl_v5`.
- Add the standard Make quality targets and supporting scripts.
- Validate the current Go code with the installed strict lint command.
