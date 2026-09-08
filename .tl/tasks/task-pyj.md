---
id: task-pyj
title: Compact dashboard task metadata into one line
status: done
priority: medium
type: chore
created_at: 2026-09-08T12:52:24Z
updated_at: 2026-09-08T12:53:51Z
created_by: pi-dashboard
assignee: null
depends_on: []
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags: []
references:
  - task-ylk
  - features/list-dashboard.feature
  - cmd/dashboard.go
---

## Description

User feedback: four metadata bullets make dashboard task blocks too tall. Replace status/priority/type/claimant bullets with one compact line, retaining all values and using a Unicode claimant icon. Keep status underscores readable, safely render arbitrary metadata, and leave filters, JSON and other task content unchanged. Update BDD expectations, regression tests and documentation.

## Notes

- 2026-09-08T12:53:51Z [pi-dashboard] note: Replaced four metadata bullets with a single line: code-formatted status, bold priority, code-formatted type, and a Unicode claimant icon plus code-formatted actor. Underscores now remain readable; missing claimant stays -. Updated the existing dashboard BDD expectations first and confirmed two red scenarios before implementation. Added compact-line and code-span regression tests, including backticks, whitespace and control characters; README documents field order. Verification: make bdd passed all 252 scenarios; make test, go vet ./... and git diff --check passed. Smoke-tested go run . list --dashboard --claimed-by pi-dashboard. JSON/filter behavior remains unchanged.
