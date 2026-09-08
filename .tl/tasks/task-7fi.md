---
id: task-7fi
title: Add tl tree to render the dependency graph
status: open
priority: medium
type: task
created_at: 2026-09-08T11:52:16Z
updated_at: 2026-09-08T11:52:16Z
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
  - features/dep-add.feature
  - cmd/dep.go
  - task-ylk
---

## Description

## Problem

"Dependency-aware — agents can ask what is ready now" is thesis point 4 of the PRD, and tl already stores and computes the DAG: `depends_on` is in the frontmatter, in JSON, and drives `tl ready`.

But the graph is write-only from a human standpoint. `tl show` lists a task direct dependencies as a flat list; `tl list` shows none. To understand how a feature decomposes — which slices exist, which are done, what is actually blocking the parent — a reader must run `tl show` repeatedly and assemble the tree by hand.

tl is the only tool that can render this. No external consumer has the DAG.

## Proposed solution

```
tl tree [TASK_ID] [--all] [--json]

task-wk5  Add RSS/Atom sources through TUI and CLI      open
├─ task-ckm  Subscribe to a direct RSS/Atom URL         done
│  └─ task-7nl  Fetch and parse real RSS/Atom feeds     done
├─ task-efl  Add-source form in the TUI                 open
└─ task-nx5  Discover advertised feeds from a website   open
```

- With `TASK_ID`: the subtree rooted at that task.
- Without: every root (a task nothing else depends on), so the whole ledger reads as a forest.
- Closed tasks hidden by default and revealed by `--all`, matching `tl list`.
- Colour follows the existing list conventions (dim for closed, priority colours) via commandColorEnabled.

## Design constraints

- Read-only; no new storage, no new field. Purely a rendering of data tl already has.
- Cycles must render without hanging — print the repeated node marked as a cycle rather than recursing. `tl doctor` owns diagnosing them; `tl tree` must merely survive them.
- A task appearing under two parents is legitimate (a shared dependency). Render it in both places rather than picking one.
- `--json` emits nested children so consumers do not re-derive the graph.

## Relationship to task-ylk

task-ylk (`tl list --dashboard`) is a flat Markdown overview grouped by status or type. This is the orthogonal axis: structure rather than grouping. They should stay separate commands, but `--dashboard` may want to reuse the tree walker once it exists.

## Use case context

Reported from rssb, where one feature ticket was split into three vertical slices with a dependency chain. The shape of that decomposition — and which slice is blocking — is currently only visible by reading four tasks one at a time.
