---
id: task-8xx
title: Register AUR account and set up SSH access for publishing tl-bin
status: open
priority: high
type: task
created_at: 2026-08-30T19:51:32Z
updated_at: 2026-09-08T19:52:13Z
created_by: pi
assignee: null
depends_on: []
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags:
  - packaging
  - aur
  - promotion
references:
  - packaging/aur/README.md
  - packaging/aur/PKGBUILD
  - packaging/aur/.SRCINFO
---

## Description

Complete the AUR account/SSH prerequisites so the tl-bin package can be published (prerequisite of the tl-bin AUR publish task).

Human-only steps (need a browser session at https://aur.archlinux.org; an agent cannot perform these):
1. Create an account at https://aur.archlinux.org/register.
2. Under **My Account**, add the SSH public key that the push will use. If no dedicated keypair exists yet, generate one with `ssh-keygen` and keep the private key available to the agent that runs the publish push.

Agent steps (only once steps 1-2 are done):
3. Verify SSH access: `ssh aur@aur.archlinux.org help`.
4. Confirm no name conflict: check https://aur.archlinux.org/packages?K=tl-bin and ?K=tl so `pkgname=tl-bin` / `provides=('tl')` / `conflicts=('tl')` are safe. If an unrelated tl/tl-bin package already exists, stop and report rather than push.

Blocking note: steps 1-2 require the human. Once the account exists and the key is registered, the agent can run steps 3-4, unblocking the publish task.

## Notes

- 2026-09-08T19:52:13Z [claude] note: Releasing this claim with --force. It is held by aho and its lease expired hours ago, so the ledger was implying someone was actively on it. Nothing here was worked on — the claim is simply stale. Context that changes this task's urgency: AUR publication is no longer the only path for Arch users. task-q32 shipped a self-hosted pacman repository (aholbreich/pacman-repo) that publishes tl-bin automatically on release, and this machine now installs from it. The AUR remains worth doing as the discoverable front door, since Arch users look there first and a third-party repository asks for more trust, but it is no longer blocking anyone.
