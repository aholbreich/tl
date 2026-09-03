# tl task ledger

This directory contains the repository's work ledger, maintained with
[tl](https://github.com/aholbreich/tl).

- `tasks/*.md` — human-readable tasks
- `events.jsonl` — append-only audit history
- `config.yaml` — ledger configuration

You can inspect these files without installing `tl`. When changing the ledger,
use the `tl` CLI so its append-only history remains complete.
