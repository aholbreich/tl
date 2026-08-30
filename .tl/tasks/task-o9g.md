---
id: task-o9g
title: Publish tl-bin initial import to AUR
status: open
priority: high
type: task
created_at: 2026-08-30T19:51:32Z
updated_at: 2026-08-30T19:51:48Z
created_by: pi
assignee: null
depends_on:
  - task-8xx
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
  - README.md
---

## Description

Push the prepared tl-bin package (0.9.0-1) to the AUR, creating the package on first push.

Steps:
1. `git clone ssh://aur@aur.archlinux.org/tl-bin.git` (initial clone may report empty).
2. Copy `packaging/aur/PKGBUILD` and `packaging/aur/.SRCINFO` into the AUR checkout.
3. Test inside the AUR checkout:
   - `makepkg --verifysource`
   - `makepkg -f`
   - `namcap PKGBUILD tl-bin-*.pkg.tar.zst` (optional, needs namcap)
4. `git add PKGBUILD .SRCINFO`, commit "Initial import: tl-bin 0.9.0-1", `git push -u origin HEAD:master`.
5. First push creates https://aur.archlinux.org/packages/tl-bin.
6. Un-mark AUR as "coming soon" in README.md so the install commands (`omarchy pkg aur add tl-bin`, `yay -S tl-bin`) are live.

Source of truth for packaging remains `packaging/aur/`; the AUR repo holds copies of PKGBUILD + .SRCINFO.
