---
id: task-a46
title: 'Remove bare Created task task-r75 create shorthand — keep only  and '
status: cancelled
priority: medium
type: task
created_at: 2026-05-31T09:51:48Z
updated_at: 2026-05-31T09:54:52Z
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

## Description

## Problem

Created task task-8ym (bare root with positional args) is silently treated as Created task task-9s0. This is implemented in `cmd/root.go`'s `RunE` and creates multiple problems:

1. **Three ways to do the same thing**: `tl "title"`, `tl create "title"`, `tl add "title"`
2. **Hidden behavior**: Not discoverable via `--help`, not shown in README command reference
3. **Root flag pollution**: Root registers `-d`, `-t`, `--priority`, `--tag`, `--ref`, `--json` — duplicated from `create` — making `tl --help` confusing
4. **No analogue elsewhere**: `tl note`, `tl block`, `tl pending` don't have bare-root shortcuts
5. **Duplicate flag registrations**: same flags declared in both `root.go` and `create.go`

## What to do

### cmd/root.go
- Remove the `RunE` function entirely (the bare-args-to-create routing) — keep only `cmd.Help()` behavior
- Remove the create-specific root-level flag vars (`title`, `description`, `priority`, `tags`, `refs`, `asJSON`)
- Remove the corresponding `root.Flags().*Var` registrations for those flags
- Keep only `--color` as the sole root persistent flag
- Clean up any unused imports (check `os`, `fmt`, etc.)

### features/create.feature
- Remove the scenario: `Creating a task records without "create" or "add"`
- This is the scenario that tests `tl "Shortcut syntax"`

### Validation
- Run `make bdd` to confirm the godog suite still passes (minus the removed scenario)
- Run `make test` for full suite
- Run `go vet ./...` and `gofmt -d .` to catch any issues
- Check that `tl --help` no longer shows create-specific flags (`-d`, `-t`, `--priority`, `--tag`, `--ref`, `--json`)

## Non-goals
- Do NOT remove `tl add` — the cobra Alias on create is fine and documented
- Do NOT change `tl create` behavior in any way
- Do NOT touch any other command

## Notes

- 2026-05-31T09:54:52Z [pi:analysis] cancelled: title garbled from shell escaping, recreating
