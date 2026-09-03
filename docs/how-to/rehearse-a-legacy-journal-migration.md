# Rehearse a Legacy Journal Migration

Use an isolated copy of a legacy events directory to validate migration before
an installation cutover. This procedure does not replace the installed binary
or live journal.

## Build the candidate

Build without installing over the active executable:

```sh
go build -o /tmp/acta-candidate .
```

## Copy the legacy events directory

Choose a private rehearsal directory on the same machine. Copy the complete
legacy events tree while the source copy is stable:

```sh
cp -a "$HOME/.acta/events" /tmp/acta-legacy-rehearsal
```

The migrator accepts the events root itself, not the containing `.acta`
directory. The source must contain only year directories and monthly JSONL
segments.

## Run two independent migrations

Ensure none of the destination or manifest paths exists, then run:

```sh
/tmp/acta-candidate migrate \
  --source /tmp/acta-legacy-rehearsal \
  --destination /tmp/acta-canonical-first \
  --manifest /tmp/acta-manifest-first.json

/tmp/acta-candidate migrate \
  --source /tmp/acta-legacy-rehearsal \
  --destination /tmp/acta-canonical-second \
  --manifest /tmp/acta-manifest-second.json
```

Migration fails rather than skipping an unsupported, malformed, partial, or
lossy legacy event. It never modifies the source.

## Compare deterministic output

Compare the complete trees and manifests:

```sh
diff -r /tmp/acta-canonical-first /tmp/acta-canonical-second
cmp /tmp/acta-manifest-first.json /tmp/acta-manifest-second.json
```

Both commands must produce no differences.

## Verify with the candidate

The normal verifier uses the fixed Acta journal location, so cutover rehearsal
validation is performed internally by `migrate` before it publishes either
destination. Inspect the manifest to confirm segment count, event count,
ordered IDs, effective-event count, correction chains, and source and
destination SHA-256 values.

Do not use a rehearsal destination as the live journal. The final cutover also
requires a writer freeze, recoverable snapshots of the installed executable
and complete journal, staged migration from the frozen source, manifest
comparison, ordered same-filesystem replacement, service restart, and smoke
checks. Retain rollback snapshots until explicit operator acceptance.
