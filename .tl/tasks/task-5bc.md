---
id: task-5bc
title: Remove bare-root tl "title" create shortcut
status: done
priority: medium
type: task
created_at: 2026-05-31T09:54:57Z
updated_at: 2026-05-31T10:04:36Z
created_by: human
assignee: null
depends_on: []
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags:
  - polish
  - cli
---

## Problem

`tl "Some task title"` (bare root command with positional arg) is silently treated as `tl create "Some task title"`. This is wired in `cmd/root.go`'s `RunE`.

Issues:

1. **Three ways to create**: `tl "title"`, `tl create "title"`, `tl add "title"` — inconsistent
2. **Hidden behavior**: Not discoverable via `--help`, not shown in README command reference
3. **Root flag pollution**: `-d`, `-t`, `--priority`, `--tag`, `--ref`, `--json` are registered on root (duplicating `create`'s flags), making `tl --help` confusing
4. **No analogue**: No other command (`note`, `block`, `pending`) has a bare-root shortcut
5. **Duplicate registrations**: Same flags declared in both `root.go` and `create.go`

## What to do

### cmd/root.go

- Remove the `RunE` function body (the bare-args-to-create routing) — let it fall through to `cmd.Help()`
- Remove root-level flag vars: `title`, `description`, `priority`, `tags`, `refs`, `asJSON`
- Remove `root.Flags().*Var` registrations for those flags
- Keep only `--color` as the sole root persistent flag
- Clean up unused imports after removal

### features/create.feature

- Remove the scenario: `Creating a task records without "create" or "add"`
  - This is the only scenario that tests `tl "Shortcut syntax"`

### Validation

- `make bdd` — godog suite passes (minus removed scenario)
- `make test` — full suite passes
- `go vet ./...` and `gofmt -d .` — clean
- `tl --help` — no longer shows create-specific flags

## Non-goals

- Do **not** remove `tl add` (the cobra Alias on `create` is fine and documented)
- Do **not** change `tl create` behavior
- Do **not** touch any other command

## Notes

- 2026-05-31T09:55:26Z [pi:analysis] note: Full implementation context written into task body. Key files: cmd/root.go (remove RunE + root flags), features/create.feature (remove shortcut scenario).
- 2026-05-31T10:04:36Z [pi] note: Implemented removal of bare-root create shortcut: root no longer registers create-specific flags, bare root titles now return exit code 2 with guidance to use tl create/tl add, create BDD now asserts no task is created, references feature uses tl add instead of bare shorthand, docs no longer say bare title is enough. Validation: gofmt, go test ./cmd -v, make bdd, make test, go vet ./..., gofmt -d ., git diff --check, and tl --help root flags checked.
