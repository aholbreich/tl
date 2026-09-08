# tl cli - Task ledger for your repository

> A Git-native task ledger for humans and AI coding agents.

[![CI](https://github.com/aholbreich/tl/actions/workflows/ci.yaml/badge.svg)](https://github.com/aholbreich/tl/actions/workflows/ci.yaml)
[![Release](https://img.shields.io/github/v/release/aholbreich/tl)](https://github.com/aholbreich/tl/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/aholbreich/tl)](https://goreportcard.com/report/github.com/aholbreich/tl)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

> **Quick start** — the shortest path to a shared human-and-agent task ledger:
>
> ```sh
> brew install aholbreich/tap/tl && tl init && tl agents --write-files
> ```
>
> Then `tl ready` to see what's open, `tl claim <id>` to take a task, and
> `tl close <id>` when it's done. Install options for every platform are
> [below](#installation-options).

<img src=".github/tl-demo.svg" alt="tl demo - init, create, ready, claim, note, close" width="100%">

## Why tl?

AI coding agents are a regular part of software teams now, and the hardest
problem in a repository is no longer writing code — it is keeping humans and
agents coordinated. Chat threads disappear. `TODO.md` files drift. GitHub
Issues live on a remote server and don't follow the code.

`tl` gives every repository a small, local task ledger that both humans and
agents read and update — no daemon, no database, no remote service. Its
differentiator from other Git-native trackers is **agent-safe coordination**:
explicit claims with time-limited leases, dependency-aware `ready` lists,
detectable stale work, and a recorded handoff trail — all in state you can
read, `diff`, and reason about with any tool.

- **Agent-safe coordination:** claims are explicit and lease-based, stale work is detectable, handoffs are recorded — agents don't silently step on each other
- **Dependency-aware:** `tl ready` only lists work whose blockers are done
- **Git-native:** state lives in `.tl/` — commit it, diff it, branch it, review it in any PR
- **Human-readable:** tasks are plain Markdown with YAML frontmatter — read or edit any task in your editor
- **Agent-readable:** every read command supports `--json`, every write can be attributed with `--actor`
- **Boring by design:** no daemon, no database, no git hooks, no automatic push — you decide when to sync

For an honest feature-by-feature comparison with the nearest tools —
[Beads](https://github.com/steveyegge/beads),
[Backlog.md](https://github.com/MrLesk/Backlog.md) — and with GitHub Issues,
see [How tl cli compares](#how-tl-cli-compares).

**Contents:** [Quickstart](#quickstart) · [Setup for agent collaboration](#setup-for-agent-collaboration) · [Installation Options](#installation-options) · [Commands](#commands) · [How tl cli compares](#how-tl-cli-compares) · [Development](#development) · [Further reading](#further-reading)

---

## Quickstart

One `tl init` per repository creates the ledger and nothing else:

```sh
tl init                                                          # create .tl/ (once per repo)
tl completion --install                                          # TAB-complete task IDs (one-time)
tl create "Add login form validation"                            # add a task
tl create "Refactor auth errors" -t chore -p low --tag auth      # with type, priority, tag
tl list                                                          # see everything
tl show <id>                                                     # full task detail
```

Take a task from `ready`, work it, and close it:

```sh
tl ready                              # unclaimed, unblocked tasks
tl claim <id>                         # take a time-limited lease (re-run = heartbeat)
tl note <id> -m "Initial pass done."  # record progress for the next person
tl close <id>                         # done and verified
```

### Setup for agent collaboration

Once the ledger exists, hand your agents the playbook in one step:

```sh
tl agents --write-files                # merge the tl workflow into AGENTS.md, CLAUDE.md, …
```

This injects a managed workflow block into the agent instruction files already
present in your repo (`AGENTS.md`, `CLAUDE.md`, `.cursorrules`, and friends),
so every agent that reads them also knows how to use `tl`. For constrained
context windows:

```sh
tl agents --compact                    # print the short version
tl agents --write-files --compact      # write the short version
```

From then on, the agent loop is:

```sh
tl ready --json                          # what's claimable right now?
tl claim <id> --actor agent-a            # take a lease (and say who you are)
tl show <id>                             # read the task in full
tl note <id> -m "Blocked on the API key; handing back." --actor agent-a
tl release <id> --actor agent-a          # step away cleanly — or tl close when done
```

Identity resolves in order: `--actor` flag > `TL_ACTOR` env > `ACTOR_NAME` env
> agent auto-detection (Claude Code, Codex, aider, Windsurf, pi, …). Setting
`TL_ACTOR` once per session is the easiest way to stay attributed.

---

## Installation Options

Latest releases are published to the
[GitHub Releases page](https://github.com/aholbreich/tl/releases/latest) as
prebuilt archives for **Linux** and **macOS** (amd64 + arm64) and **Windows**
(amd64 + arm64). Every release triggers an automatic update of the Homebrew
tap and the RPM repository.

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

### Arch Linux / Omarchy (AUR — coming soon)

> **Coming soon:** the `tl-bin` package is prepared, but publication is waiting
> for AUR account registration. The commands below will work after the package
> is published.

On Omarchy, use its AUR package helper:

```sh
omarchy pkg aur add tl-bin
```

On other Arch-based systems, use an AUR helper such as `yay`:

```sh
yay -S tl-bin
```

Or build and install directly from the AUR:

```sh
git clone https://aur.archlinux.org/tl-bin.git
cd tl-bin
makepkg -si
```

The AUR package is named `tl-bin` because it packages the prebuilt GitHub
release binary. It installs `/usr/bin/tl`, so the command stays `tl`.

### Install script (macOS / Linux)

```sh
curl -fsSL https://raw.githubusercontent.com/aholbreich/tl/main/install.sh | sh
```

Install a specific version or target directory:

```sh
curl -fsSL https://raw.githubusercontent.com/aholbreich/tl/main/install.sh | sh -s -- --version 0.9.0
curl -fsSL https://raw.githubusercontent.com/aholbreich/tl/main/install.sh | sh -s -- --bin-dir "$HOME/.local/bin"
```

### Windows

Download the latest `tl-windows-<arch>.zip` from the
[Releases page](https://github.com/aholbreich/tl/releases/latest) and unpack
`tl.exe` into a directory on your `PATH`.

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

<details>
<summary>RPM (Fedora / Red Hat) — repo-based install</summary>

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
</details>

---

## Commands

The whole surface at a glance:

```sh
# Set up
tl init                            # create the .tl/ ledger (once per repo)
tl completion --install            # enable TAB completion for task IDs

# Define work
tl create "<title>" [-t type -p prio --tag x --ref r -d "..."]  # add a task
tl refine <id> [-p prio -t title --edit]                # edit an existing task
tl refine <id> [--add-ref r --remove-ref r]             # attach/detach references
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
tl remove <id> -m "..." [--force]  # delete a mistaken task file from the active ledger
tl release <id>                    # step away cleanly (leave a note first)

# Inspect
tl list [--all --status s --tag t --mine] [--type t --priority p] [--json]
tl list --dashboard [--tag t] > tasks.md  # regeneratable Markdown overview
tl list --spec-status              # add a column for referenced .feature specs
tl show <id> [--json]              # full task detail
tl history [<id>] [--json]         # event-by-event audit trail
tl stale                           # claims whose lease has expired
tl doctor [--json] [--fix] [--force] # scan ledger for integrity issues (optionally repair)

# Agents
tl agents [--compact] [--write-files [--dry-run] [--file path]] # print or install agent workflow guide
```

### Markdown dashboard

`tl list --dashboard` writes a plain Markdown snapshot to stdout, grouped by
status in list order (pending human, blocked, in progress, open; then done and
cancelled with `--all`). Each task includes its ID, title, status, priority,
type, claimant, a one-line description capped at 240 characters, and references.
Missing types display and filter as `task`. References are displayed as text,
not interpreted as links. Notes are omitted to keep the overview compact.

Metadata occupies one line in status · priority · type · claimant order:

`in_progress` · **high** · `task` · 👤 `aho`

The 👤 icon identifies the claimant; `-` means unclaimed.

All list filters compose: `--status`, `--claimed-by`, `--mine`, `--tag`,
`--type`/`-t`, and `--priority`/`-p` (including `l`, `m`, `h` aliases).
`--status done` or `--status cancelled` includes that status without `--all`.
`--json` takes precedence over `--dashboard` and keeps the existing list JSON
format. Markdown never includes terminal colors or generation timestamps, and
listing does not modify the ledger.

```sh
tl list --dashboard --type feature --priority high > roadmap.md
tl list --dashboard --tag docs > docs-tasks.md
```

There is no dedicated area or due-date field; use tags to scope areas of work.
Watch mode and dependency-tree rendering are not part of the dashboard.

### Spec status

A reference whose path ends in `.feature` is treated as a **spec reference**.
That is the whole rule — there is no new flag on `tl create`, no frontmatter
field and no configuration. `tl show` marks such a reference `(spec)`, which
costs nothing because it is a string test.

`tl list --spec-status` (and `tl ready --spec-status`) adds a column showing
each referenced spec, whether the file is there, and how many scenarios it
holds. A `Scenario Outline` counts once, however many `Examples` rows it has.

```
ID        Status  Title                     Spec
task-7fi  open    Add tl tree               features/tree.feature (missing)
task-cys  open    Add tl agents --remove    features/agents.feature (14)
task-wke  open    Add --type field          -
```

This is the **only** path that opens a file outside `.tl/`. Plain `tl list`,
`tl ready` and `tl show` read the ledger and nothing else, so the default
output never depends on the state of your working tree. A missing or
unreadable spec renders as `(missing)` or `(unknown)` rather than failing the
listing — `tl doctor` is what complains about dead references.

In `--json`, the `spec` key is always present: `null` when the flag was not
passed, an array when it was, so a consumer's schema never changes based on
which flags were used or on whether the project writes Gherkin.

```sh
tl list --spec-status --json | jq -r '.[] | select(.spec[]?.state == "missing") | .id'
```

**What this does not tell you.** The link is to a *file*. A feature file
usually describes a capability while a task is a slice of one, so a spec
column says "this task points at a spec that exists", not "this task's
behaviour is specified" — and never that the work is done. Delivery state
stays in the task's own status, beside it. See
[`.decisions/0002-reading-referenced-files.md`](.decisions/0002-reading-referenced-files.md).

**Exit codes:** `0` success · `1` generic · `2` invalid args · `3` task not found · `4` task not ready · `5` already claimed · `7` lock failed

- Walkthrough: [`docs/usage.md`](docs/usage.md) — tl by example, flow by flow
- Behavioral spec: [`features/`](features) (one `.feature` file per command)
- Per-command flags: `tl <cmd> --help`

---

## How tl cli compares

`tl` shares a category with [Beads](https://github.com/steveyegge/beads) and
[Backlog.md](https://github.com/MrLesk/Backlog.md): Git-native task trackers for
humans **and** AI coding agents. The short version — `tl` is the files-only,
no-database option, and its one differentiator is **agent-safe coordination
with readable, Git-native state**: explicit claims, detectable stale work,
computable dependencies, recorded handoffs, everything inspectable by hand.

Feature-by-feature, including the honest "why `tl` and not Beads / Backlog.md /
GitHub Issues":
**[`docs/comparison.md`](docs/comparison.md)**.

---

## Development

```sh
make build                  # version-stamped local binary
make test                   # all Go tests
make bdd                    # godog suite only
make dists                  # local cross-platform archives for manual testing
make release VERSION=x.y.z  # validate, tag, and push; GitHub Actions publishes
make clean
```

CI runs `gofmt`, `go vet`, `make build`, `make test` on every PR and push to
`main` (see [`.github/workflows/ci.yaml`](.github/workflows/ci.yaml)).
`make release VERSION=x.y.z` only verifies that `HEAD` is clean, on `main`, and
already pushed to `origin/main`, then pushes the tag. The tag-triggered release
workflow builds all platform archives and publishes the GitHub Release.

---

## Further reading

- [`docs/usage.md`](docs/usage.md) — tl by example, flow by flow
- [`docs/tech-docs.md`](docs/tech-docs.md) — some implementation detail
- [`docs/PRD.md`](docs/PRD.md) — design intent, non-goals, status enum
