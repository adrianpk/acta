# Agent Operating Contract

Read this document before starting work in this workspace.

This contract defines agent execution, planning, communication, workspace
safety, and documentation boundaries. It is not a product or runtime
specification.

## 1. Execution

- Execute only after a clear implementation request.
- Treat questions, exploration, and design discussion as alignment, not as
  authorization to change files or delivery state.
- For high-impact changes, explain the approach and wait for explicit approval.
- Stay inside the requested scope. Record adjacent ideas instead of implementing
  them.
- After a defined workflow unit is authorized, continue through implementation,
  validation, versioning, push, and handoff unless a policy gate or blocker
  requires user input.
- Do not start, stop, or restart long-running development services unless the
  user requests it or the authorized validation requires it.
- Complete the requested scope, report the outcome, and stop.

## 2. Planning and Operational Artifacts

- Create plans only when the user requests planning or a persistent task list.
- Use `ops/<namespace>/` for operational artifacts. Resolve the namespace from
  explicit user input, then `OPS_USER`, then `default`.
- Use tickets for bounded independent work.
- Use specs, plans, trackers, slices, and reports for non-trivial coordinated
  work.
- Do not implement a formal slice until its plan and tracker exist, are approved,
  and are committed to `dev`.
- Keep operational delivery records out of `docs/` unless they become durable
  product documentation.

## 3. Communication

- Use the user's language for conversation. Use English for repository
  artifacts, code, comments, and configuration.
- Use direct, neutral language. Keep updates factual and concise.
- Explain material reasoning before proposing a significant change.
- Report exact validation results, failures, and remaining risks.
- Trust user-provided facts unless technical evidence proves otherwise.
- For English technical communication, apply the principles of ASD-STE100
  Simplified Technical English when practical. Prefer short sentences, one
  topic per sentence, active voice, and consistent terms. Keep necessary
  project-specific technical terms. This is a writing goal, not a claim of
  certified or full ASD-STE100 compliance.

## 4. State Alignment

- For a status question, confusion, or state challenge, report the verified
  state without changing files, branches, tickets, trackers, or pull requests.
- When the user asks to start, continue, resume, or proceed with an already
  defined workflow, verify the required state and continue in the same turn.
- Questions such as "shall we start?", "can we continue?", or "comenzamos?"
  count as execution authorization when the next workflow unit is already
  defined.
- Do not stop after routine state verification. Stop only for a mismatch,
  destructive action, missing authority, material user decision, or blocker.
- Once a slice is authorized, continue through its implementation, validation,
  commits, push, and pull-request handoff without intermediate approval pauses.
- When the user reports that a slice pull request was merged, verify the merge,
  close the delivered slice on `dev`, and immediately execute the next slice
  already defined by the same committed plan and tracker. Report the transition
  without requesting another confirmation.
- A merge confirmation or `OK` does not authorize a new delivery set, plan,
  branch, or unrelated work.
- Before closing, activating, or starting a slice, verify the current branch,
  remote `dev`, tracker state, and pull-request or merge state. Report mismatches
  before changing state.

## 5. Workspace and Git Safety

- When this workspace uses Git, `origin` is the canonical remote unless local
  instructions say otherwise.
- When a `mirror` remote exists, treat it as a publication remote and push only
  the refs allowed by the local versioning rules.
- Never perform destructive Git operations without explicit approval.
- Preserve unrelated user changes and report any overlap with the requested
  work.

## 6. Language and Documentation

- Keep code, comments, configuration, and repository documentation in English.
- Keep product and runtime specifications in `docs/` or the approved
  `ops/<namespace>/spec/` location.
- Keep working rules in `AGENTS.md` and `.agents/`.
- Keep delivery state in `ops/<namespace>/`.

## 7. Scope Boundary

- Keep this file focused on shared working rules.
- Keep project and stack details in the referenced `.agents/` documents.

## 8. Referenced Agent Docs

- `.agents/DOCUMENTING.md` — documentation model and placement rules.
- `.agents/CONVENTIONS.md` — coding standards and codebase boundaries.
- `.agents/ARCHITECTURE.md` — workspace structure and ownership map.
- `.agents/LIFECYCLE.md` — delivery lifecycle and handoff rules.
- `.agents/VERSIONING.md` — commit, branch, remote, and merge rules.
- `.agents/SLICES.md` — task decomposition and delivery strategy.
- `.agents/ROUTING.md` — request classification and execution modes.
- `.agents/PLAYBOOKS.md` — operational procedures.
- `.agents/REPORTING.md` — persistent implementation reporting rules.
- `.agents/LINTING.md` — validation and quality-gate rules.
- `.agents/INVARIANTS.md` — non-negotiable collaboration invariants.

## Acta Activity Record

- After creating a commit for completed work that is not hypertrivial, run
  `acta record [flags] [note]` immediately while that commit is still `HEAD`.
- Record features, fixes, refactors, meaningful documentation, tests,
  maintenance, decisions, specifications, plans, and operational milestones.
  During initial adoption, record the commit when uncertain.
- Omit only abandoned or work-in-progress commits, generated-only changes,
  punctuation-only changes, isolated formatting, or similarly trivial edits.
- Add a concise note when the commit subject lacks useful retrieval terms or
  does not explain why the change matters. Supply provenance flags only when
  their values are known reliably.
- Never record credentials, tokens, personal data, or hidden reasoning.
