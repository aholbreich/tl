---
id: task-ejg
title: Use make changelog for GitHub release notes
status: done
priority: medium
type: task
created_at: 2026-05-31T20:37:50Z
updated_at: 2026-05-31T20:38:35Z
created_by: pi:release-notes
assignee: null
depends_on: []
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags:
  - release
  - ci
references:
  - Makefile
  - .github/workflows/release.yaml
---

## Description

Make the release workflow publish the curated make changelog output as the GitHub Release body. Support tag-aware changelog ranges so tag-triggered workflows compare the previous tag to the pushed tag.

## Notes

- 2026-05-31T20:38:35Z [pi:release-notes] note: Implemented tag-aware release notes: make changelog accepts CURRENT_TAG and compares previous tag to that tag; release workflow checks out full history, writes release_notes.md, and passes it to softprops/action-gh-release via body_path with generated notes disabled. Validation: make changelog, make changelog CURRENT_TAG=0.8.1, git diff --check, make test passed.
