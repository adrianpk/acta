# Local Journal Model

Acta uses one per-user local service as the sole canonical writer. Agents and
CLI processes send bounded operations through an owner-only Unix socket. This
keeps concurrent writes serialized without requiring agents to coordinate file
access.

Events are complete JSON objects stored one per line in monthly UTC segments.
Appending to an existing segment is cheaper and easier to inspect than creating
one file for every event. Monthly rotation keeps individual files bounded and
avoids one permanent journal file that grows for the lifetime of the tool.

Corrections append a new self-contained event that supersedes the current
effective event. Published lines do not change. Search and reports operate on
the effective event in each correction chain.

Acta keeps all journal and query work local. It does not send events to an
external search service. Git remains the source of code, diffs, specifications,
plans, and detailed history; Acta stores the route back to that evidence.
