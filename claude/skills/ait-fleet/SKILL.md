---
name: ait-fleet
description: >
  Multi-repo orchestration for `ait`. Use when a project spans more than one
  git repository (e.g. separate frontend and backend repos) and you need to
  plan, track, or coordinate work across them. Triggers on mentions of
  "ait fleet", "across repos", "multi-repo ait", "frontend and backend ait",
  or any request to set up or use ait across more than one repo.
allowed-tools: "Read,Write,Bash(ait:*),Bash(test:*),Bash(ls:*),Bash(mkdir:*),Bash(cat:*),Bash(realpath:*)"
version: "0.1.0"
author: "ohnotnow <https://github.com/ohnotnow>"
license: "MIT"
---

# AIT Fleet — Multi-Repo Orchestration

`ait` is deliberately scoped to a single repo's database. When a project
spans more than one repo (frontend + backend, or backend + mobile + web),
this skill provides a convention for coordinating work across them using
`ait`'s existing `--db <path>` flag — no binary changes, just a manifest
file and a workflow.

## When to Use This Skill

- The user mentions a fleet, a multi-repo project, or "frontend and backend".
- The user asks to create, list, or update an `ait` fleet config.
- The user is planning a feature that needs work in more than one repo.
- The user asks "what's ready across all my repos?" or similar.

If the user is only working in one repo, defer to the regular `ait` skill.

## The Manifest

Fleet configs live at `~/.config/ait-fleet/<name>.json`:

```json
{
  "name": "myproject",
  "primary": "backend",
  "members": {
    "backend":  { "path": "/abs/path/to/backend" },
    "frontend": { "path": "/abs/path/to/frontend" }
  }
}
```

- `name` — the fleet identifier (used as the filename).
- `primary` — alias of the repo that owns cross-cutting initiatives.
- `members` — alias → metadata. The DB path is `<path>/.ait/ait.db` unless
  an explicit `db` field is set on the member.

The manifest doesn't dictate where repos live on disk. Sibling layout
(`project/frontend`, `project/backend`) and scattered layout
(`~/code/api`, `~/work/clients/foo/web`) are both fine.

## Setup Walkthrough

When the user asks to create a new fleet config, walk them through the
following steps. Confirm before each action that touches their filesystem.

### 1. Look for existing manifests first

```bash
ls ~/.config/ait-fleet/ 2>/dev/null
```

If a manifest already covers the repos they mention, use it instead of
creating a duplicate. If they want to add a repo to an existing fleet,
read the JSON, add the new member, and re-validate.

### 2. Gather the repo paths

Ask the user for the absolute path to each member repo. They may give
`~`-prefixed or relative paths — resolve them with `realpath` so the
manifest stores absolute paths.

If the user hasn't told you which is the "primary" (where cross-cutting
initiatives will live), ask. Default suggestion: the backend, since
initiatives often start with API or data work. The user can override.

### 3. Validate each path

For each candidate path, run these checks and surface the result to the
user before continuing:

```bash
test -d <path>            # is it a directory?
test -d <path>/.git       # is it a git repo?
test -f <path>/.ait/ait.db  # does an ait DB already exist?
```

If the path isn't a directory or isn't a git repo, stop and ask the user
to clarify — don't silently skip.

If `.ait/ait.db` doesn't exist, offer to initialise it:

```bash
(cd <path> && ait init)
```

Don't run `ait init` without explicit confirmation — it's the user's repo
and `init` writes the prefix.

### 4. Read each repo's prefix

```bash
ait --db <path>/.ait/ait.db config
```

Show the user each repo's prefix. **If two repos share a prefix, stop.**
The cross-repo `<alias>:<id>` convention loses meaning when prefixes
collide. Suggest re-running `ait init --prefix <distinct>` in one of them
before continuing.

### 5. Pick aliases

The alias is the short name used inside the manifest and in cross-repo
references. Default to the repo's directory basename. Ask the user to
confirm or rename — they may want something shorter or more descriptive
(e.g. `api` instead of `myproject-server`).

### 6. Show the proposed manifest and confirm

Print the JSON. Ask the user to confirm before writing. This is the last
chance to catch a wrong path or alias.

### 7. Write the manifest

```bash
mkdir -p ~/.config/ait-fleet
```

Then use the `Write` tool to create `~/.config/ait-fleet/<name>.json`
with the confirmed JSON.

### 8. Show the user what they can do next

Tell them the fleet is ready and remind them of the basic operations
(below). Suggest they create a cross-cutting initiative in the primary
repo if there's already a feature they're planning.

## Working with a Fleet

Once a manifest exists, treat it as the source of truth for which DBs to
talk to. At the start of any fleet-aware session, read the manifest and
re-validate paths (`test -d`) — repos move and get renamed.

### List available fleets

```bash
ls ~/.config/ait-fleet/
```

### Fleet-wide ready

For each member in the manifest:

```bash
ait --db <member-db> ready
```

Merge the output mentally and report it back to the user grouped by
member alias. Cross-repo blockers (notes — see below) won't filter through
`ait ready`, so scan recent notes for the soft-blocker pattern separately.

### Fleet-wide status

```bash
ait --db <each-db> status
```

### Show a cross-cutting initiative

The initiative lives in the primary repo. Run `ait --db <primary-db>
show <id>` for the strategic context, then for each linked member epic
listed in its description or notes, run `ait --db <member-db> show
<epic-id>` for that slice's tasks.

## Cross-Repo Conventions

`ait` does not natively link issues across DBs. This skill encodes those
links as text so they survive without any binary changes. The conventions
are intentionally greppable — you can scan a member's notes for a string
pattern when `ait ready` can't see across DBs.

### Initiative + linked epics

1. Create the initiative in the **primary** repo:

   ```bash
   ait --db <primary-db> create --type initiative \
     --title "Add feature X" \
     --description @initiative.md
   ```

   In the description, list the affected members:

   ```markdown
   Fleet: myproject
   Affected members: backend, frontend

   See linked epics in each member repo (notes on this initiative).
   ```

2. For each affected member, create an epic in that repo's DB:

   ```bash
   ait --db <member-db> create --type epic \
     --title "Feature X — backend slice" \
     --description "Part of initiative <full-initiative-id> in fleet 'myproject' (primary repo)."
   ```

3. Add a note to the primary initiative for each linked epic:

   ```bash
   ait --db <primary-db> note add <init-id> \
     "Linked epic: backend:<full-epic-id>"
   ```

   The literal `Linked epic: <alias>:<id>` form keeps it greppable.

4. Add tasks under each member epic in the usual way. Tasks stay scoped
   to their own repo's DB — no cross-repo linking needed at the task level
   unless a real blocker exists (see below).

### Cross-repo blockers

`ait dep add` only works within a single DB. When a frontend task depends
on a backend task landing first, express it as a note on the blocked task:

```bash
ait --db <frontend-db> note add <task-id> \
  "Blocked by backend:<full-task-id>"
```

Use the literal prefix `Blocked by <alias>:<id>` so the convention is
greppable. Before claiming or starting work in one repo, read recent notes
on candidate tasks for that pattern.

This is **advisory only** — `ait ready` won't filter on it. The skill
relies on you (or the agent) to actually look.

### Reconciliation

When a member epic finishes:

```bash
ait --db <member-db> close <epic-id> --cascade
ait --db <primary-db> note add <init-id> \
  "backend slice complete (closed YYYY-MM-DD)"
```

When all member epics are done, close the initiative in the primary repo.

## Delegating Per-Repo Slices

When handing a single member's slice to a sub-agent, use the regular
`ait` skill's `DELEGATION.md` export/reconcile pattern, scoped to that
member's DB:

```bash
ait --db <member-db> export <epic-id> --output briefing.md
```

The sub-agent works in that repo with the briefing — it doesn't need to
know about the fleet at all. Reconcile the closed tasks back into that
member's DB when done, and update the initiative note in the primary.

## Tips and Pitfalls

- **Validate manifest paths each session.** Repos move, get renamed, or
  get cloned to different paths. If `test -d` on a member path fails, tell
  the user and offer to update the manifest before doing anything else.
- **Distinct prefixes are mandatory.** If two repos share a prefix, the
  `<alias>:<id>` convention collapses. Catch this at setup time.
- **Don't try to share IDs across repos.** Each repo has its own ID space.
  Cross-references are textual only.
- **Keep the fleet small.** Two or three repos is the sweet spot. Beyond
  that, coordination via notes gets noisy and the lack of native dependency
  tracking starts to bite — at which point a different tool may be the
  right answer.
