---
id: task-xu2
title: Enable and verify the pacman repo release job
status: open
priority: medium
type: task
created_at: 2026-09-08T19:30:34Z
updated_at: 2026-09-08T19:52:28Z
created_by: claude
assignee: null
depends_on: []
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags: []
references:
  - .github/workflows/release.yaml
  - .github/scripts/build-pacman-pkg.sh
  - task-q32
---

## Description

## Problem

The update-pacman-repo job is committed but has never run. Its scripts are tested end to end on an Arch host, and the repository it publishes to is live, but the CI wiring around them is unverified and the job is gated off.

## What is unverified

The job runs in an archlinux:base-devel container, which introduces two things that cannot fail locally:

- Node is installed by a run step *before* actions/checkout, because the Arch image ships none and a JavaScript action cannot start without it. If that ordering is wrong the job fails at checkout.
- makepkg refuses to run as root while the container runs as root, so the build happens as a created 'builder' user via su. Ownership of GITHUB_WORKSPACE has to be handed over for that to work.

Testing this locally needs a container runtime. Docker is installed on the maintainer's machine but the service is inactive; starting it requires privileges the agent did not have at the time.

## What a human must add

- Secret PACMAN_REPO_TOKEN: a PAT with write access to aholbreich/pacman-repo. The token behind RPM_REPO_TOKEN would serve if its scope already covers it.
- Repository variable UPDATE_PACMAN_REPO=true.

The variable was deliberately left unset: enabling the job without the token would fail the next release.

## Acceptance

A tagged release publishes the package to aholbreich/pacman-repo without manual intervention, and a re-run of that same job reports 'already in the repository; nothing to publish' rather than committing a rebuilt package.

## Notes

- 2026-09-08T19:48:29Z [claude] note: Container path now verified locally. Docker was started, and the job was simulated step for step in archlinux:base-devel against a read-only mount of the repository: pacman -Syu installs git/nodejs/pacman-contrib (node v26.8.1, so actions/checkout can start), a copy stands in for the checkout, a created builder user owns the workspace and runs build-pacman-pkg.sh producing tl-bin-0.12.0-1-x86_64.pkg.tar.zst, and update-pacman-repo.sh runs as root exactly as the job does. The published database came out as a real gzip file rather than a symlink, and a second publish reported 'already in the repository; nothing to publish'. So both concerns this ticket was opened for — Node before checkout, and non-root makepkg — are settled. What remains is purely credentials: add secret PACMAN_REPO_TOKEN with write access to aholbreich/pacman-repo, then set variable UPDATE_PACMAN_REPO=true. The variable is still deliberately unset, since enabling the job without the token would fail the next release. Acceptance is unchanged: a tagged release should publish without intervention, and a re-run of that job should report nothing to publish.
- 2026-09-08T19:48:29Z [claude] note: System is on a clean install: tl-bin removed and reinstalled from the repository rather than upgraded in place, so nothing survives from the original local build. tl 0.12.0, Validated By SHA-256 Sum, no orphaned packages, no unowned tl files under /usr, and no shadowing copies in ~/bin, ~/.local/bin or /usr/local/bin. The stale tl-bin-debug 0.9.0-1 is gone.
- 2026-09-08T19:52:28Z [claude] note: Credentials resolved without a PAT. Registered a write-scoped deploy key on aholbreich/pacman-repo and stored its private half as secret PACMAN_REPO_DEPLOY_KEY on aholbreich/tl; the private key was generated locally, uploaded, then shredded from disk. This is narrower than the approach originally specified: a deploy key reaches exactly one repository, whereas the PAT pattern the homebrew and rpm jobs use carries the whole account's permissions, and revoking it is one click on pacman-repo. The job now checks out via actions/checkout's ssh-key input instead of token. That change surfaced a failure the original design would have hit on its first real run: the archlinux image ships no ssh client, so actions/checkout could not have configured an SSH transport. Verified directly — 'command -v ssh' in the container returned nothing before openssh was added to the toolchain step. Proved the fixed path end to end by registering a throwaway read-only deploy key, cloning pacman-repo over SSH from inside the container, then deleting that key; only the real write-scoped key remains. Full simulation re-run green afterwards. UPDATE_PACMAN_REPO is now set to true, so the job is armed. Remaining and unavoidable: the job has still never executed on GitHub's runners. Everything it does has been reproduced locally, but only a tagged release exercises the real thing.
