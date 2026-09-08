---
id: task-8xx
title: Register AUR account and set up SSH access for publishing tl-bin
status: in_progress
priority: high
type: task
created_at: 2026-08-30T19:51:32Z
updated_at: 2026-09-08T12:28:45Z
created_by: pi
assignee: null
depends_on: []
claim:
  actor: aho
  claimed_at: 2026-09-08T12:28:45Z
  expires_at: 2026-09-08T13:28:45Z
  heartbeat_at: 2026-09-08T12:28:45Z
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
