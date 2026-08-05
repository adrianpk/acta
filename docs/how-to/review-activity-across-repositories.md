# Review Activity Across Repositories

Use search for discovery, show for one exact event, and report for a bounded
calendar interval.

## Locate prior work

Search commit subjects, notes, correction reasons, repository names, harnesses,
and models:

```sh
acta search --query "redaction" --limit 20
```

Narrow the result when necessary:

```sh
acta search \
  --repository example-project \
  --recorded-from 2026-08-01T00:00:00Z \
  --recorded-to 2026-09-01T00:00:00Z \
  --order ascending \
  --json
```

## Resolve the exact event

```sh
acta show --commit <full-or-unambiguous-prefix> --json
```

Add `--repository` when the same abbreviated commit matches more than one
repository. Use the returned repository root and commit to inspect the source
with Git.

## Summarize a period

```sh
acta report \
  --from 2026-08-01 \
  --to 2026-08-08 \
  --timezone Europe/Paris \
  --json
```

Reports group effective events by local recording date and repository. They
describe recorded milestones, not hours worked or productivity. Check the
`truncated` field before treating JSON output as complete.
