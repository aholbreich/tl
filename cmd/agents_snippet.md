## TaskLedger Workflow

This repository uses TaskLedger (`tl`) for local task coordination between humans and agents.

Set `TL_ACTOR` once at the start of your session so you don't need `--actor` on each command:

```sh
export TL_ACTOR=claude-code:<purpose>
```

When starting work:

1. Pick a task:
   - `tl ready --json` for unclaimed work, or `tl ready --tag <role> --json` to filter by role-ish tags.
   - `tl show <task-id>` when handed a specific task.
   - `tl history <task-id>` if the task was previously worked on; read prior notes before starting.
2. Claim it before editing files:
   `tl claim <task-id>`
3. Inspect the task details:
   `tl show <task-id>`
4. Do the work. Re-run `tl claim <task-id>` periodically on long work — it extends the lease (heartbeat pattern).
   Record important context, decisions, blockers, or handoff notes:
   `tl note <task-id> -m "..."`
5. Pick the correct exit:
   - `tl close <task-id>` — work is done and verified.
   - `tl cancel <task-id> -m "<reason>"` — work won't be done.
   - `tl block <task-id> -m "<blocker>"` — external blocker; claim is released.
   - `tl pending <task-id> --question "..."` — you need a human decision; claim is released.
   - `tl release <task-id>` — you're stepping away cleanly; leave a comprehensive note first.

   (`cancel`, `block`, `pending` are spec'd in `features/`; check the current `@implemented` set with `make bdd` before relying on them.)

Rules:

- Do **not** work on a task claimed by another active actor unless explicitly told.
- If your work uncovers a separable piece of work, create a follow-up task with `tl create` rather than silently expanding scope.
- Prefer tasks from `tl ready`; blocked, pending, done, cancelled, or actively claimed tasks are not ready.
- Leave notes for partial progress, failed approaches, decisions, and handoffs.
- Do **not** edit `.taskledger/events.jsonl` manually.
- Ask before editing `AGENTS.md` or other project instruction files.
- If `.taskledger/` is missing, ask the human whether to run `tl init`.
