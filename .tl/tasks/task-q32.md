---
id: task-q32
title: Publish a self-hosted pacman repository for tl-bin
status: done
priority: high
type: task
created_at: 2026-09-08T18:36:23Z
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
  - packaging/aur/PKGBUILD
  - README.md
  - task-o9g
  - task-8xx
---

## Description

## Problem

AUR account registration has been failing for weeks, so task-o9g (publish tl-bin to the AUR) is blocked with no date. Meanwhile Arch and Omarchy users have no supported install path at all: the README promises the AUR package 'coming soon', and the only working option is building the in-repo PKGBUILD by hand, which is undocumented.

Worse, a locally built package is a dead end. pacman tracks it but nothing can upgrade it, because there is no repository to compare against. That is why this project's own machine sat on 0.9.0 while 0.12.0 shipped — structural, not drift.

## Why not Homebrew

Homebrew was considered and rejected for this audience. It is already shipped and automated (UPDATE_HOMEBREW_TAP=true), and it is the right answer on macOS. On Linux it installs into /home/linuxbrew/.linuxbrew outside pacman entirely, so an Arch user would trade a package pacman cannot upgrade for one pacman does not know exists, plus a second package manager that omarchy update also ignores.

## Proposed solution

Mirror what already works for Fedora. aholbreich/rpm-repo is a self-hosted repository on GitHub Pages, updated automatically on every release, with no distro gatekeeper. The identical pattern works for pacman and needs no AUR account.

- New repository aholbreich/pacman-repo, served from GitHub Pages, laid out by architecture so a Server line can use $arch.
- A build script alongside build-rpm.sh that produces the package from the published release artifacts.
- A release job, gated on UPDATE_PACMAN_REPO like the other two, that builds the package, runs repo-add against the existing database, and pushes.
- README documents both the repository and the makepkg fallback.

## Design constraints

- Reuse packaging/aur/PKGBUILD as the single source of packaging truth. A second definition would drift from the AUR one.
- Checksums must be regenerated for both architectures, not just the build host's, or the PKGBUILD shipped to AUR users would carry stale sums for the other arch.
- repo-add must update the existing database rather than replace it: the repository is intended to host more than one project, as rpm-repo already does.
- Gate the job on a repository variable, matching UPDATE_HOMEBREW_TAP and UPDATE_RPM_REPO, so a fork does not attempt to push.

## Relationship to task-o9g

This does not replace the AUR. Arch users look there first, and a third-party repository asks for more trust. Keep task-o9g open as the front door for whenever registration unblocks; this is the track that works now.

## Notes

- 2026-09-08T18:42:53Z [claude] note: Pipeline built and tested. Two scripts alongside build-rpm.sh: build-pacman-pkg.sh retargets packaging/aur/PKGBUILD at VERSION, refreshes sums with updpkgsums and runs makepkg; update-pacman-repo.sh copies the package into a pacman-repo checkout and runs repo-add. PKGBUILD stays the single source of packaging truth. Repository aholbreich/pacman-repo created, Pages enabled, tl-bin 0.12.0 published and verified live over HTTPS. Three defects found by testing, none of which would have surfaced from reading the code. 1. repo-add leaves .db and .files as SYMLINKS to the .tar.gz. GitHub Pages does not resolve symlinks — it serves the git blob, whose content is the target filename — so pacman would fetch holbreich.db and receive the 20-byte string 'holbreich.db.tar.gz' where a gzip archive is expected. The script now replaces them with real copies; confirmed by curling the live URL and getting 'gzip compressed data'. 2. makepkg stamps a build date into the package, so rebuilding the same tag produces a byte-different file for identical software. Re-running the release job would have republished it, changed the database checksum and forced every user to re-download for nothing. Added a guard that skips when the exact pkgver-pkgrel is already in the database. 3. repo-add leaves .old database backups that the workflow's git add -A would have committed and accumulated. Now removed. The guard itself was wrong on first write — the suffix strip left tl-bin-0.12.0 against a database entry of tl-bin-0.12.0-1, so it silently never fired. Only caught because the idempotency test still showed a modified file. Both paths are now tested: re-running the same version is a no-op, a pkgrel bump publishes. Also bumped the committed PKGBUILD from 0.9.0 to 0.12.0 — three releases stale, so anyone acting on task-o9g would have published 0.9.0 to the AUR. NOT tested: the GitHub Actions job itself. Docker is installed here but the service is inactive and starting it needs sudo, so the Arch-container path (node install before actions/checkout, non-root makepkg via a builder user) is unverified. It also cannot run until PACMAN_REPO_TOKEN and UPDATE_PACMAN_REPO are added.
- 2026-09-08T19:30:34Z [claude] note: Reinstalled on this machine from the repository. tl is now 0.12.0-1, Validated By changed from 'None' to 'SHA-256 Sum' (pacman verified the checksum from the repository database rather than trusting a local build), and 'pacman -Sl holbreich' lists it as installed — meaning it is finally known to a repository and therefore upgradeable by pacman -Syu and omarchy update. That was the structural problem the whole ticket existed to fix. Note on escalation: sudo could not prompt because the agent shell has no TTY. pkexec worked, which is what the omarchy skill prescribes for agent-launched privileged commands. Leftover: tl-bin-debug 0.9.0-1 is still installed from the old local build. The repository deliberately does not publish debug packages, so nothing will ever upgrade it, and its symbols are for a binary that is no longer installed. It should be removed.
