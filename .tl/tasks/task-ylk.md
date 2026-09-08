---
id: task-ylk
title: Add tl list --dashboard flag for human-readable Markdown overview
status: done
priority: medium
type: task
created_at: 2026-06-04T13:05:07Z
updated_at: 2026-09-08T12:36:33Z
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

## Notes

- 2026-09-08T11:53:59Z [claude] note: Cross-reference from an external adopter (rssb). A grouped overview was wanted there too and this ticket already covers it, so no duplicate was filed. Two notes on scope. 1. Grouping and structure are different axes. The dashboard as described groups a flat list by status or type. The other thing readers ask for is the dependency shape — which slices a feature decomposed into, and which one is blocking. That is filed separately as task-7fi (`tl tree`), since it is a different rendering of different data. If both land, `--dashboard` may want to reuse the tree walker rather than growing its own nesting. 2. The dashboard will want references. In rssb the useful overview joined each ticket to the `.feature` file specifying it, which needs references present in bulk output — currently missing, filed as task-ps0. Worth treating as a prerequisite if the dashboard is meant to link tasks to their artefacts.
- 2026-09-08T12:29:20Z [pi-dashboard] note: Reviewed context and added features/list-dashboard.feature first. Implementing status grouping, compact descriptions/references, deterministic Markdown, existing claim/status/tag filters plus type/priority filters, and JSON precedence. Current Task model has no area/due-date fields; watch is optional and dependency trees belong to task-7fi. Bulk JSON references remain task-ps0; Markdown can read Task.References directly.
- 2026-09-08T12:36:33Z [pi-dashboard] note: Completed the scoped first dashboard after user requested continuation. Added BDD spec first (initial red run confirmed missing --dashboard), then cmd/dashboard.go renderer and list flag/filter integration. Status-grouped Markdown includes ID/title/status/priority/type/claimant, 240-code-point description summaries and references; escapes Markdown/HTML/control characters, omits notes/colors/timestamps, and never mutates ledger data. Added --type/-t and --priority/-p filters for all list formats; existing filters compose, --json takes precedence, and ID tie-breaking stabilizes identical timestamps. README documents usage and scope. Verification: make bdd passed all 247 scenarios; make test, go vet ./..., gofmt check and git diff --check passed. Unit tests cover Unicode truncation/escaping, ordering, JSON compatibility, invalid priorities and read-only snapshots. Smoke-tested go run . list --dashboard --type task --priority medium --tag aur with redirected Markdown. Dedicated area/due-date design preserved in follow-up task-n6n; optional watch is deferred, dependency trees remain task-7fi, bulk JSON references remain task-ps0. No unrelated code or existing ledger changes reverted; no commit made.
