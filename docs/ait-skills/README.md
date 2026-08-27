# AIT Shared Skill Bundle

Shared, agent-agnostic port of the `ait` skills. This bundle is the **single
live contract** for using `ait` (Agent Issue Tracker) from any coding agent:
Codex, Claude, OpenCode, Droid, Pi, and others. It is a copy adapted from the
Claude-only template, not a fork — one contract, no per-agent variants.

> Rule: keep exactly one live contract. Do not fork per-agent copies, do not
> introduce v2/v3 of this bundle. Edit the files here and re-verify parity.

## Bundle layout

```
docs/SKILLS/ait/
├── README.md            # this manifest
├── ait/                 # core tracker skill
│   ├── SKILL.md         # commands, hierarchy, roles, contract
│   └── DELEGATION.md    # export → delegate → reconcile workflow
├── ait-fleet/
│   └── SKILL.md         # multi-repo orchestration via --db + manifest
└── ait-recap/
    └── SKILL.md         # friendly recap of recent tracker activity
```

## Parity with the original

Source of truth: `TOOLS/ait/claude/skills/` (Claude-only template, never
edited by this bundle).

| File | vs original | Changes |
|---|---|---|
| `ait/SKILL.md` | adapted | stripped Claude-only frontmatter (`allowed-tools`, `author`, `license`); claim identity canonicalised to `PASEO_AGENT_ID` (removed the playful "pick a fun name" paragraph); added **Role permissions** (Lead/Peer/Supervisor), **No automatic issue creation**, **Review & Handback** (review optional/risk-based), and **Bookkeeping, not staffing** (hierarchy is optional bookkeeping) notes; command contract re-verified against the installed `ait` binary (exit codes 64 usage / 65 conflicts-validation / 1 uninitialised, error envelopes, `--human`/`--tree` exclusion, `hidden_count` always present under default `list`, sibling rekeying on delete); `init` marked Lead-only |
| `ait/DELEGATION.md` | adapted | export marked read-only (all seats); reconcile actions mapped to the role matrix; kept the export→delegate→reconcile workflow |
| `ait-fleet/SKILL.md` | adapted | stripped Claude-only frontmatter (large `allowed-tools` list); tool wording generalised ("use the Write tool" → "write the file"); member-repo `ait init` marked Lead-only; kept fleet manifest + cross-repo conventions |
| `ait-recap/SKILL.md` | adapted | stripped Claude-only frontmatter; noted read-only so it is safe for Supervisor; recap templates preserved unchanged |

Deliberately **not** copied: `claude/skills/.DS_Store` and any Claude-specific
frontmatter fields.

## Frontmatter contract

Only three portable fields survive in every `SKILL.md` frontmatter:

```yaml
name: <skill-name>
description: > ...
version: "1.0.0"
```

Fields that are Claude-only (`allowed-tools`) were removed. Do not add
agent-specific fields back; keep the frontmatter portable so every agent's
skill loader accepts it.

## Role permissions (Lead / Peer / Supervisor)

The tracker is a cooperative advisory surface, not a lock. Command boundaries:

| Command group | Lead | Peer | Supervisor |
|---|---|---|---|
| `init` | ✅ | ❌ | ❌ |
| `close --cascade` / `cancel` / `flush` / `delete` / `log purge` / `reopen` | ✅ | ❌ | ❌ |
| `claim` (assign to others), `unclaim` of others | ✅ | ❌ | ❌ |
| `claim`/`unclaim` own work, `update`, `note add`, `close` own task | ✅ | ✅ | ❌ |
| `ready --type task` | ✅ | ✅ | ❌ |
| `status`, `list --long`, `show`, `search`, `export` | ✅ | ✅ | ✅ |
| `config` | ✅ | ❌ | ✅ |
| `log` (flush history) | ✅ | ✅ | ✅ (inspect only) |

Canonical claim identity is the `PASEO_AGENT_ID` environment variable —
never a display name or alias.

## No automatic issue creation

`ait` never creates issues from conversation, prompts, or session activity,
and there is no hook wired in this bundle to auto-create them. Issue creation
is explicit via `ait create`. The skills never claim otherwise.

## Review is optional

Independent peer review is not a default gate on every handback. It is
risk-based and usually opened at plan closure with an explicit independent
question and a proof bar. A stable handback is accepted by the owning seat
directly. Any owner-directed review on a specific episode is a scoped decision
for that episode, not the standing workflow.

## Verified command contract

Verified against `/usr/local/bin/ait` (version `dev`) while authoring this
bundle. Key confirmed behaviours:

- Claim conflict → exit code **65**, error code `conflict`; cycle in `dep add`
  → exit **65**, code `validation`; `delete` without `--force` → exit **65**,
  code `confirmation`.
- Uninitialised database → exit **1**, error code `uninitialised`.
- Usage errors (e.g. combined `--human`/`--tree`) → exit **64**, code `usage`.
- `close --note` works (`--reason` is an alias); `create --description @file`
  reads description from a file; `export <id> --output briefing.md` produces a
  self-contained markdown briefing.
- `init` writes `.ait/` into `.gitignore` in a git repo and reports
  `schema_version` (currently 4) in `ait config`.
- `status` counts: `{cancelled, closed, in_progress, open, ready, total}`.

## Per-agent consumption

- **OpenCode**: register each `SKILL.md` (or this bundle) in the agent's
  skills dir; the minimal frontmatter (`name`, `description`, `version`) is
  accepted as-is.
- **Claude**: drop the bundle into the skills dir; the removed `allowed-tools`
  field is the only loss, which is intentional (no allowlist in the shared
  contract).
- **Codex / Droid / Pi**: no native skill registry — point the agent at this
  folder (or the relevant `SKILL.md`) as an instruction/reference file.
- **Any agent**: the workflow, role matrix, and command contract in
  `ait/SKILL.md` are provider-agnostic; claim identity is `$PASEO_AGENT_ID`.

## Origin & licence

Adapted from the MIT-licensed `ait` skills by ohnotnow
(`<https://github.com/ohnotnow>`). Source template:
`TOOLS/ait/claude/skills/`. This adaptation keeps the MIT licence terms; see
`TOOLS/ait/LICENCE`.
