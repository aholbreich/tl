---
id: task-qf4
title: Add install.sh for release binaries
status: done
priority: medium
created_at: 2026-05-23T20:58:30Z
updated_at: 2026-05-23T20:59:41Z
created_by: human
assignee: null
depends_on: []
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags:
  - installer
---

## Description

Create a shell installer that detects OS/architecture, downloads the matching tl release archive from GitHub, installs the tl binary, and document usage in README.

## Notes

### 2026-05-23T20:59:41Z - pi:install-script

Implemented install.sh for GitHub release binaries. It detects linux/darwin and amd64/arm64, supports --version and --bin-dir, uses curl/wget, installs tl, and README documents curl | sh usage. Verified with sh -n, shellcheck, go test ./..., and installing release 0.4.4 into a temp directory.
