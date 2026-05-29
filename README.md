# tl tool - Task ledger for your repository

> A Git-native task ledger for humans and AI coding agents.

## Why tl?

Chat history disappears. TODO files are not dependency-aware. GitHub Issues are remote-first.
`tl` gives every repository a small local task ledger that both humans and agents can read and update. It anchors the work to be done around one repository or project.

- Git-native: state lives in `.tl/` folder
- Human-readable: tasks are Markdown with YAML frontmatter
- Agent-readable: read commands support `--json`
- Coordination-safe: claims use leases and stale work is detectable
- Handoff-friendly: notes preserve context across sessions
- Flexible: tasks are the unit of work — `tl` adapts to your (agentic) flow.
- Boring by design: no daemon, no database, no automatic push

**Contents:** [How it compares](#how-tl-compares) · [Installation Options](#installation-options) · [Quickstart](#quickstart) · [Commands](#commands) · [Implementation status](#implementation-status) · [Development](#development) · [Further reading](#further-reading)

---

## How tl compares

`tl` sits in a small but real category: Git-native task trackers built for
humans **and** AI coding agents. Its nearest neighbours are
[Beads](https://github.com/steveyegge/beads) and
[Backlog.md](https://github.com/MrLesk/Backlog.md). The honest shape of the
trade-offs (these tools move fast — verify before quoting):

|                                   | **tl**                                | **Beads**                    |   **Backlog.md**          | **GitHub Issues**   |
| --------------------------------- | ------------------------------------- | ---------------------------  | ----------------------- | ------------------- |
| State lives in                    | Markdown + append-only `events.jsonl` | Embedded Dolt DB (`.beads/`) | Markdown + YAML         | A remote service    |
| Read / edit a task in your editor | ✅                                    | ❌ (binary DB)               | ✅                      | Via the web UI      |
| Inspect history with `git diff`   | ✅                                    | ～ (DB diffs)                 | ✅                      | ❌                  |
| Dependency-aware `ready`          | ✅                                    | ✅                           | ～                      | ～                  |
| Coordination primitive            | Leases + stale-work detection         | `ready` / gates / routing    | Status + Kanban board   | Assignees           |
| Sync model                        | Plain `git` — you stay in control     | `bd dolt push` / `pull`      | `git`                   | Always online       |
| Extra surface area                | None — one static binary              | Embedded database            | React web UI + MCP server | SaaS + API        |
| Works offline / at a commit       | ✅                                    | ✅                           | ✅                      | ❌                  |

**Why `tl` and not Beads?** Beads keeps its ledger in an embedded Dolt
database under `.beads/` and syncs it with `bd dolt push/pull`. That buys a
genuine dependency graph, "memory decay", and cross-machine sync — at the price
of a binary store you commit to git but cannot read or `diff` as text. `tl`
makes the opposite bet: every task is a Markdown file you can read, `grep`,
edit, and `git diff` with zero tooling, and the only moving parts are those
files plus an append-only event log. Want a database that remembers things for
your agent? Use Beads. Want a ledger you can read at any commit and reason about
yourself? Use `tl`.

**Why `tl` and not Backlog.md?** Backlog.md shares the plain-Markdown,
Git-native philosophy but centers on a Kanban board and ships a web UI and an
MCP server. `tl` is deliberately leaner: no board, no server, no daemon. Its
core primitive is *coordination* — explicit claims with leases, detectable
stale work, and a recorded handoff trail — not visualization.

The shared backbone, in `tl`'s words:

> **Agent-safe task coordination with readable, Git-native state.**
> Claims are explicit · stale work is detectable · dependencies are computable
> · handoffs are recorded · humans can inspect everything.

**What `tl` deliberately leaves out** (see the
[non-goals](docs/PRD.md)): no web UI or Kanban board, no embedded database or
automatic sync, no hosted backend, and it does not run the agent itself. `tl`
tracks and coordinates the work; Git and your agent do the rest.

---



## Installation Options


### Homebrew (macOS / Linux)

```sh
brew install aholbreich/tap/tl           # latest stable release
brew install --HEAD aholbreich/tap/tl    # or: build from current main
```

If you install multiple tools from the same tap, you can tap once:

```sh
brew tap aholbreich/tap
brew install tl
```

Prebuilt binaries are available for **macOS (Intel + Apple Silicon)** and **Linux (amd64 + arm64)**.

### RPM (Fedora / Red Hat)

Add the Holbreich RPM repository:

```sh
# Documentation: https://aholbreich.github.io/rpm-repo/#installation-fedora-centos-redhat
echo '[Holbreich]
name=Holbreich Repository
baseurl=https://aholbreich.github.io/rpm-repo/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/holbreich.repo
```

Install `tl`:

```sh
sudo dnf install tl
tl --version
```

If you run into issues with the RPM repository, see the
[rpm-repo project](https://github.com/aholbreich/rpm-repo).

### Install script (macOS / Linux)

```sh
curl -fsSL https://raw.githubusercontent.com/aholbreich/tl/main/install.sh | sh
```

Install a specific version or target directory:

```sh
curl -fsSL https://raw.githubusercontent.com/aholbreich/tl/main/install.sh | sh -s -- --version 0.4.4
curl -fsSL https://raw.githubusercontent.com/aholbreich/tl/main/install.sh | sh -s -- --bin-dir "$HOME/.local/bin"
```


### From source

```sh
git clone https://github.com/aholbreich/tl
cd tl
make install                # installs `tl` to $HOME/bin
```

Cross-platform release archives:

```sh
make dists                  # tl-linux-amd64.tar.gz, tl-darwin-arm64.tar.gz, …
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

Actor identity resolves in order: `--actor` flag > `TL_ACTOR` env >
`ACTOR_NAME` env > `BEADS_ACTOR` env > agent auto-detection.

### Shell completion

`tl` ships completions for bash, zsh, fish, and PowerShell. Pressing TAB on
a `TASK_ID` argument suggests the actual task IDs from the current ledger.

```sh
tl completion --install            # auto-detect shell from $SHELL
tl completion --install bash       # or pick one explicitly
```

The script is written to the canonical XDG path for the chosen shell:
`~/.local/share/bash-completion/completions/tl` (bash),
`~/.config/fish/completions/tl.fish` (fish), `~/.zsh/completions/_tl` (zsh —
plus an fpath line to add to `~/.zshrc`). Open a new shell to activate.

For a one-off in the current session: `source <(tl completion bash)`.

---

## Commands

The whole surface at a glance:

```sh
# Set up
tl init                            # create the .tl/ ledger (once per repo)
tl completion --install            # enable TAB completion for task IDs

# Define work
tl create "<title>" [-t type -p prio --tag x -d "..."]  # add a task
tl refine <id> [-p prio -t title --edit]                # edit an existing task
tl dep add <id> --on <id>                               # declare a dependency
tl dep remove <id> --on <id>                            # drop one

# Do the work
tl ready [--tag x] [--json]        # unclaimed, unblocked tasks
tl claim <id>                      # take a time-limited lease (re-run = heartbeat)
tl note <id> -m "..."              # record progress / handoff context
tl close <id>                      # done and verified

# When it doesn't just finish
tl block <id> -m "..."             # external blocker; releases the claim
tl unblock <id>                    # blocker cleared; back to open
tl pending <id> --question "..."   # need a human decision; releases the claim
tl resolve <id> --answer "..."     # human answers; task reopens
tl cancel <id> -m "..."            # won't be done
tl release <id>                    # step away cleanly (leave a note first)

# Inspect
tl list [--all --status s --tag t --mine] [--json]      # browse tasks
tl show <id> [--json]              # full task detail
tl history [<id>] [--json]         # event-by-event audit trail
tl stale                           # claims whose lease has expired
tl agents                          # print a paste-ready agent workflow guide
```

- Walkthrough: [`docs/usage.md`](docs/usage.md) — tl by example, flow by flow
- Behavioral spec: [`features/`](features) (one `.feature` file per command)
- Per-command flags: `tl <cmd> --help`

---

## Implementation status

Implemented commands carry the `@implemented` tag in their feature file.
`make bdd` runs only the implemented suite; untagged features are the
binding contract for unimplemented commands.

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

---

## Further reading

- [`docs/usage.md`](docs/usage.md) — tl by example, flow by flow
- [`docs/tech-docs.md`](docs/tech-docs.md) - some implementation detail
- [`docs/PRD.md`](docs/PRD.md) — design intent, non-goals, status enum


