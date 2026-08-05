# Delivery Lifecycle Contract

This document defines lifecycle rules for work in this workspace.

## Work Phases

1. Alignment
2. Planning when requested
3. Implementation after approval
4. Validation
5. Versioning and push within the authorized workflow
6. Report and handoff

## Entry Gates

- Alignment starts when task context is provided.
- Persistent plans, trackers, specs, and tickets require an explicit request.
- Implementation starts only after the task or plan is approved.
- Production deployment, billing mutations, credential changes, and final
  merges require explicit authorization.

## Exit Gates

A work unit is complete only when:

- the requested scope is implemented;
- relevant checks and exact failures are reported;
- versioned artifacts reflect the actual state;
- the working tree and branch state are explicit;
- remaining risks and external prerequisites are recorded.

## Blocked State

When blocked:

- capture the failing command and concise error context;
- exhaust safe in-scope alternatives;
- record the blocker in the active tracker or ticket when one exists;
- stop before speculative, destructive, or out-of-scope changes.

## Handoff Artifacts

- Durable plans live in `ops/<namespace>/plan/`.
- Execution state lives in `ops/<namespace>/tracker/`.
- Slice reports live in `ops/<namespace>/report/slices/`.
- End-to-end cycle reports live in `ops/<namespace>/report/`.
- Temporary command output belongs under an ignored temporary directory, not in
  versioned operational artifacts.
