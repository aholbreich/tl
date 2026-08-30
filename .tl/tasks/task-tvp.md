---
id: task-tvp
title: Test AUR build and installation end-to-end
status: open
priority: medium
type: task
created_at: 2026-08-30T19:51:32Z
updated_at: 2026-08-30T19:51:48Z
created_by: pi
assignee: null
depends_on:
  - task-o9g
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags:
  - packaging
  - aur
  - promotion
  - testing
references:
  - packaging/aur/README.md
  - packaging/aur/PKGBUILD
  - packaging/aur/.SRCINFO
---

## Description

Verify the published tl-bin package builds cleanly and installs correctly on a fresh machine.

Build tests (from a clean clone of the AUR repo):
- `makepkg -f` on x86_64 (already verified once locally for 0.9.0-1; re-verify from the published repo).
- aarch64: `makepkg --printsrcinfo` and `makepkg --verifysource` to confirm the arm64 tarball + LICENSE checksums resolve; full build if an arm64 machine is available.
- `namcap` lint clean (or document accepted warnings).

Install tests:
- `sudo pacman -U tl-bin-0.9.0-1-x86_64.pkg.tar.zst` (direct archive install).
- AUR helper path: `yay -S tl-bin` and/or `omarchy pkg aur add tl-bin`.
- Verify `provides=('tl')` / `conflicts=('tl')` behave (no other `tl` package can be co-installed).

Post-install verification:
- `pacman -Ql tl-bin` lists /usr/bin/tl, /usr/share/licenses/tl-bin/LICENSE, and bash/zsh/fish completion files.
- `tl --version` runs.
- `tl completion --install` succeeds and completions load in bash, zsh, and fish.

Record results as a note; any packaging fixes go back into `packaging/aur/` and a version bump to the AUR repo.
