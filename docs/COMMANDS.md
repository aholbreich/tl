# Command reference

User-visible flags for each implemented `tl` command. For canonical per-command behavior, read the matching file in [`../features/`](../features); the feature file wins if this page drifts. Every command also accepts
`--help` at the terminal.

---

## `tl init`

Initialize a `.taskledger/` ledger in the current directory.

```
(no flags)
```

## `tl create [title] [options]`

Create a new task. The title is a required positional argument or passed via
`--title`.

```
    --title              Task title (required, or positional argument)
-d, --description        Task description (stored under ## Description)
-p, --priority           Task priority (l|low, m|medium, h|high) [default: medium]
-t, --type               Task type
    --tag                Tag to apply (repeatable)
    --actor              Creator actor [default: human]
    --json               Emit JSON output
```

## `tl list`

List active tasks in the ledger, sorted by priority then identifier. Closed
statuses (`done`, `cancelled`) are hidden by default. Human output includes
`ID`, `Status`, `Priority`, `Claimed By`, and `Title`.

```
-a, --all                Include closed tasks (done and cancelled)
    --claimed-by         Only show tasks claimed by this actor
    --json               Emit JSON output
```

## `tl show TASK_ID`

Show a task in detail. Human output includes the identifier, title, status,
priority, dependencies, claim state, and Markdown body content such as notes.
Accepts both full IDs (`task-k5g`) and bare short codes (`k5g`).

```
    --json               Emit JSON output
```

## `tl ready`

List tasks that are ready to be claimed. A task is ready when it is `open`
(or `in_progress` with an expired claim), all dependencies are `done`, and
no active claim exists.

```
    --json               Emit JSON output
```

## `tl claim TASK_ID`

Claim a task with a time-limited lease. Sets status to `in_progress` and
records claim expiry. Rejects if another actor holds an active claim (exit
5) or dependencies are unmet (exit 4). Uses the same actor resolution chain
as `note`.

```
    --actor              Claiming actor (optional; auto-resolved if unset)
    --ttl                Lease duration, e.g. 60m or 2h (default from config)
    --json               Emit JSON output
```

## `tl dep add TASK_ID --on TASK_ID`

Add a dependency link between tasks. Both tasks must exist (exit 3 if not).
Idempotent — adding the same dependency twice is a no-op.

```
    --on                 Target task to depend on (required)
```

## `tl dep remove TASK_ID --on TASK_ID`

Remove a dependency link. Both tasks must exist (exit 3 if not).
Idempotent — removing a non-existent dependency is a no-op.

```
    --on                 Target task to drop as a dependency (required)
```

## `tl note TASK_ID`

Append a timestamped note to a task's body under a `## Notes` section. Notes
are the human-facing audit trail — use them for progress updates, handoff
context, and decision records.

```
-m, --message            Note message (required)
    --actor              Actor writing the note (optional; auto-resolved)
```

## `tl close TASK_ID`

Mark a task as `done`. Unclaimed open tasks may be closed by any actor.
Claimed tasks may be closed by the claiming actor, or by another actor with
`--force`. Rejects blocked and already-done tasks.

```
    --actor              Actor closing the task (optional; auto-resolved)
    --force              Close even when another actor holds an active claim
    --json               Emit JSON output
```

## `tl release TASK_ID`

Voluntarily release a claim on a task, returning it to `open`. Only the
claiming actor may release unless `--force` is used.

```
    --actor              Actor releasing the claim (optional; auto-resolved)
    --force              Release even when another actor holds the claim
    --json               Emit JSON output
```

## `tl agents`

Print a recommended `AGENTS.md` snippet for TaskLedger-aware agents. Writes
only to stdout and never edits `AGENTS.md` for you. Commands in the snippet
are formatted as Markdown code spans, for example `tl ready --json`.

```
(no flags)
```
---

## Setup errors

Commands that need ledger state are non-interactive for agent safety. If
`.taskledger/` is missing, they exit with code `1` and print:

```text
TaskLedger is not initialized in this repository.
Run `tl init` from the repository root to create .taskledger/.
```
