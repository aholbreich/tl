## tl workflow compact
Use tl agents and tl <cmd> --help to learn the workflow and available commands.
tl ready [--tag x] --json - find unclaimed ready work.
tl show <id> && tl history <id> - read scope, references, notes, and prior events.
Inspect References before editing files.
Claim tasks:
tl claim <id> --actor agent-name - take/renew a lease before editing or implementing.
Keep progress in tasks:
tl note <id> -m "..." --actor agent-name - record progress, decisions, and handoffs.
tl dep add <id> --on <id> / tl dep remove <id> --on <id> - manage dependencies.
tl block/unblock, tl pending/resolve — handle blockers and human decisions.
tl close/cancel/release <id> --actor agent-name - end with explicit ledger state.

Use --json on read commands for automation.
Do not edit .tl/events.jsonl manually.
