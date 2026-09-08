---
id: task-npt
title: Add --actor to resolve and refine for consistency
status: done
priority: low
type: task
created_at: 2026-09-08T13:35:06Z
updated_at: 2026-09-08T17:54:27Z
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
  - cmd/resolve.go
  - cmd/refine.go
  - features/resolve.feature
  - features/actor.feature
  - features/update-actor.feature
---

## Description

## Problem

Mutating commands are inconsistent about explicit actor attribution. `create`, `claim`, `cancel`, `pending`, `note`, `close`, `block` and `release` accept `--actor`. `resolve` and `refine` do not — they resolve the actor only from TL_ACTOR / ACTOR_NAME / BEADS_ACTOR or auto-detection.

Found while dogfooding: recording a human decision with `tl resolve` required setting an environment variable inline, while the surrounding `cancel` and `close` calls in the same session took `--actor` directly. CLAUDE.md instructs agents to 'use an explicit actor on mutating commands', which these two commands cannot honour.

`resolve` is the sharper case of the two: it records a human's answer, so attribution is the point of the event, and getting it wrong silently credits the answer to whichever actor the environment happens to resolve to.

## Proposed solution

Add the standard `--actor` flag to both commands, with the same precedence as everywhere else: flag beats TL_ACTOR beats ACTOR_NAME beats BEADS_ACTOR beats auto-detection.

## Design constraints

- Purely additive. Omitting the flag keeps today's resolution order exactly.
- Reuse the existing ResolveActor path rather than adding a second precedence rule.

## Acceptance

features/resolve.feature and features/refine.feature gain a scenario each asserting the recorded event carries the actor passed on the flag.

## Notes

- 2026-09-08T17:52:56Z [pi-actor] note: Claimed and reviewed all references. Added features/update-actor.feature first for fallback/no-op/reference attribution and added explicit-actor scenarios to resolve/refine plus editor coverage. Existing code inspection found refine events are currently anonymous, including reference changes; events.Append does not resolve actors. Implementing the requested shared ResolveActor fallback for these paths as well as explicit flags. Initial BDD red run confirms missing flags and missing editor attribution.
- 2026-09-08T17:54:27Z [pi-actor] note: Implemented --actor on resolve and refine via the shared ResolveActor chain: flag > TL_ACTOR > ACTOR_NAME > BEADS_ACTOR > detection. Resolution attributes both the pending_resolved event and answer note. Refinement attributes regular field updates, editor updates, and reference-added/removed events; no-op and validation behavior preserved. Refine previously emitted anonymous events, now receives the requested standard fallback. Added the new BDD feature first plus required resolve/refine scenarios and editor coverage; initial red tests confirmed missing flags/attribution. Updated README and usage examples. Verification: make bdd passed all 306 scenarios; make test, go vet ./..., gofmt and git diff --check passed. Built-binary smoke test in a temporary ledger verified --actor aho overriding TL_ACTOR in JSON-mode resolve/refine, reference events and answer note. Dogfooded the new refine flag to attach features/update-actor.feature to this task. No commit made.
