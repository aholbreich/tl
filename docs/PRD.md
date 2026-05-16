# PRD: Lightweight Git-Native Task Ledger for Coding Agents

## 1. Product Name

Working name: **TaskLedger**

For now, use **TaskLedger** in code and documentation.

---

## 2. Summary

TaskLedger is a lightweight, Git-native task management CLI for humans and AI coding agents.

It stores tasks as readable Markdown files with YAML frontmatter inside the repository. It provides dependency-aware ready queues, safe task claiming, stale claim detection, handoff notes, verification commands, and JSON-first output for automation.

The goal is to capture the useful coordination ideas from tools like Beads while keeping the operating model simple, transparent, and repo-local.

TaskLedger is **not** an orchestration platform, not a daemon, not a database server, and not a replacement for Jira, Linear, GitHub Issues, or Trello.

It is a small, reliable work ledger for one repository and one to five active humans or coding agents.

---

## 3. Problem

Coding agents often lose context across sessions. They need a durable way to know:

- What work exists.
- What is blocked.
- What is ready.
- Which task they should work on next.
- Who is already working on what.
- What was already attempted.
- What counts as done.
- How to safely hand work back to a human or another agent.

Existing solutions have tradeoffs:

- Built-in agent task lists are often session-local.
- Markdown TODO files are readable but lack dependency and claim semantics.
- Full issue trackers are too external and not optimized for agent automation.
- Heavier agent coordination systems can introduce databases, daemons, sync complexity, or hidden state.
- File-based tools like `ticket` are transparent, but leave room for stronger agent-native claiming, verification, and handoff workflows.

---

## 4. Target Users

### Primary user

A software developer or software architect using Pi, Claude Code, Codex, Cursor, or similar coding agents inside a Git repository.

### Secondary users

- Small indie hackers.
- Small product teams.
- Technical founders.
- Engineering managers coordinating agent-assisted work.
- Developers who want local-first task tracking without SaaS overhead.

---

## 5. Core Product Thesis

A good agent task system should be:

1. **Repo-local**  
   The task state lives inside the Git repository.

2. **Human-readable**  
   A developer can inspect and edit tasks with a normal editor.

3. **Machine-readable**  
   Every command supports JSON output.

4. **Dependency-aware**  
   Agents can ask: “What is ready now?”

5. **Claim-safe**  
   Agents can claim work with leases and stale-claim detection.

6. **Verification-oriented**  
   A task should define how completion is proven.

7. **Small and predictable**  
   No daemon, no hidden database, no automatic remote push, no AGENTS.md magic.

---

## 6. Non-Goals

TaskLedger will not initially support:

- Real-time collaboration.
- A web app.
- A hosted backend.
- Complex role hierarchies.
- Multi-repository orchestration.
- Automatic Git pushing.
- Automatic merging.
- Long-running background workers.
- tmux/session management.
- Full Jira/Linear/GitHub Issues replacement.
- AI agent execution itself.

The tool tracks and coordinates work. It does not run the agent in v1.

---

## 7. MVP Scope

The MVP should provide a CLI with the following commands:

```bash
taskledger init
taskledger create
taskledger show
taskledger list
taskledger ready
taskledger dep add
taskledger dep remove
taskledger claim
taskledger release
taskledger stale
taskledger note
taskledger close
taskledger verify
```

Short alias:

```bash
tl
```

Example:

```bash
tl ready --json
tl claim task-abc123 --actor claude-code:frontend
tl verify task-abc123
tl close task-abc123 --evidence verify.log
```

---

## 8. Storage Model

TaskLedger stores all state under:

```text
.taskledger/
  config.yaml
  tasks/
    task-abc123.md
  events.jsonl
```

### 8.1 Task file

Each task is a Markdown file with YAML frontmatter.

Example:

```markdown
---
id: task-abc123
title: Add login form validation
status: open
priority: medium
type: feature
created_at: 2026-05-16T12:00:00Z
updated_at: 2026-05-16T12:00:00Z
created_by: human
assignee: null

depends_on:
  - task-def456

claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null

verify:
  commands:
    - npm test
    - npm run lint
  evidence_required:
    - test_output
    - changed_files
    - summary

tags:
  - frontend
  - auth
---

## Description

Add client-side validation to the login form.

## Acceptance Criteria

- Email field validates email format.
- Password field is required.
- Errors are visible and accessible.
- Existing tests pass.

## Notes

No notes yet.
```

---

## 9. Status Model

Supported task statuses:

```text
open
in_progress
blocked
pending_human
done
cancelled
```

### Meaning

| Status | Meaning |
|---|---|
| `open` | Work exists but is not currently claimed. |
| `in_progress` | Work is claimed by an actor. |
| `blocked` | Work cannot continue because of another task or technical blocker. |
| `pending_human` | Work requires human clarification or decision. |
| `done` | Work is completed and verified. |
| `cancelled` | Work is intentionally abandoned. |

---

## 10. Dependency Model

A task may depend on zero or more other tasks.

A task is **ready** when:

- status is `open`
- all dependencies are `done`
- it has no active claim
- it is not `pending_human`
- it is not `blocked`

Command:

```bash
tl ready
```

Example output:

```text
task-abc123  medium  Add login form validation
task-ghi789  low     Refactor auth error messages
```

JSON output:

```json
[
  {
    "id": "task-abc123",
    "title": "Add login form validation",
    "priority": "medium",
    "status": "open",
    "depends_on": []
  }
]
```

---

## 11. Claiming Model

Claims are first-class.

An actor can claim a ready task:

```bash
tl claim task-abc123 --actor claude-code:frontend
```

Default claim duration:

```text
60 minutes
```

Custom duration:

```bash
tl claim task-abc123 --actor claude-code:frontend --ttl 120m
```

A claim sets:

```yaml
claim:
  actor: claude-code:frontend
  claimed_at: 2026-05-16T12:00:00Z
  expires_at: 2026-05-16T13:00:00Z
  heartbeat_at: 2026-05-16T12:00:00Z
```

A claim is stale when:

```text
now > claim.expires_at
```

Command:

```bash
tl stale
```

A stale claim can be released:

```bash
tl release task-abc123 --force
```

A non-stale claim should not be overwritten unless `--force` is used.

---

## 12. Human Input State

Agents need a clean way to stop and ask for help.

Command:

```bash
tl pending task-abc123 --question "Which auth provider should be supported first?"
```

This sets:

```yaml
status: pending_human
pending:
  question: "Which auth provider should be supported first?"
  requested_by: claude-code:frontend
  requested_at: 2026-05-16T12:20:00Z
```

A human can resolve it:

```bash
tl resolve task-abc123 --answer "Use GitHub OAuth first."
```

This appends the answer as a note and returns the task to `open`.

This command can be added after the MVP if needed, but the schema should allow it from the beginning.

---

## 13. Notes and Handoffs

Agents and humans can append notes:

```bash
tl note task-abc123 --actor claude-code:frontend --message "Implemented validation, but tests fail due to missing mock."
```

Notes should be appended to the Markdown body under `## Notes`.

Each note should include:

```text
timestamp
actor
message
```

Example:

```markdown
### 2026-05-16T12:42:00Z — claude-code:frontend

Implemented validation, but tests fail due to missing mock.
```

Future version may store notes as separate append-only files to reduce merge conflicts, but MVP may keep notes inside the task file.

---

## 14. Verification

Tasks may define verification commands.

Command:

```bash
tl verify task-abc123
```

This runs the commands from the task’s `verify.commands`.

Example:

```yaml
verify:
  commands:
    - go test ./...
    - go vet ./...
```

The command should:

- Run each verification command in order.
- Stop on first failure.
- Return non-zero exit code if verification fails.
- Print human-readable output by default.
- Return structured output with `--json`.

Example JSON:

```json
{
  "task_id": "task-abc123",
  "success": false,
  "commands": [
    {
      "command": "go test ./...",
      "success": false,
      "exit_code": 1
    }
  ]
}
```

---

## 15. Closing Tasks

A task can be closed only when:

- It exists.
- It is not blocked.
- It either has no verification commands or verification has passed.
- The actor owns the claim, unless `--force` is used.

Command:

```bash
tl close task-abc123 --actor claude-code:frontend
```

Optional evidence:

```bash
tl close task-abc123 --actor claude-code:frontend --evidence verify.log
```

Closing sets:

```yaml
status: done
closed_at: 2026-05-16T13:00:00Z
closed_by: claude-code:frontend
```

---

## 16. Event Journal

Every mutating command appends an event to:

```text
.taskledger/events.jsonl
```

Example events:

```json
{"time":"2026-05-16T12:00:00Z","event":"created","task_id":"task-abc123","actor":"human"}
{"time":"2026-05-16T12:10:00Z","event":"claimed","task_id":"task-abc123","actor":"claude-code:frontend"}
{"time":"2026-05-16T12:42:00Z","event":"note_added","task_id":"task-abc123","actor":"claude-code:frontend"}
{"time":"2026-05-16T13:00:00Z","event":"closed","task_id":"task-abc123","actor":"claude-code:frontend"}
```

The task file is the current state.  
The event journal is the audit trail.

MVP does not need full replay support, but events should be valid JSONL from day one.

---

## 17. CLI Design Principles

### 17.1 Human-readable by default

```bash
tl ready
```

Output:

```text
ID            Priority  Title
task-abc123   medium    Add login form validation
```

### 17.2 JSON for automation

```bash
tl ready --json
```

Output:

```json
[
  {
    "id": "task-abc123",
    "priority": "medium",
    "title": "Add login form validation"
  }
]
```

### 17.3 No hidden behavior

The tool must not:

- Start daemons.
- Push to Git remotes.
- Modify AGENTS.md automatically.
- Run agents automatically.
- Open network connections unless explicitly requested.

---

## 18. Configuration

Config file:

```yaml
version: 1
default_claim_ttl: 60m
id_prefix: task
default_verify_commands: []
actors:
  require_actor: true
```

Location:

```text
.taskledger/config.yaml
```

---

## 19. Agent Usage Example

Claude Code can use this flow:

```bash
tl ready --json
```

Pick the highest-priority ready task.

```bash
tl claim task-abc123 --actor claude-code:main --json
```

Read task:

```bash
tl show task-abc123
```

Implement work.

Add note:

```bash
tl note task-abc123 --actor claude-code:main --message "Implemented initial version."
```

Run verification:

```bash
tl verify task-abc123 --json
```

Close:

```bash
tl close task-abc123 --actor claude-code:main
```

If blocked:

```bash
tl pending task-abc123 --actor claude-code:main --question "Should this support GitHub OAuth only or Google too?"
```

---

## 20. AGENTS.md Snippet

The tool should provide a command:

```bash
tl prime
```

It prints recommended instructions but does not modify files automatically.

Example output:

```markdown
## TaskLedger Workflow

This repository uses TaskLedger for agent-visible work tracking.

Before starting work:

1. Run `tl ready --json`.
2. Choose one ready task.
3. Claim it with `tl claim <id> --actor <your-actor-name>`.
4. Read details with `tl show <id>`.
5. Add notes with `tl note <id> --actor <your-actor-name> --message "..."`
6. Run `tl verify <id>` before closing.
7. Close with `tl close <id> --actor <your-actor-name>`.

Do not work on tasks claimed by another actor unless the claim is stale or a human explicitly instructs you.
Do not edit `.taskledger/events.jsonl` manually.
```

---

## 21. Implementation Language

Use **Go**.

Reasons:

- Good fit for a CLI.
- Single static binary.
- Easy cross-platform distribution.
- Fast startup.
- Strong YAML, Markdown, JSON, and filesystem support.
- Easier to maintain than Bash for this level of structure.
- Simpler development than Rust for an MVP.

Suggested libraries:

```text
cobra       CLI framework
viper       config loading, optional
yaml.v3     YAML parsing
goldmark    Markdown handling, optional
ulid        sortable unique IDs
```

Keep dependencies minimal.

---

## 22. Technical Requirements

### 22.1 Repository detection

Commands should find `.taskledger/` by walking upward from the current working directory.

### 22.2 Atomic writes

Task file updates should be atomic:

1. Write temp file.
2. fsync if reasonable.
3. Rename temp file over original.

### 22.3 Locking

Use a simple repo-local lock file for mutating commands:

```text
.taskledger/.lock
```

The lock should prevent simultaneous writes from corrupting files.

MVP can use best-effort file locking.

### 22.4 ID generation

Use stable unique IDs:

```text
task-01HXABC...
```

Prefer ULID or short hash.

### 22.5 Exit codes

Use predictable exit codes:

| Code | Meaning |
|---:|---|
| 0 | Success |
| 1 | Generic error |
| 2 | Invalid arguments |
| 3 | Task not found |
| 4 | Task not ready |
| 5 | Task already claimed |
| 6 | Verification failed |
| 7 | Lock failed |

---

## 23. Success Metrics

MVP is successful if:

- A developer can initialize a repo and create tasks in under one minute.
- Claude Code, Codex, Pi can reliably query ready tasks using JSON.
- Two agents do not accidentally work on the same task if claims are respected.
- Stale claims can be detected and released.
- All task state can be reviewed in Git.
- No hidden service or database is required.
- The project remains understandable after one hour of code reading.

---

## 24. V1 Acceptance Criteria

The first complete version is done when:

- `tl init` creates `.taskledger/config.yaml`, `.taskledger/tasks/`, and `.taskledger/events.jsonl`.
- `tl create` creates a valid Markdown task file.
- `tl list` lists all tasks.
- `tl ready` returns only unblocked, unclaimed open tasks.
- `tl dep add` and `tl dep remove` update dependencies.
- `tl claim` sets claim metadata and changes status to `in_progress`.
- `tl release` clears claim metadata and returns status to `open`.
- `tl stale` lists expired claims.
- `tl note` appends a timestamped note.
- `tl verify` runs configured verification commands.
- `tl close` marks a task as `done`.
- All read commands support `--json`.
- All mutating commands append to `events.jsonl`.
- The README explains the agent workflow clearly.
- The tool never modifies AGENTS.md automatically.
- The tool never pushes to Git automatically.

---

## 25. Suggested Coding Agent Implementation Plan

### Phase 1: Project skeleton

Create a Go CLI project:

User godog as depndency and BDD ad method.

```bash
go mod init github.com/YOURNAME/taskledger
```

Add Cobra commands:

```text
cmd/root.go
cmd/init.go
cmd/create.go
cmd/list.go
cmd/show.go
cmd/ready.go
cmd/claim.go
cmd/release.go
cmd/stale.go
cmd/note.go
cmd/verify.go
cmd/close.go
```

Internal packages:

```text
internal/repo
internal/task
internal/store
internal/events
internal/lock
internal/output
internal/verify
```

---

### Phase 2: Storage and parsing

Implement:

```go
type Task struct {
    ID        string
    Title     string
    Status    string
    Priority  string
    Type      string
    CreatedAt time.Time
    UpdatedAt time.Time
    CreatedBy string
    Assignee  *string
    DependsOn []string
    Claim     Claim
    Verify    Verify
    Tags      []string
    Body      string
}
```

Implement load/save for Markdown with YAML frontmatter.

---

### Phase 3: Core commands

Implement:

```text
init
create
list
show
ready
```

Focus on correctness before polish.

---

### Phase 4: Dependencies and claims

Implement:

```text
dep add
dep remove
claim
release
stale
```

Add tests for:

- blocked dependencies
- ready queue behavior
- active claim
- stale claim
- force release

---

### Phase 5: Notes, verification, close

Implement:

```text
note
verify
close
```

Add tests for:

- note append
- verification success
- verification failure
- close blocked by failed verification
- close with force

---

### Phase 6: Documentation

Create:

```text
README.md
docs/agent-workflow.md
docs/schema.md
docs/comparison.md
```

README headline:

```markdown
# TaskLedger

A transparent Git-native task ledger for humans and AI coding agents.
```

---

## 26. Initial README Positioning

Use this positioning:

```markdown
TaskLedger is a small Git-native task ledger for coding agents.

It stores tasks as Markdown files in your repository, gives agents a dependency-aware ready queue, supports safe claim leases, and records every change in an append-only event log.

No daemon.  
No hidden database.  
No automatic push.  
No AGENTS.md magic.
```

---

## 27. Key Differentiator

The strongest differentiator is not “file-based tasks.”

The strongest differentiator is:

> **Agent-safe task coordination with readable Git-native state.**

That means:

- Claims are explicit.
- Stale work is detectable.
- Dependencies are computable.
- Verification is part of the task.
- Handoffs are recorded.
- Humans can inspect everything.

This is the main reason to build your own version instead of only using `ticket`.
