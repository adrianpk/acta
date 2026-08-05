# Work Request Routing

This document defines how incoming requests are classified.

## Execution Modes

- exploration
- planning
- implementation
- investigation
- maintenance or meta-work

## Rules

- Exploration covers questions, possibilities, design discussion, and tradeoff
  analysis. It does not create persistent artifacts or delivery state.
- Planning occurs only when requested. A persistent plan, tracker, ticket, or
  specification requires explicit authorization.
- Implementation starts only with explicit instruction and stays within the
  agreed scope.
- A request to start, continue, resume, or proceed authorizes execution when the
  next workflow unit is already defined.
- Do not turn exploration, a plan, a prior merge, or a status acknowledgment
  into open-ended autonomous work.
- Ambiguous requests require clarification before code or state changes when a
  reasonable assumption could materially change the result.
- During exploration, work inside the direction the user is evaluating. Map how
  the idea could fit and name concrete tradeoffs. Push back only when there is a
  specific technical, safety, scope, or product risk.
- Use `.agents/playbooks/work-modes.md` whenever work may move from exploration
  into delivery.

## Delivery Mode

- An approved slice includes implementation, validation, versioning, push, and
  its pull-request handoff.
- A change outside an approved slice uses direct `dev` delivery unless the user
  explicitly requests a branch and pull request.
- Do not infer branch or pull-request delivery from task size or risk.

## Non-goal

This document does not define product or runtime behavior.
