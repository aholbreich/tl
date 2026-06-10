---
id: task-wke
title: Add --type field for semantic task categorization
status: open
priority: medium
type: task
created_at: 2026-06-04T13:05:26Z
updated_at: 2026-06-04T13:05:26Z
created_by: human
assignee: null
depends_on: []
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags: []
---

## Description

## Problem

tl currently treats every entry as a generic task. But in practice, unresolved work items fall into distinct semantic categories:

- **task** — concrete, actionable, can be claimed and completed
- **decision** — a choice or scope boundary that needs resolution
- **research** — investigation, evaluation, learning that produces knowledge, not code
- **waiting** — blocked on an external person, system, or date
- **risk** — security, legal, compliance, or operational risk to monitor
- **question** — unknown fact or context that needs clarification

Each category has different coordination needs. For example, a decision doesn't need claiming — it needs a resolution and a record. A risk doesn't get closed — it gets monitored or mitigated. A question doesn't need a lease — it needs an answer.

Without a type field, tl consumers must infer semantics from tags, descriptions, or naming conventions. This works at small scale but breaks down when multiple agents and humans coordinate across different work categories.

## Proposed solution

Add an optional --type flag to tl create and tl refine, supporting these values:

| Type | Coordination behavior | Claim needed? |
|---|---|---|
| task | Standard task workflow | Yes |
| decision | Record options + resolution | No (lightweight) |
| research | Track investigation scope + findings | Optional |
| waiting | Track blocker + expected resolution | No |
| risk | Monitor + severity tracking | No |
| question | Track answer + resolution | No |

Non-task types skip the claim/stale/lease machinery. tl list and tl ready can filter by type. For example:
- tl ready → shows only task-type items (actionable)
- tl list --type decision → shows unresolved decisions for review
- tl list --type research → shows active investigations

The --dashboard output groups tasks by type, giving readers immediate context about what kind of work each item represents.

## Design constraints

- Backward compatible — existing tasks without --type default to task
- Type is stored in the YAML frontmatter of the task file
- tl list supports --type filtering with comma-separated values
- tl ready excludes non-task types by default (only actionable work is ready)
- Type transitions are tracked in the event log (e.g., from research to task when investigation is done and implementation begins)

## Use case context

A project maintains a single overview of unresolved work: actionable tasks, pending decisions, active research, open questions, monitored risks. Currently this lives in a single Markdown file. The project wants:

1. Agent-safe coordination for actionable tasks (claims, leases) — tl's existing strength
2. Lightweight tracking for non-actionable items (decisions, research, questions) without coordination overhead
3. A single tool that handles both, so humans and agents don't switch between a task tracker and a notes file
