---
id: task-n6n
title: Decide whether tasks need dedicated area and due-date fields
status: open
priority: low
type: decision
created_at: 2026-09-08T12:36:25Z
updated_at: 2026-09-08T12:36:25Z
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
  - internal/task/task.go
  - features/list-dashboard.feature
---

## Description

Follow-up from task-ylk: the dashboard proposal mentioned --area and due dates, but neither exists in the current task model. The first dashboard uses existing --tag scoping and omits due dates, without introducing new storage fields. Decide whether dedicated metadata is needed or tags/no dates are sufficient. If approved, define area versus tag semantics, due-date format and timezone/overdue behavior, creation/refinement and filter interfaces, legacy-task defaults, and JSON/dashboard representation before scheduling implementation. Record the decision; do not silently add schema fields as part of rendering.
