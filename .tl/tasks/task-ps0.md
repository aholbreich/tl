---
id: task-ps0
title: Include references in list and ready JSON output
status: done
priority: high
type: task
created_at: 2026-09-08T11:51:58Z
updated_at: 2026-09-08T12:38:26Z
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
  - cmd/json.go
  - features/list.feature
  - features/references.feature
  - docs/PRD.md
  - task-ir9
---

## Description

## Problem

`tl show --json` emits `references`. `tl list --json` and `tl ready --json` do not — the field is silently absent, not empty.

This breaks the PRD thesis "Machine-readable — every read command supports `--json`": the flag is there, but bulk reads lose a field that single reads keep. Any consumer wanting to answer "which tasks reference this file?" must currently shell out to `tl show` once per task, which is what stops such tooling from being written at all.

## Cause

`compactTasksJSON` in cmd/json.go:28 copies ID, Title, Status, Priority, Type, CreatedAt, UpdatedAt, CreatedBy, Assignee, DependsOn, Claim, Pending, Tags, Description, Notes and Sections — but not References, although `task.Task.References` carries a `json:"references"` tag (internal/task/task.go:35).

This looks like collateral rather than intent. task-ir9 introduced the compact DTO to omit the *raw body* while explicitly preserving parsed description and notes; references simply never got carried across.

## Proposed solution

Add `References` to `compactTaskJSON`, emitting `[]` (not null) when empty, matching the shape `tl show --json` already produces.

## Design constraints

- Additive and backward compatible: consumers reading the current shape keep working.
- Same JSON shape as `tl show --json` for the same field — no second serialization of the same concept.
- Empty is `[]`, consistent with the existing "JSON output emits an empty references array" scenario in features/references.feature.

## Acceptance

features/list.feature and features/ready.feature gain a scenario asserting references survive bulk JSON output, mirroring the existing show scenario.

## Why this is first

Every other reference-aware idea (spec links, dashboards, cross-artefact views) is gated on this. Without it those tools cost one subprocess per task; with it they are a single command plus a filter.

## Use case context

Reported from rssb, which uses tl to track features specified as Gherkin. The goal is an overview joining ticket status with the `.feature` file each ticket specifies. That join needs references in a bulk read.

## Notes

- 2026-09-08T12:38:26Z [claude] note: Implemented. Added References to compactTaskJSON (cmd/json.go), so tl list --json and tl ready --json now emit the same references field tl show --json does. Found and fixed a second-order gap: store.List did not normalize nil References the way store.Read does (yaml omitempty), so a naive DTO copy would have emitted null instead of []. Normalized in store.List for symmetry with Read; compactReferences() is a belt-and-braces guard at the JSON boundary. Tests: 4 new scenarios, two per file in features/list.feature and features/ready.feature, mirroring how references.feature splits the populated and empty cases. New array-scoped step defs in bdd/bulk_json_test.go fail on absent field, null, and non-array alike. Verified negatively: reverting the one-line DTO change fails exactly those 4 and no others. Suite 251/251, go vet clean. Scope decision: references were NOT added to the plain tl list / tl ready tables. References are multi-valued (up to 5 here, mixing paths and URLs) and would wrap or overflow a one-line-per-task row; the tables are the what-should-I-pick-up scan, and tl show already carries the context you want after picking. tl list --dashboard already renders references (cmd/dashboard.go:33), so JSON was the only real gap. If a table view is wanted later it should be an opt-in --refs flag with its own ticket. Unblocks task-kd0.
