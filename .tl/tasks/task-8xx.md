---
id: task-8xx
title: Finalize AUR registration and SSH access
status: open
priority: high
type: task
created_at: 2026-08-30T19:51:32Z
updated_at: 2026-08-30T19:51:48Z
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

Complete the AUR account/SSH prerequisites so the tl-bin package can be published.

Steps:
1. Create an account at https://aur.archlinux.org/register (human action — agent cannot create the account).
2. Add the SSH public key under **My Account**.
3. Verify SSH access with `ssh aur@aur.archlinux.org help`.
4. Confirm no name conflict: check https://aur.archlinux.org/packages?K=tl-bin and ?K=tl so `pkgname=tl-bin` / `provides=('tl')` / `conflicts=('tl')` are safe.

Blocking note: step 1 requires a human; once the account exists and the key is registered, the agent can run the SSH verification and proceed to the publish task.
