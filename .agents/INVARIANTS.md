# Collaboration Invariants

This document defines non-negotiable execution invariants.

## Scope Integrity

- Execute only approved scope.
- Record out-of-scope ideas instead of implementing them.
- Do not infer authorization from task size, technical value, or prior work.

## Approval Gates

- Questions and exploration are not implementation requests.
- High-impact changes require explicit approval.
- Formal slice execution requires an approved, committed plan and tracker.
- Production deployment and external billing, credential, or access mutations
  always require explicit approval.

## Safe Delivery

- Prefer incremental, reversible changes.
- Keep the repository in a working state between delivery units.
- Do not use destructive Git operations without explicit approval.
- Keep secrets and local agent configuration out of version control.
- Preserve unrelated user changes.

## Quality and Truthfulness

- Do not mark work complete with failing required checks.
- Do not describe focused checks as the full validation suite.
- Distinguish local implementation from deployed or externally verified state.
- Record checks that could not run and the reason.
- Keep repository, tracker, commit, and pull-request state consistent.

## Documentation Separation

- Product behavior belongs in `docs/` or approved specifications.
- Working rules belong in `AGENTS.md` and `.agents/`.
- Delivery state belongs in `ops/<namespace>/`.
