---
id: task-npp
title: 'Finish tl tree presentation: priority colours and --spec-status'
status: done
priority: medium
type: task
created_at: 2026-09-08T17:22:52Z
updated_at: 2026-09-08T17:25:12Z
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
  - cmd/tree.go
  - cmd/color.go
  - features/tree.feature
  - task-7fi
  - task-kd0
---

## Description

## Problem

Two gaps left by task-7fi and task-kd0, both reported from use.

**Colours.** task-7fi specified 'colour follows the existing list conventions (dim for closed, priority colours)'. Only the dim-for-closed half was built, and it is invisible on the default path because closed tasks are hidden unless --all. So `tl tree` renders monochrome in practice. There is also nothing to colour: the tree has no priority column, unlike `tl list`.

**Spec status.** `tl list --spec-status` and `tl ready --spec-status` resolve referenced .feature specs, but `tl tree` does not accept the flag. A tree is where the question 'which slice is blocking, and is it specified?' is most naturally asked, so the omission is felt more here than in a flat list.

## Proposed solution

- Add a Priority column to `tl tree` and colour it with colorListPriority, exactly as `tl list` does, so the two views read the same way.
- Keep dimming closed rows under --all.
- Accept --spec-status on `tl tree`, reusing resolveSpecsFor and specCell unchanged.

## Design constraints

- Colour must be applied after the tabwriter has flushed. ANSI escapes counted as column width misalign every row; cmd/list.go and cmd/tree.go both already do it in this order.
- --spec-status stays the only path that opens a file outside .tl/ (decision 0002). Plain `tl tree` must read the ledger and nothing else.
- No new colour vocabulary. Reuse the existing helpers so list, show and tree cannot drift apart.

## Acceptance

features/tree.feature gains scenarios for the priority column, for forced colour on a priority value, and for the spec column under --spec-status.

## Notes

- 2026-09-08T17:25:12Z [claude] note: Implemented. tl tree gained a Priority column coloured with colorListPriority, so high/medium/low read red/yellow/blue exactly as in tl list; closed rows still dim under --all, and the two effects compose the same way colorListRow composes them. Also added --spec-status, reusing resolveSpecsFor and specCell unchanged, so the flag behaves identically across list, ready and tree. One thing that differs from list and is worth knowing: colorListRow locates the priority cell by finding the status first, which is safe there because status precedes title in that table. A tree row carries the title BEFORE the status, so the same logic would mistake a task titled 'open the door' for the status column. treePriorityStart anchors on the exact title first, then status, then priority. Verified by hand against a task titled 'open the door and close it' with priority high: colour lands on high. Colour is still applied only after the tabwriter has flushed, since ANSI escapes counted as column width misalign every row. Tests: 6 new scenarios in features/tree.feature covering the priority column, forced-colour priority values, dimmed closed rows, spec resolution at root and at depth, and that plain tl tree resolves nothing. Suite 283/283, go vet clean. README example updated — it previously showed the pre-priority output.
