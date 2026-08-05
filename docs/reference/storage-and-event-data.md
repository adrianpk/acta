# Storage and Event Data Reference

Acta stores per-user data locally and uses fixed paths in the current release.
There is no user-editable configuration file yet.

## Paths

| Path | Purpose |
| --- | --- |
| `$HOME/.acta/events/YYYY/MM.jsonl` | Canonical monthly UTC journal segment. |
| `$HOME/.acta/pending.json` | Recovery record for one interrupted append. |
| `$HOME/.acta/verification.json` | Month of the last successful scheduled full verification. |
| `$XDG_RUNTIME_DIR/acta/acta.sock` | Per-user Unix service socket. |

Acta creates data directories with mode `0700` and owned files with mode
`0600`. The supplied systemd units use `UMask=0077`. The service is the sole
canonical journal writer.

## Milestone event fields

Each JSONL line contains one complete `acta/v1` event.

| Field | Meaning |
| --- | --- |
| `schema_version` | Durable event schema, currently `acta/v1`. |
| `id` | Time-sortable ULID assigned by the service. |
| `type` | `milestone` or `correction`. |
| `recorded_at` | Service recording time in UTC. |
| `execution` | Optional known harness, harness version, and model provenance. |
| `repository.name` | Repository directory name. |
| `repository.root` | Canonical absolute worktree root recorded on this machine. |
| `repository.identity` | Normalized remote URL, or a `file` URL for a repository without a remote. |
| `repository.remote` | Optional normalized remote URL. |
| `git.commit` | Full lowercase Git object ID. |
| `git.branch` | Optional branch name; absent for a detached `HEAD`. |
| `git.subject` | Commit subject. |
| `git.authored_at` | Git author time normalized to UTC. |
| `git.committed_at` | Git committer time normalized to UTC. |
| `git.worktree_clean` | Worktree cleanliness observed during recording. |
| `note` | Optional significance note. |

A correction repeats the complete repository and Git evidence. It also stores
`target_id`, a required reason, and the effective replacement note.

## Repository remote normalization

Acta removes passwords, query parameters, and fragments from repository
remotes. It removes user information from HTTP and HTTPS URLs. SSH and Git URLs
retain a non-secret username when one exists. The normalized identity still
contains the remote host, optional port, and repository path.

The journal therefore contains local and repository-identifying information.
Keep it private. Do not commit journal files to a repository or attach them to
an issue without reviewing and redacting them.

## Data not collected automatically

Acta does not collect prompts, conversations, hidden reasoning, diffs, patches,
file contents, environment variables, terminal history, hours, effort, or
productivity. A user or agent can still place sensitive text in a note, so the
capture policy must prohibit secrets and personal data.

## Bounds and rotation

- Journal segments rotate by UTC month.
- One encoded event line can be at most 8 MiB.
- One local service request can be at most 64 KiB.
- One local service response can be at most 8 MiB.
- Search and report return at most 1000 effective events per request.
- A note or correction reason can contain at most 500 Unicode code points.
