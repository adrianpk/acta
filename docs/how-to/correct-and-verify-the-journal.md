# Correct and Verify the Journal

Append a correction when a milestone note is inaccurate. Verify the journal
when storage integrity is in doubt or after restoring local data.

## Find the effective event ID

```sh
acta show --commit <commit> --json
```

Read the event `id` from the result.

## Append a correction

Replace the note:

```sh
acta correct \
  --target <event-id> \
  --reason "Clarify the affected behavior" \
  --note "The corrected significance note"
```

Omit `--note` to retain the current note. Pass `--note ""` to clear it. A
correction must target the current effective event. Acta appends a new event
and does not rewrite the earlier line.

## Verify journal integrity

Verify the active monthly segment:

```sh
acta verify
```

Verify every canonical segment:

```sh
acta verify --full
```

Use `--json` when an agent needs structured results. On corruption, Acta
reports the segment, line, byte offset, and reason. Verification does not
perform general repair.
