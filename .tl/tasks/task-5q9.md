---
id: task-5q9
title: Add low-noise tl ledger branding and discovery
status: done
priority: medium
type: feature
created_at: 2026-09-03T22:16:30Z
updated_at: 2026-09-03T22:25:12Z
created_by: pi
assignee: null
depends_on: []
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags:
  - ux
  - ledger
references:
  - internal/repo/repo.go
  - internal/doctor/doctor.go
  - features/init.feature
  - features/doctor.feature
---

## Description

Brand tl's on-disk ledger format without adding noise to individual tasks or events. New ledgers should include a human-facing .tl/README.md and a machine-readable format: tl marker in config.yaml. Existing unbranded ledgers remain valid; tl doctor should report the missing identity metadata and --fix should add it explicitly.

## Notes

- 2026-09-03T22:25:12Z [pi] note: Implemented v1 ledger identity: new ledgers get .tl/README.md plus a commented format: tl config marker; older configs remain accepted; doctor reports missing identity as fixable warnings, rejects conflicting format markers, and --fix safely adds both without overwriting an existing guide. Updated BDD specs, unit tests, and storage docs. Verified on this repository by running the new doctor --fix. Tests: make test (235 BDD scenarios / 1212 steps plus all unit tests), go vet ./..., git diff --check — all pass with Go 1.25.9.
