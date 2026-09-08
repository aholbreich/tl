---
id: task-xu2
title: Enable and verify the pacman repo release job
status: open
priority: medium
type: task
created_at: 2026-09-08T19:30:34Z
updated_at: 2026-09-08T19:30:34Z
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
