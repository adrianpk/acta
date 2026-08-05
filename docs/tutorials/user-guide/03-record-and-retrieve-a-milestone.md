# Record and Retrieve a Milestone

Record one completed commit, then retrieve it through the local chronicle.

## Create a milestone commit

Use an example Git repository with a completed change. Confirm the commit that
you intend to record:

```sh
git show --no-patch --oneline HEAD
```

## Record the commit

Run Acta from anywhere inside that worktree:

```sh
acta record "Established the first reusable workspace policy"
```

Acta prints `recorded`, an event ID, the repository name, and an abbreviated
commit. Repeating the command for the same repository and commit prints
`already recorded` and does not append a duplicate milestone.

## Find the milestone

Search the stored subject and note text:

```sh
acta search --query "workspace policy"
```

The result identifies the recording time, event type, repository, abbreviated
commit, and commit subject.

Use the commit from that result to retrieve the complete effective event:

```sh
acta show --commit <commit> --json
```

The JSON result supplies the repository and Git evidence that an agent can use
to inspect the original implementation.

## Review a bounded period

Use local calendar dates to group recorded milestones by day and repository:

```sh
acta report --from 2026-08-01 --to 2026-08-08 --timezone UTC
```

The `from` date is inclusive. The `to` date is exclusive.

Return to the [user guide](index.md), or use the [how-to
guides](../../how-to/index.md) for recurring tasks.
