---
id: task-thh
title: Split scenario anchors before validating references in doctor
status: cancelled
priority: high
type: task
created_at: 2026-09-08T13:14:17Z
updated_at: 2026-09-08T13:34:37Z
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
  - internal/doctor/doctor.go
  - features/doctor.feature
  - features/references.feature
  - task-kd0
---

## Description

## Problem

`checkReferences` (internal/doctor/doctor.go:396) classifies any reference containing `/` as a repo-relative path and stats it. A scenario-anchored spec reference such as `features/login.feature#Signing in with a valid password` contains `/`, so the stat fails and the reference is reported dead — with `Fixable: true` and fix kind `fixDeadRef`.

`tl doctor --fix` therefore **deletes every anchored spec reference in the ledger**, silently, on its first run.

No anchored reference exists yet, so this is not a live bug. It is a prerequisite: any ticket that introduces anchor syntax must land after this one, or adopting the convention destroys the links it just created.

## Proposed solution

Split a path-shaped reference at the first `#` before validating, and stat only the path part.

- A missing file part is still a dead reference — unchanged behaviour.
- The anchor is not validated here. Whether the named scenario actually resolves requires reading the file, which is a separate question gated on decision 0002.
- References with no `#` behave exactly as they do today.

## Design constraints

- **Fragment splitting applies only to path-shaped references.** URL-shaped references already skip this check and must keep their fragments intact — a `#` in a URL is not a scenario anchor.
- **No file reading.** This ticket is solely about not destroying data; it deliberately adds no knowledge of Gherkin.
- **The fix must remove the whole reference**, anchor included, not just the path part.
- A file path may legitimately contain `#`. Splitting on the first `#` makes such a path unaddressable; accept that and note it, rather than adding escaping.

## Acceptance

`features/doctor.feature` gains scenarios covering: an anchored reference whose file exists is clean; an anchored reference whose file is missing is reported dead; `--fix` removes the entire anchored string rather than leaving a fragment behind.

## Why this is first

It is small, needs no decision, and is the one item on this path that can cause data loss. Landing it early means the anchor convention can be adopted incrementally afterwards without a flag day.

## Notes

- 2026-09-08T13:34:37Z [claude] cancelled: Decision 0002 rejected scenario anchors in references (alternative 5), so no reference will contain a '#' and doctor has no anchor to mis-parse. The underlying hazard is real but latent: checkReferences treats anything containing '/' as a path, so a hand-written 'features/x.feature#Name' is still reported dead and 'tl doctor --fix' would delete it. Revive this ticket only if anchors are ever adopted as a second resolver.
