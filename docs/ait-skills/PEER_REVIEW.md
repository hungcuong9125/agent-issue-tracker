# Peer Review — AIT Shared Skill Bundle

- **Reviewer**: Peer (opencode-go/deepseek-v4-flash), `PASEO_AGENT_ID=fa2e3893-1615-4db6-9df2-701a24e34a46`
- **Issue**: pdt-UkLWZ.4.1 (independent peer review)
- **Scope reviewed**: `docs/SKILLS/ait/` (README.md, ait/SKILL.md, ait/DELEGATION.md, ait-fleet/SKILL.md, ait-recap/SKILL.md)
- **Template (source of truth)**: `TOOLS/ait/claude/skills/` (Claude-only, read-only)
- **Binary under test**: `/usr/local/bin/ait` (version `dev`)
- **Method**: full-file `diff` per file; live command-contract verification against the installed binary using throwaway DBs under `/tmp`; grep sweep for forbidden wording; YAML frontmatter parse. Workspace `.ait` DB untouched; no Lead-part files modified.
- **Date**: 2026-08-03

## Verdict

**APPROVED — FOLLOW-UP CLOSED**

The bundle is structurally correct, parity-preserving, role-correct, and
empirically accurate on every required command-contract item (exit codes 65/1,
error envelopes, `--human`/`--tree` exclusion, `close --note`/`--reason`,
`create --description @file`, `export --output`, `show` shape, `status`
counts, `config` schema_version, `init` gitignore behaviour). No original
capability was lost in the adaptation, no role boundary is breached, and no
false auto-issue claim exists. The previously recorded M1/L1/L2/N1 findings
are resolved in the current bundle. **Open findings: 0.**

---

## (a) File-by-file parity table vs the template

| Bundle file | Template file | Result | Notes |
|---|---|---|---|
| `ait/SKILL.md` | `claude/skills/ait/SKILL.md` | **Adapted, no capability lost** | Stripped `allowed-tools`/`author`/`license`; `version` 0.1.0→1.0.0; `name` `ait-usage`→`ait`; removed "pick a fun name" paragraph; added Role Permissions, No Automatic Issue Creation, canonical `PASEO_AGENT_ID` identity, explicit exit codes, `init` Lead-only. Every command block (view/create/update/close/cancel/reopen/dep/note/flush/delete/log/claim/unclaim/init/config/export/status/ready/search) preserved verbatim except documented edits. |
| `ait/DELEGATION.md` | `claude/skills/ait/DELEGATION.md` | **Adapted, additions only** | Added "`ait export` is read-only / safe for every seat"; reconcile actions mapped to the role matrix (cascade close = Lead only, partial close = Peer own / Lead any, note add = Peer+Lead, creation = Lead). Export→delegate→reconcile workflow intact. |
| `ait-fleet/SKILL.md` | `claude/skills/ait-fleet/SKILL.md` | **Adapted, no capability lost** | Stripped large `allowed-tools` list; `version` 0.1.0→1.0.0; "use the Write tool"→"write the file"; member-repo `ait init` marked Lead-only; `close --cascade` marked Lead action. Manifest format, cross-repo conventions (initiative+linked epics, greppable blockers, reconciliation) intact. |
| `ait-recap/SKILL.md` | `claude/skills/ait-recap/SKILL.md` | **Adapted, no capability lost** | Stripped `allowed-tools`; `version` 0.1.0→1.0.0; added read-only-for-every-seat note; "Find candidate databases with `Glob`"→"by scanning for the tracker files"; claim tip reworded to `PASEO_AGENT_ID`. All recap output templates preserved unchanged. |
| `README.md` | *(no template equivalent — new manifest)* | — | Parity table (lines 30–35) cross-checked against actual diffs and is accurate, including the deliberate-not-copied `.DS_Store` note. |

**No original command, flag, workflow step, or semantic content was silently
lost.** All changes are additive (role matrix, contract precision, identity
canonicalisation) or deletions of Claude-specific metadata (`allowed-tools`,
`author`, `license`) as expected.

---

## (b) Findings by severity

### Critical
None.

### High
None.

### Medium

- **M1 — `hidden_count` omission rule is factually wrong**
  `docs/SKILLS/ait/ait/SKILL.md:234–238` ("`list` and `hidden_count`") says
  "the `hidden_count` field is omitted when nothing is hidden."
  Observed behaviour: the default `list` **always** includes `hidden_count`
  (shows `0` when nothing is hidden, `N` when N closed/cancelled issues are
  filtered); `list --all` omits the field entirely.
  Suggested fix: change to "the `hidden_count` field is always present in the
  default `list` (0 when nothing is hidden) and is omitted under `--all`."
  (Also inherited by the template, so parity is preserved — this is an
  accuracy fix, not a divergence.)

### Low

- **L1 — sibling rekeying on `delete` is undocumented**
  `docs/SKILLS/ait/ait/SKILL.md:108–125` (Delete) and :187–195 (Hierarchical
  IDs) present IDs as stable and "visible directly in the identifier."
  Observed behaviour: `delete <child> --force` rekeys surviving siblings
  (`.2`→`.1`, `.3`→`.2`). A reader who recorded sibling IDs before a delete
  would reference the wrong issue afterwards.
  Suggested fix: add one sentence to the Delete section — "deleting a child
  also rekeys its surviving siblings (e.g. `.2` becomes `.1`), so re-read
  IDs after a delete."

- **L2 — README role table grants Peer `config`, which SKILL.md/protocol do not**
  `docs/SKILLS/ait/README.md:65` lists `config` in the shared-read row with
  all three seats ✅. But `docs/SKILLS/ait/ait/SKILL.md:292` Peer read list
  and `WORKSPACE_PROTOCOL.md` Peer list are
  `show`/`list`/`log`/`export`/`search`/`status` (no `config`); `config` is
  part of the **Supervisor** read subset only. `config` is read-only, so no
  safety impact, but the matrix is internally inconsistent with the protocol.
  Suggested fix: split the README read row per seat (Peer:
  `show`/`list`/`log`/`export`/`search`/`status`; Supervisor:
  `status`/`list --long`/`show`/`search`/`config`/`export`) or drop `config`
  from the Peer column.

### Nit

- **N1 — usage-error exit code 64 not surfaced**
  `docs/SKILLS/ait/ait/SKILL.md:213–214` says errors exit "typically 65 for
  conflicts/validation, 1 for uninitialised." Combined `--human`/`--tree`
  returns exit **64** (code `usage`). The word "typically" makes this not
  wrong, but the contract is precise elsewhere.
  Suggested fix (optional): add "64 for usage errors" alongside 65/1 in the
  Output Modes paragraph and the README "Verified command contract".

---

## (c) Command-contract verification evidence

All commands run against `/usr/local/bin/ait` (version `dev`) with throwaway
DBs under `/tmp/ait-peer-test`. Observed output is verbatim where relevant.

| Skill claim | Command run | Observed (verbatim) | Exit | Verdict |
|---|---|---|---|---|
| Exit **65** `conflict` on claim of a held issue | `ait --db x.db claim peertest-UkLWZ.1 "agent-B"` (after agent-A claimed) | `{"error":{"code":"conflict","message":"issue peertest-UkLWZ.1 is already claimed by agent-A"}}` | 65 | ✅ matches (holder's name included, as documented) |
| Exit **65** `validation` on dep cycle | `ait --db x.db dep add peertest-UkLWZ.2 peertest-UkLWZ.1` (A→B already added) | `{"error":{"code":"validation","message":"adding this dependency would create a cycle"}}` | 65 | ✅ matches |
| Exit **65** `confirmation` on `delete` without `--force` | `ait --db x.db delete peertest-UkLWZ.1` | `{"error":{"code":"confirmation","message":"delete is permanent and unrecorded; it removes the issue along with its notes and dependency links. Re-run with --force to confirm."}}` | 65 | ✅ matches |
| Exit **1** `uninitialised` | `ait --db /tmp/ait-peer-test/x.db list` (fresh DB) | `{"error":{"code":"uninitialised","message":"no ait database at /tmp/ait-peer-test/x.db — run 'ait init' first"}}` | 1 | ✅ matches SKILL.md:249 exactly |
| `--human` + `--tree` mutually exclusive | `ait --db x.db list --human --tree` | `{"error":{"code":"usage","message":"--human and --tree are mutually exclusive"}}` | 64 | ✅ mutually exclusive; usage exit 64 is documented |
| `close --id --note` (and `--reason` alias) | `ait --db x.db close rk-UkLWZ.1 --note "x"`; `ait --db x.db close peertest-UkLWZ --reason "Alias works"` | both close successfully; `--reason` closed epic with note recorded (appears in later export "**Notes:**") | 0 | ✅ both work as documented |
| `create --description @file` | `ait --db x.db create --title "Epic One" --type epic --description @/tmp/ait-peer-test/spec.md` | created; `show` returns full multi-line description "Multi-line description\nfrom file." | 0 | ✅ reads description from file |
| `export <id> --output` | `ait --db x.db export peertest-UkLWZ --output briefing.md` | produced briefing.md: title/ID/P2 header, description, **Notes:** list, `- [ ]` task checklist, `## Summary` counts | 0 | ✅ self-contained markdown briefing |
| `show` returns children, blockers, notes | `ait --db x.db show peertest-UkLWZ` | `"children"`, `"blockers": []`, `"blocks": []`, `"notes": []` keys all present | 0 | ✅ SKILL.md:40 wording ("children, blockers, notes") accurate |
| `status` counts keys | `ait --db x.db status` | `{"counts":{"cancelled":0,"closed":1,"in_progress":0,"open":1,"ready":1,"total":2}}` | 0 | ✅ exact `{cancelled, closed, in_progress, open, ready, total}` |
| `config` reports schema_version 4 | `ait --db x.db config` | `{"prefix":"peertest","schema_version":4}` | 0 | ✅ |
| `init` gitignore + prefix inference | `ait --db gitrepo/.ait/ait.db init --prefix gitproj`; `ait init` (no prefix) | `.ait/` appended to `.gitignore`; no-prefix init inferred `infer-dir` | 0 | ✅ |
| Task directly under initiative rejected | `ait --db x.db create --title "Bad task" --parent <init-id>` | `{"error":{"code":"validation","message":"tasks cannot be direct children of initiatives — create an epic under the initiative first, then add tasks to that epic"}}` | 65 | ✅ matches SKILL.md "Common mistake" |
| `ready --type task` excludes epics | `ait --db x.db ready --type task` | returned only the task, not the epic | 0 | ✅ |
| `delete --force` refuses with children unless `--cascade` | `ait --db x.db delete peertest-UkLWZ --force` | `{"error":{"code":"validation","message":"issue peertest-UkLWZ has 2 descendant issue(s); pass --cascade to delete the whole subtree"}}` | 65 | ✅ |
| `delete --force` response shape | `ait --db x.db delete peertest-UkLWZ.1 --force` | `{"deleted":[{"id":"peertest-UkLWZ.1","title":"Task A","status":"open","type":"task","priority":"P2"}]}` | 0 | ✅ matches `{ "deleted": [refs] }` |
| Errors → stderr, data → stdout | `delete peertest-UkLWZ.2` (missing id) | stdout empty; `{"error":{...}}` captured from stderr | 66 | ✅ envelope on stderr; note `not_found` exits 66 (not in skill's "typically 65" list; harmless) |
| `--long` on mutating command returns full record | `ait --db x.db update <init> --priority P1 --long` | returned full record incl. `type`, `parent_id`, timestamps | 0 | ✅ |
| `list --all` / `hidden_count` | `list` vs `list --all` | default list shows `"hidden_count": 0/1`; `--all` omits the field | 0 | ✅ documented accurately |

Not exercised (outside reviewer's allowed command set; Lead-only): `init` on a
member repo, `flush`, `log purge`, `close --cascade`, `cancel`, `reopen`,
claim-to-assign, `unclaim`. These are documented consistently with
`WORKSPACE_PROTOCOL.md` (see (d)).

### Structural checks (Dimension 1)

- Layout mirrors template exactly: `ait/` (SKILL.md + DELEGATION.md),
  `ait-fleet/SKILL.md`, `ait-recap/SKILL.md`, plus README.md manifest. ✅
- All three `SKILL.md` files have valid YAML frontmatter containing **only**
  `name`, `description`, `version`. ✅ (parsed via PyYAML — keys exactly
  `['name','description','version']`, version `1.0.0`.)
- Grep sweep: no `allowed-tools`, `author`, `license` fields anywhere in the
  bundle (README mentions them only to document the removal). ✅
- Grep sweep: no "pick a fun name"/persona/display-name claim guidance
  remains; claim identity is `$PASEO_AGENT_ID`/`PASEO_AGENT_ID` everywhere. ✅

### Role separation (Dimension 4)

Cross-checked SKILL.md:286–297 and README.md:58–66 against
`WORKSPACE_PROTOCOL.md` (Issue tracking §Role command boundaries):

- Lead-only set matches exactly: `init`, `close --cascade`, `cancel`,
  `flush`, `delete`, `log purge`, `reopen`, claim-to-assign, unclaim of
  others. ✅
- Peer set matches exactly: claim/unclaim own work, `update`, `note add`,
  `close` own task (no cascade), `ready --type task`, read commands
  `show`/`list`/`log`/`export`/`search`/`status`. ✅
- Supervisor set matches exactly: `status`, `list --long`, `show`, `search`,
  `config`, `export`; `log` only to inspect flush history. ✅
- DELEGATION.md reconcile mapping is consistent with the matrix (creation =
  Lead, cascade close = Lead, own-task close = Peer). ✅
- `close --cascade` correctly flagged "Lead only" in SKILL.md:205, 315 and
  fleet SKILL.md:255. ✅
- README role table separates Peer read commands from Supervisor `config`. ✅

### No false auto-issue claims (Dimension 5)

- `ait/SKILL.md:271–278` states the explicit negation only: "`ait` never
  creates issues from conversation, prompts, or session activity. There is no
  hook that turns a user mention into a tracker entry. Issue creation is
  explicit only: someone runs `ait create`."
- `README.md:71–75` states the same negation ("no hook wired in this bundle
  to auto-create them").
- Grep sweep confirms no positive auto-creation claim anywhere in the bundle;
  the only "hook" mentions are the negations. ✅

---

## Follow-up closure

The current bundle now reflects all four previously suggested corrections:

1. Default `list` always documents `hidden_count`; `--all` omits it.
2. `delete` documents sibling rekeying and requires IDs to be re-read.
3. README role permissions keep `config` out of the Peer command set.
4. Usage errors explicitly document exit code 64.

No structural, parity, role, or safety blocking issues remain. The bundle is
ready for Lead's final acceptance.
