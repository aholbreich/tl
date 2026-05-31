---
id: task-phl
title: 'Launch checklist: distribution channels'
status: open
priority: medium
type: feature
created_at: 2026-05-29T16:52:45Z
updated_at: 2026-05-31T21:11:50Z
created_by: human
assignee: null
depends_on:
  - task-z31
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags:
  - promotion
references:
  - README.md
  - docs/comparison.md
  - .github/workflows/release.yaml
---

## Description

GATED on green CI + demo + 'why tl' section landing first (one shot per channel). Sequence: (1) awesome-go + awesome-cli-apps PRs, (2) Show HN with the demo gif, (3) r/golang, r/commandline, (4) Lobste.rs, (5) agentic-coding communities — Claude Code / Cursor / aider Discords & subreddits (the real target audience). Track which channel converted. Do NOT fire before the storefront converts.

## Notes

- 2026-05-31T21:11:38Z [pi:promotion-plan] note: Refined launch checklist for later execution. Gates before promotion: 1. Green CI on main. 2. Latest GitHub Release looks good: binaries attached, release notes use curated `make changelog`, install paths verified. 3. README storefront is ready (task-z31): lead with agent-safe coordination, clear install path, `tl init && tl agents --write-files`, demo visible. 4. Comparison page remains sharp and honest for Beads / Backlog.md positioning. 5. Smoke-test install flow: `brew install aholbreich/tap/tl`, `tl init`, `tl agents --compact`. Execution waves: 1. Curated lists: awesome-go, awesome-cli-apps, and possibly agent-tooling lists. Angle: Git-native task ledger for humans and AI coding agents; Markdown tasks, explicit claims, stale-work detection, dependency-aware ready queue, JSON output. 2. Show HN after README is strong. Suggested title: "Show HN: tl — a Git-native task ledger for humans and coding agents". Keep post short: problem, why not Issues/TODO/chat, differentiator, install command, demo GIF. 3. Reddit / Lobsters, staggered over days: r/golang, r/commandline, maybe r/selfhosted; Lobsters tags go/tools/release. Use "I built this because..." tone, not a marketing blast. 4. Agent communities: Claude Code, Cursor, aider, Copilot/agentic-coding communities. Angle: shared local ledger for multi-agent/human repos with claims, notes, stale leases, JSON-readable state. Tracking per channel: - Add a note with URL, date/time, post title, and signal. - Track stars, GitHub traffic/referrers, release downloads, Homebrew/RPM installs if available, issues opened, and useful comments/questions. Do not fire all channels at once; use one-shot posts and adapt copy based on conversion/comments.
