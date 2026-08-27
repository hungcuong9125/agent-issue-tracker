---
name: ait
description: >
  Local-first issue tracker for coding agents. Shared across agents — Codex,
  Claude, OpenCode, Droid, Pi — with a single portable contract and no
  agent-specific allowlists. Use when planning work, tracking multi-step
  tasks, modelling dependencies, coordinating between agents, or resuming
  after session loss or conversation compaction. Use this skill if the user
  mentions `ait`.
version: "1.0.0"
---

# AIT (`ait`) — Agent Issue Tracker Quick Reference

AIT is a lightweight, SQLite-backed CLI issue tracker designed for coding
agents. Use it instead of TodoWrite for tracking work that spans sessions,
has dependencies, or involves multiple agents.

This skill is the single live contract for `ait` across agents. Keep exactly
one contract: do not fork per-agent copies, do not version this file. Change
it here and re-verify parity with the original skill bundle when you do.

## Essential Commands

### View Issues
```bash
ait status                          # Overview: counts by status
ait list                            # List open issues (slim JSON, 5 fields)
ait list --long                     # Full JSON record (all fields)
ait list --human                    # Compact tabular view for humans
ait list --tree                     # Parent-child hierarchy with connectors
ait list --type epic                # Filter by type
ait list --status open              # Filter by status
ait list --priority P1              # Filter by priority
ait list --parent <id>              # Children of a specific epic
ait list --all                      # Include closed/cancelled issues
ait ready                           # Unblocked issues, ordered by priority
ait ready --type task               # Unblocked tasks only (excludes epics)
ait ready --long                    # Unblocked issues with full detail
ait show <id>                       # Full details, children, blockers, notes
ait search "keyword"                # Search issues by text
ait config                          # Show project prefix and schema version
```

### Create Issues
```bash
ait create --title "Title"                              # Basic task
ait create --title "Title" --type initiative             # Initiative (strategic vision)
ait create --title "Title" --type epic                   # Epic (group of tasks)
ait create --title "Title" --type epic --parent <init-id>  # Epic under an initiative
ait create --title "Title" --parent <epic-id>            # Task under an epic
ait create --title "Title" --description "Details..."    # With description
ait create --title "Title" --description @spec.md        # Description from file
ait create --title "Title" --priority P1                 # With priority
```

### Update Issues
```bash
ait update <id> --status in_progress   # Start working
ait update <id> --status open          # Back to open
ait update <id> --title "New title"    # Change title
ait update <id> --priority P0          # Change priority
```

### Close / Cancel / Reopen
```bash
ait close <id>                # Close a single issue
ait close <id> --cascade      # Close an epic and all its descendants
ait close <id> --note "Done — merged in PR #42"    # Add a closing note then close (--reason is an alias)
ait cancel <id>               # Cancel an issue
ait reopen <id>               # Reopen a closed or cancelled issue
```

### Dependencies
```bash
ait dep add <id> <blocker-id>      # <id> is blocked by <blocker-id>
ait dep remove <id> <blocker-id>   # Remove a dependency
ait dep list <id>                  # Show dependencies for an issue
ait dep tree <id>                  # Show dependency tree
```
Cycle detection is built in — adding a dependency that would create a cycle is
rejected (exit code 65, error code `validation`).

### Notes
```bash
ait note add <id> "Note body text"   # Attach a note to an issue
ait note list <id>                   # List notes for an issue
```

### Flush (Housekeeping)
```bash
ait flush              # permanently delete all closed/cancelled issues
ait flush --dry-run    # preview what would be deleted without changing anything
ait flush --summary "Fixed pg compatibility, added API docs"  # with editorial note
```
Flush records all flushed issues to a history log before deleting them.
The `--summary` flag attaches a short editorial note to the history entry —
useful for giving future sessions a quick description of what was accomplished.

Flush removes root-level issues whose entire descendant tree is also closed or
cancelled. Notes and dependencies are cascade-deleted automatically.

**Important:** If the `skipped` list in the response is non-empty, it means a
closed epic has open or in-progress children — something that probably needs
human attention. Flag this to the user and suggest they review the skipped
issues before deciding what to do.

### Delete (mistakes only)
```bash
ait delete <id> --force             # permanently remove a single issue
ait delete <id> --force --cascade   # remove the issue and its whole subtree
```
`delete` is for genuine mistakes — a fat-fingered duplicate, a throwaway, an
issue created against the wrong parent. It is **irreversible** and, unlike
flush, records **nothing**: the issue, its notes, and its dependency links are
gone for good. The response is `{ "deleted": [refs] }`, listing exactly what
was removed.

It is deliberately guarded:
- nothing happens without `--force` (exit 65, error code `confirmation`);
- an issue that has children is refused unless you also pass `--cascade`.

Deleting a child also rekeys its surviving siblings (e.g. `.2` becomes `.1`),
so re-read the IDs of siblings after a delete rather than trusting a
previously recorded ID.

Reach for `cancel` or `close` (which keep an auditable record) when you're
closing out real work. Use `delete` only when an issue genuinely should never
have existed.

### Flush History
```bash
ait log                           # summary: date, summary, root items, item count
ait log --long                    # full detail: all items with close reasons
ait log --last 3                  # most recent 3 flush events
ait log --since 2026-04-01        # flushes since a date
ait log --search "migration"      # find items by title or close reason
ait log --search "auth" --long    # search with full detail
ait log purge --keep 20           # compact: keep summaries, drop items for old entries
ait log purge --keep 10 --full    # fully delete old entries
ait log purge --before 2026-01-01 # compact entries older than a date
```

The default `log` output is slim: each flush entry shows its date, summary,
total item count, and only root-level items. Use `--long` for all items
including children and close reasons.

Use `--search` when the user mentions past work vaguely ("we changed the
migrations a while back") — it matches against item titles and close reasons.

`log purge` defaults to **compact** mode: summary rows are kept, per-issue
items are dropped. Use `--full` to delete entries entirely. Scope with
`--keep <n>` or `--before <date>` (mutually exclusive).

### Claiming (Multi-Agent)
```bash
ait claim <id> "$PASEO_AGENT_ID"   # Claim an issue (prevents duplicate work)
ait unclaim <id>                   # Release the claim
```
If another agent already holds the claim, `claim` returns a conflict error
(exit code 65, error code `conflict`) with the current holder's name.

**Canonical identity.** Always claim with the canonical agent identity
`PASEO_AGENT_ID` — the environment variable that every agent session exposes.
Never use a display name, a model label, or a creative alias as the claim
name. One agent, one claim identity, across every provider. Before claiming,
read the issue and check the holder so you do not fight an active owner.

## Issue Types & Hierarchy

The three issue types form a strict hierarchy: **initiative > epic > task**.

- `initiative` — the strategic "why": vision, goals, and key decisions behind a group of epics. **Top-level only** (cannot have a parent).
- `epic` — container for related tasks. Can be top-level or a child of an initiative. **Cannot** be a child of another epic or task.
- `task` (default) — a unit of work. Child of an epic or another task (for subtasks). **Cannot** be a direct child of an initiative.

To build a full three-tier structure:
1. Create the initiative: `ait create --title "Vision" --type initiative`
2. Create an epic under it: `ait create --title "Phase 1" --type epic --parent <initiative-id>`
3. Create tasks under the epic: `ait create --title "Do X" --parent <epic-id>`

**Common mistake**: trying to add a task directly under an initiative. You need an epic in between.

**Bookkeeping, not staffing.** The initiative → epic → task hierarchy is
optional bookkeeping, not a mandatory staffing or decomposition structure.
A single outcome keeps a continuous owner from discovery through
implementation to proof/validation. Create a child issue or a second seat only
when there is an independent decision to make, or a disjoint writable/proof
boundary that is genuinely frozen. Do not split work by file or phase, and do
not bolt extra roles (e.g. scout/council) onto the three tracker seats. If
the work in front of you outgrows the outcome, flag that to the Lead — do not
silently expand the tracker.

**Connecting to intake.** If this workspace uses `docs/FEATURE_INTAKE.md`, its
lane maps onto the hierarchy for `epic` vs `task`, but not for `initiative`:
`tiny`/`normal` work is a plain task (no epic needed); `high-risk` work that
requires an active ExecPlan (`docs/PLANS.md`) becomes an epic, with each
"Direction And Work Units" item as a task underneath it. Neither lane decides
`initiative` — that is a separate, rarer call about whether a whole
multi-epic program shares one strategic vision, not about any single
feature's risk.

**Before creating a new epic**, run `ait list --type initiative`. Nest the
new epic under an existing initiative only when it genuinely shares that
initiative's "why" — not by default just because one exists. A top-level epic
is normal, not a loose end; do not nest to avoid a bare root. Several open
initiatives in one project is normal too — initiative count is not a target to
minimize or a structure to consolidate into one. Create an initiative only
when a *second* epic needs that shared "why" to be readable in one place; a
single-epic initiative is an epic with extra ceremony. If this workspace has
`framework/issue-policy.md`, it is the canonical statement of this rule —
read it before restructuring an existing tracker.

## Priorities
- `P0` — critical / urgent
- `P1` — high priority
- `P2` — normal (default)
- `P3` — low priority
- `P4` — nice to have

## Hierarchical IDs

IDs are auto-generated with the project prefix:
- Root issue: `<prefix>-<sqid>` (e.g. `ait-AXs1i`)
- First child: `<prefix>-<sqid>.1` (e.g. `ait-AXs1i.1`)
- Grandchild: `<prefix>-<sqid>.1.1`

The parent-child structure is visible directly in the identifier. For a full
three-tier setup: `proj-abc` (initiative) -> `proj-abc.1` (epic) -> `proj-abc.1.1` (task).

## Workflow Pattern

1. **Start of session**: `ait log --last 3` for recent context, then `ait ready` to see what is unblocked
2. **Pick work**: `ait claim <id> "$PASEO_AGENT_ID"` to claim an issue
3. **Check context**: `ait show <id>` for full details and notes. If the issue belongs to an initiative, read the initiative description to understand the strategic intent.
4. **Mark in progress**: `ait update <id> --status in_progress`
5. **Do the work**: implement, test, iterate
6. **Leave notes**: `ait note add <id> "what was done / what remains"`
7. **Close**: `ait close <id>` (or `--cascade` for an epic and its children, Lead only)
8. **Next**: `ait ready` again to find the next unblocked item
9. **End of session**: `ait status` for an overall summary

## Output Modes

By default all commands return JSON — compact and token-efficient for agents.
Data goes to stdout; failures are a JSON `{"error": {"code", "message"}}`
envelope on **stderr** with a non-zero exit code (64 for usage errors, 65 for
conflicts/validation, 1 for uninitialised), so pipelines and `$(...)`
captures never ingest an error as data. Same contract as `ant`.

- `--long` adds all fields (description, timestamps, claimed_by, etc.)
- `--human` gives a compact tabular view grouped by epic
- `--tree` shows parent-child hierarchy with tree connectors
- `--human` and `--tree` are mutually exclusive (combined use is a usage error, exit 64, code `usage`)
- All display modes support the same filters (`--type`, `--status`, `--priority`)

### Mutating commands return slim refs

`create`, `update`, `close`, `cancel`, `reopen`, `claim`, `unclaim` return a
slim ref `{id, title, status, type, priority}` — enough to chain into the next
command without burning context on a full record echo. `dep add`, `dep remove`,
and `note add` return a slim ack `{ok: true, ...ids}`. Pass `--long` on any of
these to get the full record back when you actually need it (e.g. confirming
a description was set, or reading `claimed_by` after a claim).

### `list` and `hidden_count`

By default `ait list` excludes closed and cancelled issues. The response
always includes a `hidden_count` field telling you how many issues are being
filtered out (0 when nothing is hidden), so an empty-looking response when the
project is full of closed work is never a surprise. Pass `--all` to see
everything; under `--all` the `hidden_count` field is omitted.

## Initialisation (Lead only)

```bash
ait init --prefix myproject    # Create the database and set the ID prefix
```
`init` is the only command that creates the database. Every other command
refuses until it has been run once, returning exit code 1 and:

```json
{"error": {"code": "uninitialised", "message": "no ait database at <path> — run 'ait init' first"}}
```

`init` is a Lead-only command in this workspace's protocol. Peers and
Supervisor seats never run it. If you hit `uninitialised`, don't run `init`
reflexively — a project without an ait database may simply not use ait. Check
with the Lead or Human before initialising a project that isn't yours.

If no prefix is set, one is inferred from the directory name. The prefix can be
changed later with `init --prefix` — existing IDs are re-keyed automatically.
In a git repository, `init` also ensures `.ait/` is in `.gitignore`; outside
one, its output carries a `note` saying that step was skipped.

## Custom Database Path

```bash
ait --db /path/to/other.db list   # Use a different database file
```
Useful for git worktrees (pointing back to the main repo's database) or keeping
separate databases for different subsystems. Run with cwd at the project root,
or pass an explicit shared `--db`, so the correct `.ait/ait.db` is read.

## No Automatic Issue Creation

`ait` never creates issues from conversation, prompts, or session activity.
There is no hook that turns a user mention into a tracker entry. Issue
creation is explicit only: someone runs `ait create`. If you ever see
issues appear without an explicit `create`, the mechanism is an external
integration (a hook or orchestrator) — not `ait` itself — and this skill does
not assume or advertise one.

## Role Permissions

This workspace separates three seats. The tracker is a cooperative advisory
surface, not a lock: `closed` is a work-surface signal, never technical
acceptance truth. Command boundaries per seat:

- **Lead** (full write, decomposition owner): `init`, `close <id> --cascade`,
  `cancel`, `flush`, `delete`, `log purge`, `reopen`, claim-to-assign
  (`ait claim <id> <PASEO_AGENT_ID>` before dispatch), and `unclaim` of
  issues held by others.
- **Peer** (one execution seat): claim/unclaim **own** work, `update`,
  `note add`, `close` own task (no cascade), `ready --type task`, and the
  read commands `show`/`list`/`log`/`export`/`search`/`status`. Never
  cascade-close, never unclaim an issue you do not hold, never run
  destructive housekeeping (flush/delete/log purge/cancel).
- **Supervisor** (event-driven read-only spot-checks): `status`,
  `list --long`, `show`, `search`, `config`, `export`; `log` only to inspect
  flushed history. Never polling, never acceptance, never issue writes.

Ready-work selection: `ait ready --type task` surfaces unblocked issues from
the tracker's bookkeeping. It is a suggestion, not a work-order dispatch — the
Lead still owns scope and assignment, so pick only work that is genuinely
yours to take. `ready` includes already-claimed issues and open/in_progress
items; treat a claim conflict (exit 65) as a signal to pick a different task,
not an error.

## Review & Handback

Peer review is optional and risk-based, not a default gate on every handback.
A stable handback can be completed and accepted by the owning seat directly;
open a separate review only when a decision carries real risk — typically at
plan closure — and state both the independent question being checked and the
proof bar it must meet. A reviewer validates that specific question; it does
not re-run the work, block routine handback, or become a standing workflow
step. A given episode may still add an owner-directed review — that is a
scoped decision for that episode, not the default contract.

## Delegating Work

If you are supervising sub-agents or delegating work to agents that don't have
access to `ait`, see `DELEGATION.md` (in this skill directory) for the
export → delegate → reconcile workflow.

## Tips

- Prefer `ait ready` over `ait list` when deciding what to work on — it filters
  to unblocked issues and sorts by priority.
- Use notes liberally — they survive session loss and conversation compaction.
- Use `--cascade` on close to avoid closing children one by one (Lead only).
- Run `ait flush` periodically to keep the database lean — the tracker is for
  ephemeral work, so there is no need to keep completed issues forever. Use
  `--summary` to leave a note for future sessions about what was accomplished.
- Use `ait log --search` when the user references past work — it searches
  titles and close reasons across all flush history.
- `ait show <id>` returns children, blockers, and notes in one call — use it to
  get full context before starting work.
- The database lives at `.ait/ait.db` in the git root. It is a plain SQLite file
  and easy to inspect or back up.
