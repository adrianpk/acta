# Choose an Agent Capture Policy

Define the capture policy in each workspace `AGENTS.md`. Different workspaces
can use different granularity. Apply one policy consistently so humans and
agents can interpret the resulting activity record.

## Use a selective milestone policy

Use this policy when the journal should remain a compact index of important
technical landmarks:

```md
After creating a significant commit, immediately run `acta record [flags] [note]` while that commit is still `HEAD`.
Record completed milestones, reusable decisions, accepted specifications, and non-obvious fixes.
Use the note to state why the commit matters when the subject does not make that clear.
Supply `--harness`, `--harness-version`, and `--model` only when those values are known reliably.
Do not record routine, mechanical, generated, formatting-only, punctuation-only, or work-in-progress commits.
Never include credentials, tokens, personal data, or hidden reasoning in the note.
```

This policy makes milestone counts sparse. A missing event does not mean that
no work occurred between recorded milestones.

## Use a broad intentional-change policy

Use this policy when recall is more important than a small journal:

```md
After every completed intentional commit, immediately run `acta record [flags] [note]` while that commit is still `HEAD`.
Include features, fixes, refactors, documentation, tests, maintenance, and small deliberate changes.
Use the note only when the commit subject does not explain the purpose or later retrieval terms.
Supply `--harness`, `--harness-version`, and `--model` only when those values are known reliably.
Do not record generated-only, temporary, work-in-progress, or abandoned commits.
Never include credentials, tokens, personal data, or hidden reasoning in the note.
```

This policy produces denser activity counts. Compare repositories only when
their capture policies are sufficiently similar.

## Choose what the agent supplies

The agent may supply:

- a note of at most 500 Unicode code points;
- `--harness` with the agent harness name;
- `--harness-version` when the exact harness version is known;
- `--model` with a provider-qualified model identifier when known.

`--harness` is required when the agent supplies a harness version or model.
The agent should not guess provenance. Acta derives repository and Git fields
from the current worktree.

Place flags before the optional note:

```sh
acta record \
  --harness example-harness \
  --harness-version 1.2.3 \
  --model provider/example-model \
  "Resolved a reusable service-boundary problem"
```
