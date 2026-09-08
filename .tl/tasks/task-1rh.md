---
id: task-1rh
title: tree --spec-status resolves specs then discards them on the --json path
status: done
priority: high
type: task
created_at: 2026-09-08T21:53:59Z
updated_at: 2026-09-08T22:22:54Z
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
  - internal/spec/spec.go
  - .decisions/0002-reading-referenced-files.md
  - task-ps0
  - features/tree.feature
---

## Description

## Problem

`tl tree --spec-status --json` reports no spec information. The flag is accepted, the files are read, and the result is thrown away — so the caller pays the cost of the file reads and receives nothing for it.

The text view is unaffected and correct:

```
task-6y5   Organize sources into topics   open  high  features/topics.feature (15), features/topic_sidebar.feature (11)
```

The same command with `--json` emits only id, title, status, priority, cycle and children.

This is the same class of gap as task-ps0 (references missing from bulk JSON), and the same PRD thesis is at stake: "Machine-readable — every read command supports `--json`". A flag whose entire purpose is to surface derived data should not surface it in the human view only, since the derived view is precisely what a dashboard or an external overview wants to consume.

## Cause

cmd/tree.go:47 computes `specs := resolveSpecsFor(ledger, tasks, specStatus)` before the output branch. The JSON branch at :49-52 calls `treeForestJSON(forest)`, which takes only the forest — `specs` is never passed. The text branch at :59 passes it to renderForest.

`treeNodeJSON` (cmd/tree.go:171) has no spec field.

The work is therefore performed and discarded: with `--spec-status --json` the referenced files are opened and parsed, and the output is byte-identical to plain `--json`.

## Why it looks like an oversight rather than a decision

`spec.Spec` already carries JSON tags on every field — `json:"path"`, `json:"state"`, `json:"scenarios"` (internal/spec/spec.go). The type was written to be serialized. Nothing consumes that serialization today.

## Proposed solution

Give `treeNodeJSON` an optional `specs` array, populated only when `--spec-status` is set and omitted entirely otherwise, so the default JSON shape does not change:

```json
{
  "id": "task-6y5",
  "title": "Organize sources into topics",
  "status": "open",
  "priority": "high",
  "specs": [
    {"path": "features/topics.feature", "state": "present", "scenarios": 15},
    {"path": "features/topic_sidebar.feature", "state": "present", "scenarios": 11}
  ],
  "children": []
}
```

Serialize `spec.Spec` directly rather than restating its fields, so the two views cannot drift and `state` keeps reporting missing and unknown as the text view does.

## Design constraints

- Additive: without `--spec-status` the JSON shape is unchanged and no file is read, preserving the default-path guarantee in decision 0002.
- `state` must survive into JSON. A consumer needs to distinguish "spec present with 0 scenarios" from "spec file missing" — collapsing both to an absent entry would hide exactly the drift this feature exists to reveal.
- No new interpretation. Per decision 0002 and the internal/spec doc comment, nothing concludes that work is done; the JSON reports observations only.

## Alternative considered

Rejecting `--spec-status` together with `--json` as an unsupported combination would at least be honest, but it removes the only machine-readable path to the data and makes the flag useless to the tooling most likely to want it. Not recommended.

## Acceptance

features/tree.feature gains scenarios for: spec data present in JSON under `--spec-status`; the field absent without the flag; and a missing spec file surfacing as `"state": "missing"` rather than vanishing.

## Use case context

Found in rssb while linking a four-ticket feature to three `.feature` files. The text view answered the question immediately; building any regenerated overview on top of it is currently blocked on re-implementing the resolution externally.

## Notes

- 2026-09-08T22:22:54Z [claude] note: Fixed. treeNodeJSON now carries a spec array and treeForestJSON threads the resolved map through every node, so specs appear at any depth rather than being resolved and discarded. Always an array, never null, matching list and ready, so a consumer walking the tree needs no presence check at any level. The report was accurate on cause and on intent: spec.Spec already carried JSON tags on every field, so the type was written to be serialized and nothing consumed it. Landed alongside the removal of --spec-status, which made the waste unconditional rather than opt-in and so raised the priority of fixing it. Three scenarios added: spec on a root, spec on a nested dependency, and an empty array for a task with no spec reference. Mutation-checked — hardcoding the node's spec to an empty slice fails two of them. Suite 308/308.
