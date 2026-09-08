---
id: task-1re
title: Should tl read referenced spec files at display time?
status: pending_human
priority: medium
type: decision
created_at: 2026-09-08T11:52:57Z
updated_at: 2026-09-08T13:15:50Z
created_by: claude
assignee: null
depends_on: []
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
pending:
  question: 'Decision 0002 is drafted at .decisions/0002-reading-referenced-files.md, proposing that read commands MAY open referenced files under four limits: the reference is the trigger (no project-level Gherkin detection); existence is a cheap stat but content parsing happens only for anchored refs or behind a flag; enrichment is automatic in human output but schema-stable in JSON; and tl reports structure only, never deciding what a tag means. Accept as drafted, or fall back to alternative 2 (doctor only)? Three open questions inside need your call too: whether a missing referenced file reads as ''no spec'' or ''broken link'', what tl tree renders for an unresolved spec, and which resolver is second.'
  requester: claude
  requested_at: 2026-09-08T13:15:50Z
tags: []
references:
  - docs/PRD.md
  - internal/doctor/doctor.go
  - task-kd0
  - task-ahk
  - .decisions/0002-reading-referenced-files.md
  - task-thh
---

## Description

## The question

Once tl recognises a `.feature` reference as a spec (task-kd0), the obvious next step is to report something *about* that spec — for example, that the file carries an `@implemented` tag, or how many scenarios it holds.

That requires tl to open and parse a file outside `.tl/`. This ticket is to decide whether that is acceptable, before task-ahk is built on the assumption that it is.

## Why it is not obviously fine

PRD thesis point 7 is "Small and predictable — no daemon, no hidden database, no automatic remote push". Reading arbitrary repo files at display time weakens two properties:

- **Predictability.** `tl list` output would depend on working-tree state outside the ledger. The same ledger at the same commit could render differently with a dirty tree.
- **Cost and failure modes.** Every listed task potentially becomes a file read. Missing, unreadable, or malformed files need defined behaviour rather than an error that breaks the whole listing.

## Why it may be fine anyway

- **`tl doctor` already does it.** checkReferences stats referenced paths today (internal/doctor/doctor.go:396). The boundary has been crossed; the question is only whether it extends from `doctor` to the read commands.
- **Repo-local and deterministic.** This is a file read in the same working tree, not a network call, a subprocess, or a daemon. It is closer to the existing `doctor` behaviour than to anything on the non-goals list.
- **The alternative is worse.** The other way to know a spec status is to store it in the ledger — derived state that goes stale silently the moment someone edits the feature file. That contradicts "Human-readable" and creates exactly the drift tl exists to prevent.

## Options

1. **Never.** Recognition only (task-kd0). Consumers join spec status themselves from `--json` plus their own test runner. tl stays minimal.
2. **`doctor` only.** Extend the existing check to validate spec references more deeply; read commands stay pure. Conservative, and `doctor` is already the "look at the world" command.
3. **Opt-in on read commands.** Reads happen only behind an explicit flag (`tl list --spec-status`), never on the default path. Default output stays pure and cheap; the cost is paid only when asked for.
4. **Always.** Spec status appears wherever a spec reference does.

## Recommendation

Option 3, with option 2 as the fallback if that still feels too invasive. It preserves the default-path guarantee — plain `tl list` reads nothing but the ledger — while making the joined view a single command rather than a scripting exercise. Option 4 should be rejected: it makes every listing pay for a feature most projects will not use.

Whatever is chosen, record it in `.decisions/` alongside 0001, since it sets a precedent for how far tl reaches outside its own directory.
