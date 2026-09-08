---
id: task-npt
title: Add --actor to resolve and refine for consistency
status: open
priority: low
type: task
created_at: 2026-09-08T13:35:06Z
updated_at: 2026-09-08T13:35:06Z
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
