---
id: task-1re
title: Should tl read referenced spec files at display time?
status: done
priority: medium
type: decision
created_at: 2026-09-08T11:52:57Z
updated_at: 2026-09-08T13:34:55Z
created_by: claude
assignee: null
depends_on: []
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
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

## Notes

- 2026-09-08T13:34:55Z [human] resolved: Alternative 3, narrowed. Read commands may open referenced files ONLY behind an explicit flag (tl list --spec-status); the default path reads nothing outside .tl/. Linking stays at FILE granularity. Scenario anchors in references were rejected: an exact match on a human-written title breaks silently on rename, and the new syntax would be destroyed by tl doctor --fix. Ledger IDs tagged on scenarios were rejected: the specification is the durable artefact and the ledger is working state, so embedded ticket IDs outlive their meaning. Accepted ceiling: tl reports that a story points at a spec file which exists and carries certain tags, never that this story's behaviour is specified. JSON keeps a schema-stable 'spec' key, always present, null when the flag was not passed. Recorded as decision 0002.
