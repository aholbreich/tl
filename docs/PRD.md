# PRD: TaskLedger

**TaskLedger** (`tl`) is a Git-native task ledger for humans and AI coding
agents.

This document captures the design intent and the parts of the spec that are
not derivable from the code or the Gherkin features. For per-command
behavior, read the corresponding file in [`features/`](../features). For data
shapes and on-disk format, read [`internal/task/task.go`](../internal/task/task.go)
and the create-command output.

---

## 1. Summary

TaskLedger stores tasks as readable Markdown files with YAML frontmatter
inside the repository and exposes a CLI with safe claim leases,
dependency-aware ready queues, handoff notes, and JSON-first output for
automation.

It is **not** a Jira / Linear / GitHub Issues replacement, not an
orchestration platform, not a daemon, not a hosted service.

---

## 2. Problem

Coding agents lose context across sessions. They need a durable, repo-local
way to know:

- What work exists.
- What is blocked.
- What is ready.
- Who is already working on what.
- What was already attempted.
- What counts as done.
- How to safely hand work back to a human or another agent.

Existing options trade off poorly:

- Built-in agent task lists are session-local.
- Plain Markdown TODO files lack dependency and claim semantics.
- Issue trackers are external and not optimized for agent automation.
- Heavier agent-coordination systems introduce databases, daemons, sync
  complexity, or hidden state.

---

## 3. Target users

Primary: a developer or architect using Claude Code, Codex, Cursor, or
similar agents inside a Git repository.

Secondary: indie hackers, small product teams, technical founders,
engineering managers coordinating agent-assisted work, and developers who
want local-first task tracking without SaaS overhead.

---

## 4. Core product thesis

A good agent task system should be:

1. **Repo-local** — state lives in the Git repository.
2. **Human-readable** — a developer can edit tasks with a normal editor.
3. **Machine-readable** — every read command supports `--json`.
4. **Dependency-aware** — agents can ask "what is ready now?".
5. **Claim-safe** — claims have leases with stale-claim detection.
6. **Handoff-oriented** — notes preserve context between humans and agents.
7. **Small and predictable** — no daemon, no hidden database, no automatic
   remote push, no AGENTS.md magic.

---

## 5. Non-goals

TaskLedger will not initially support:

- Real-time collaboration.
- A web app or hosted backend.
- Complex role hierarchies.
- Multi-repository orchestration.
- Automatic Git pushing or merging.
- Long-running background workers.
- tmux/session management.
- A full Jira / Linear / GitHub Issues replacement.
- AI agent execution itself.

The tool tracks and coordinates work. It does not run the agent in v1.

---

## 6. Command surface

Each `.feature` file in [`features/`](../features) is the V1 acceptance test
for one command. Implementation status is tracked by the `@implemented` tag
on each feature.

| Command           | Spec                                                          |
| ----------------- | ------------------------------------------------------------- |
| `tl init`         | [features/init.feature](../features/init.feature)             |
| `tl create`       | [features/create.feature](../features/create.feature)         |
| `tl list`         | [features/list.feature](../features/list.feature)             |
| `tl show`         | [features/show.feature](../features/show.feature)             |
| `tl ready`        | [features/ready.feature](../features/ready.feature)           |
| `tl dep add`      | [features/dep-add.feature](../features/dep-add.feature)       |
| `tl dep remove`   | [features/dep-remove.feature](../features/dep-remove.feature) |
| `tl claim`        | [features/claim.feature](../features/claim.feature)           |
| `tl release`      | [features/release.feature](../features/release.feature)       |
| `tl stale`        | [features/stale.feature](../features/stale.feature)           |
| `tl note`         | [features/note.feature](../features/note.feature)             |
| `tl close`        | [features/close.feature](../features/close.feature)           |
| `tl pending`      | [features/pending.feature](../features/pending.feature)       |
| `tl resolve`      | [features/resolve.feature](../features/resolve.feature)       |
| `tl agents`       | [features/agents.feature](../features/agents.feature)         |

User-visible flag reference lives in the [README](../README.md).

---

## 7. Storage layout

```
.taskledger/
  config.yaml      # defaults (claim TTL, id prefix, …)
  tasks/
    task-<3>.md    # one Markdown+YAML file per task
  events.jsonl     # append-only audit journal
```

Task schema: [`internal/task/task.go`](../internal/task/task.go) (`Task`
struct and `MarshalMarkdown`). A concrete example file is what `tl create`
writes to `tasks/`.

---

## 8. Status model

| Status          | Meaning                                                  |
| --------------- | -------------------------------------------------------- |
| `open`          | Work exists but is not currently claimed.                |
| `in_progress`   | Work is claimed by an actor.                             |
| `blocked`       | Work cannot continue (dependency or technical blocker).  |
| `pending_human` | Work requires human clarification or decision.           |
| `done`          | Work is completed.                                       |
| `cancelled`     | Work is intentionally abandoned.                         |

A task is **ready** (eligible for `tl ready` / `tl claim`) when:

- status is `open`, **and**
- all dependencies are `done`, **and**
- it has no active claim, **and**
- it is not `pending_human`, **and**
- it is not `blocked`.

---

## 9. Claims and events

- **Claims** are first-class: actor, claimed_at, expires_at, heartbeat_at.
  Default lease 60 minutes. Stale claims (past `expires_at`) require
  `--force` to release. Spec: [`features/claim.feature`](../features/claim.feature),
  [`features/release.feature`](../features/release.feature),
  [`features/stale.feature`](../features/stale.feature).
- **Events.** Every mutating command appends a JSON line to
  `.taskledger/events.jsonl`. The task file is the current state; the event
  journal is the audit trail.

---

## 10. Exit codes

| Code | Meaning              |
| ---: | -------------------- |
|    0 | Success              |
|    1 | Generic error        |
|    2 | Invalid arguments    |
|    3 | Task not found       |
|    4 | Task not ready       |
|    5 | Task already claimed |
|    7 | Lock failed          |

---

## 11. Implementation notes

- **Language:** Go 1.25+; single static binary, distributed via `make dists`.
- **Major libs:** `spf13/cobra` (CLI), `gopkg.in/yaml.v3` (frontmatter),
  `cucumber/godog` (BDD acceptance tests).
- **ID generation:** `task-<3 lowercase alphanumeric>`, generated with
  `crypto/rand` and a collision-retry loop. Namespace ≈ 47k; well above the
  realistic ceiling for the project sizes TaskLedger targets.
- **Atomic writes:** task files write to `<id>.md.tmp` and `rename` over the
  target.
- **Locking:** a repo-local `.taskledger/.lock` will guard mutating commands.
  Not yet implemented.
- **Repository detection:** commands walk upward from CWD to find
  `.taskledger/`.

---

## 12. Success metrics

MVP is successful if:

- A developer can initialize a repo and create tasks in under one minute.
- Claude Code, Codex, or similar agents can reliably query ready tasks via
  JSON.
- Two agents do not accidentally work on the same task when claims are
  respected.
- Stale claims can be detected and released.
- All task state is reviewable in Git.
- No hidden service or database is required.
- The project remains understandable after one hour of code reading.

---

## 13. Key differentiator

The strongest differentiator is not "file-based tasks". It is:

> **Agent-safe task coordination with readable Git-native state.**

That means:

- Claims are explicit.
- Stale work is detectable.
- Dependencies are computable.
- Handoffs are recorded.
- Humans can inspect everything.
