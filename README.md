# TaskLedger

> A Git-native task ledger for humans and AI coding agents.

TaskLedger (`tl`) stores tasks as Markdown files with YAML frontmatter inside
your repository, gives agents a dependency-aware ready queue, supports safe
claim leases with automatic actor resolution, and records every change in an
append-only event journal.

No daemon. No hidden database. No automatic push. No AGENTS.md magic.

---

## Install

From source (Go 1.25+):

```sh
git clone https://github.com/aholbreich/taskledger
cd taskledger
make install                # installs `tl` to $HOME/bin
```

Cross-platform release archives:

```sh
make dists                  # produces tl-linux-amd64.tar.gz, tl-darwin-arm64.tar.gz, …
```

---

## Quickstart

```sh
tl init                                                          # one-time per repo
tl create "Add login form validation"
tl create "Refactor auth errors" -t chore -p low --tag auth
tl list
tl show <id>                                                     # full id or bare short code
```

Agent workflow:

```sh
tl ready --json                                                  # what's available?
tl claim <id>                                                    # take a lease (actor auto-detected)
tl show <id>                                                     # read the details
tl note <id> -m "Initial implementation done."                   # record a handoff note
tl close <id>                                                    # mark as done
```

Actor identity is resolved automatically: `--actor` flag > `TL_ACTOR` env >
`ACTOR_NAME` env > `BEADS_ACTOR` env > agent auto-detection. Use `--actor`
to override.

---

## Commands

Full behavioral spec lives under [`features/`](features). Flags currently
exposed by the implemented commands:

### `tl init`

Initialize a `.taskledger/` ledger in the current directory.

```
(no flags)
```

### `tl create [title] [options]`

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

### `tl list`

List active tasks in the ledger, sorted by priority then identifier. Closed
statuses (`done`, `cancelled`) are hidden by default. Human output includes
`ID`, `Status`, `Priority`, `Claimed By`, and `Title`.

```
-a, --all                Include closed tasks (done and cancelled)
    --claimed-by         Only show tasks claimed by this actor
    --json               Emit JSON output
```

### `tl show TASK_ID`

Show a task in detail. Human output includes the identifier, title, status,
priority, dependencies, claim state, and Markdown body content such as notes.
Accepts both full IDs (`task-k5g`) and bare short codes (`k5g`).

```
    --json               Emit JSON output
```

### `tl ready`

List tasks that are ready to be claimed. A task is ready when it is `open` (or
`in_progress` with an expired claim), all dependencies are `done`, and no
active claim exists.

```
    --json               Emit JSON output
```

### `tl claim TASK_ID`

Claim a task with a time-limited lease. Sets status to `in_progress` and
records claim expiry. Rejects if another actor holds an active claim (exit 5)
or dependencies are unmet (exit 4). Uses the same actor resolution chain as
`note`.

```
    --actor              Claiming actor (optional; auto-resolved if unset)
    --ttl                Lease duration, e.g. 60m or 2h (default from config)
    --json               Emit JSON output
```

### `tl dep add TASK_ID --on TASK_ID`

Add a dependency link between tasks. Both tasks must exist (exit 3 if not).
Idempotent — adding the same dependency twice is a no-op.

```
    --on                 Target task to depend on (required)
```

### `tl note TASK_ID`

Append a timestamped note to a task's body under a `## Notes` section. Notes
are the human-facing audit trail — use them for progress updates, handoff
context, and decision records.

```
-m, --message            Note message (required)
    --actor              Actor writing the note (optional; auto-resolved)
```

### `tl close TASK_ID`

Mark a task as `done`. Unclaimed open tasks may be closed by any actor. Claimed
tasks may be closed by the claiming actor, or by another actor with `--force`.
Rejects blocked and already-done tasks.

```
    --actor              Actor closing the task (optional; auto-resolved)
    --force              Close even when another actor holds an active claim
    --json               Emit JSON output
```

### `tl agents`

Print a recommended `AGENTS.md` snippet for TaskLedger-aware agents. The command
writes only to stdout and never edits `AGENTS.md` for you. Commands in the
snippet are formatted as Markdown code spans, for example `tl ready --json`.

```
(no flags)
```

### Not yet implemented

`release`, `stale`, `dep remove`, `pending`, `resolve` — specified in
[`features/`](features), implementation in progress. See [`docs/PRD.md`](docs/PRD.md)
§6 for the command index.

### Setup errors

Commands that need ledger state are non-interactive for agent safety. If
`.taskledger/` is missing, they exit with code `1` and print:

```text
TaskLedger is not initialized in this repository.
Run `tl init` from the repository root to create .taskledger/.
```

---

## Implementation Status

The BDD suite runs features tagged `@implemented`.

| Area | Status |
|---|---|
| `tl init` | ✅ Implemented |
| `tl create` | ✅ Implemented |
| `tl list` | ✅ Implemented |
| `tl show` | ✅ Implemented (bare short codes supported) |
| `tl ready` | ✅ Implemented (stale-claim aware) |
| `tl dep add` | ✅ Implemented |
| `tl claim` | ✅ Implemented (auto actor resolution) |
| `tl note` | ✅ Implemented |
| `tl close` | ✅ Implemented |
| `tl agents` | ✅ Implemented |
| Actor identity resolution | ✅ Implemented (`--actor` > `TL_ACTOR` > `ACTOR_NAME` > `BEADS_ACTOR` > auto-detect) |
| Friendly missing-ledger hint | ✅ Implemented |
| `tl dep remove` | Specified, pending |
| `tl release` / `tl stale` | Specified, pending |
| `tl pending` / `tl resolve` | Specified, pending |

---

## Storage

```
.taskledger/
  config.yaml      # defaults
  tasks/
    task-<3>.md    # one file per task (Markdown + YAML frontmatter)
  events.jsonl     # append-only audit trail
```

A created task looks like:

```markdown
---
id: task-x3n
title: Add login validation
status: open
priority: medium
type: ""
created_at: 2026-05-17T00:45:40Z
updated_at: 2026-05-17T00:45:40Z
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

Validate email format and require a password.
```

---

## Exit codes

`0` success · `1` generic · `2` invalid args · `3` task not found ·
`4` task not ready · `5` already claimed · `7` lock failed

---

## Development

```sh
make build                  # version-stamped local binary
make test                   # all Go tests
make bdd                    # godog suite only
make dists                  # cross-platform release archives
make clean
```

CI runs `gofmt`, `go vet`, `make build`, `make test` on every PR and push to
`main` (see [`.github/workflows/ci.yaml`](.github/workflows/ci.yaml)).
Tag-triggered releases build all platforms and publish a GitHub Release.

The BDD suite lives in [`bdd/`](bdd/) and runs the features tagged
`@implemented`; the rest are pending-implementation specs.

---

## Further reading

- [`docs/PRD.md`](docs/PRD.md) — design intent, non-goals, status enum, exit codes
- [`features/`](features/) — Gherkin behavioral spec, one file per command
- [`AGENTS.md`](AGENTS.md) — leading doc for any agent working in this repo
- [`docs/gherkin-guidelines.md`](docs/gherkin-guidelines.md) — Gherkin style rules
