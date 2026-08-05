# Command Line Reference

Acta provides `serve`, `record`, `search`, `show`, `report`, `correct`, and
`verify`. All client commands require the local service socket to be available.

## `acta serve`

```text
acta serve
```

Run the per-user service. It adopts the supplied systemd socket when socket
activation is present. Otherwise, it creates the Unix socket directly.

## `acta record`

```text
acta record [--harness <name>] [--harness-version <version>] [--model <id>] [note]
```

Record the current worktree `HEAD`. Run the command inside a Git worktree.
Flags must appear before the optional positional note.

| Option | Meaning |
| --- | --- |
| `--harness <name>` | Record the agent harness name. |
| `--harness-version <version>` | Record the exact harness version. Requires `--harness`. |
| `--model <id>` | Record a provider-qualified model identifier. Requires `--harness`. |
| `note` | Record optional significance text of at most 500 Unicode code points. |

Acta derives the repository root, normalized repository identity and remote,
branch, full commit ID, subject, author and committer times, and current
worktree cleanliness. Recording the same repository identity and commit again
is idempotent.

## `acta search`

```text
acta search [flags]
```

| Option | Meaning |
| --- | --- |
| `--query <text>` | Match normalized text in repository name, subject, note, correction reason, harness, or model. |
| `--repository <value>` | Match an exact repository name or normalized identity, without case sensitivity. |
| `--commit <prefix>` | Match a lowercase hexadecimal full or abbreviated commit. |
| `--harness <name>` | Match a normalized harness name. |
| `--model <id>` | Match a normalized model identifier. |
| `--type <list>` | Match comma-separated `milestone` or `correction` event types. |
| `--recorded-from <time>` | Include events recorded at or after an RFC3339 time. |
| `--recorded-to <time>` | Exclude events recorded at or after an RFC3339 time. |
| `--authored-from <time>` | Include events authored at or after an RFC3339 time. |
| `--authored-to <time>` | Exclude events authored at or after an RFC3339 time. |
| `--committed-from <time>` | Include events committed at or after an RFC3339 time. |
| `--committed-to <time>` | Exclude events committed at or after an RFC3339 time. |
| `--limit <count>` | Return 1 to 1000 events. The default is 20. |
| `--order <value>` | Use `ascending` or `descending`. The default is `descending`. |
| `--json` | Write the structured result. |

Search returns only the effective event in each correction chain. Results can
also be truncated to keep the encoded response within 8 MiB. JSON clients must
check `returned_count` and `truncated`.

## `acta show`

```text
acta show --commit <commit> [--repository <repository>] [--json]
```

Return one effective event. `--commit` accepts a full commit or an unambiguous
lowercase hexadecimal prefix. Use `--repository` when a prefix matches events
from several repositories.

## `acta report`

```text
acta report --from <date> --to <date> [--timezone <name>] [--limit <count>] [--json]
```

Group effective events by local recording date and repository.

| Option | Meaning |
| --- | --- |
| `--from <date>` | Inclusive date in `YYYY-MM-DD` form. |
| `--to <date>` | Exclusive date in `YYYY-MM-DD` form. |
| `--timezone <name>` | IANA timezone name. The process local timezone is the default. |
| `--limit <count>` | Return 1 to 1000 events. The default is 1000. |
| `--json` | Write the structured grouped result. |

JSON clients must check `returned_count` and `truncated`.

## `acta correct`

```text
acta correct --target <event-id> --reason <reason> [--note <note>] [--json]
```

Append a correction to the current effective event.

| Option | Meaning |
| --- | --- |
| `--target <event-id>` | Select the current effective event by ULID. |
| `--reason <reason>` | Explain the correction in at most 500 Unicode code points. |
| `--note <note>` | Replace the significance note. Omit it to preserve the current note; pass an empty value to clear the note. |
| `--json` | Write the structured result. |

The target cannot be missing, stale, or from a different correction-chain
position.

## `acta verify`

```text
acta verify [--full | --scheduled] [--json]
```

| Option | Meaning |
| --- | --- |
| `--full` | Verify every canonical monthly segment. |
| `--scheduled` | Verify the active segment daily and the complete journal on the first scheduled run of each local month. |
| `--json` | Write the structured result. |

Without a scope option, Acta verifies the active monthly segment. `--full` and
`--scheduled` are mutually exclusive. A corruption diagnostic identifies the
segment, one-based line, zero-based byte offset, and reason.
