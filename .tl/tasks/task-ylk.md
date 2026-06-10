---
id: task-ylk
title: Add tl list --dashboard flag for human-readable Markdown overview
status: open
priority: medium
type: task
created_at: 2026-06-04T13:05:07Z
updated_at: 2026-06-04T13:05:07Z
created_by: human
assignee: null
depends_on: []
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags: []
references:
  - docs/comparison.md
---

## Description

## Problem

tl stores tasks as individual Markdown files in .tl/tasks/. This is clean for version control (git diff shows exactly what changed) but removes the single-page overview that many teams rely on. For example, LLM-curated wikis and knowledge bases often maintain an open-loops or tasks page as a single rendered Markdown file that serves as a dashboard — a human can scan all active work in one view, agents can ingest it in one read, and git history shows the evolution of the whole work set.

Currently, to get this overview from tl, a human or agent must run tl list, then read each task file individually. This breaks the single-page dashboard workflow.

## Proposed solution

tl list --dashboard renders all tasks (or a filtered subset) as a single, well-structured Markdown block:

- Each task is a compact block: title, status, priority, claimant, due date, short description
- Grouped by status or type (open / in-progress / blocked / pending)
- Output goes to stdout — can be redirected to a file for committing to a wiki or knowledge base
- Supports --type and --area filters to select which tasks appear
- Optional --watch mode could regenerate on file changes for local preview

## Use case context

A project maintains a curated Markdown knowledge base. Some unresolved items are actionable tasks (fix a CI failure, evaluate a tool), others are decisions, research questions, or waiting states. The project wants:

1. Agent-safe task coordination (claims, leases, stale detection) — tl already provides this
2. A single human-readable dashboard page — currently missing
3. The dashboard should be regeneratable (not manually maintained) so it never drifts from the task ledger

## Design constraints

- Output is valid Markdown suitable for direct inclusion in a wiki page
- Works with --json for programmatic consumers AND as plain Markdown for humans
- Does not introduce a new storage format — tasks remain in .tl/tasks/
- Respects existing filters (--type, --priority, --status) to produce scoped dashboards
