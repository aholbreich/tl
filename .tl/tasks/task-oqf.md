---
id: task-oqf
title: Add a demo (asciinema/gif) to top of README
status: done
priority: high
type: feature
created_at: 2026-05-29T16:52:29Z
updated_at: 2026-05-29T19:09:52Z
created_by: human
assignee: null
depends_on: []
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags:
  - docs
  - promotion
---

## Description

Highest-leverage missing asset. Record a ~20s flow: tl init -> create -> ready --json -> claim -> note -> close, ideally showing an agent + human handoff. Embed a gif/asciinema near the top of the README so visitors see it working in 5 seconds. Star conversion driver #1.

## Notes

- 2026-05-29T19:09:52Z [pi:promotion] note: Created animated terminal demo SVG (.github/tl-demo.svg) showing the full workflow: init → create (2 tasks) → dep add → show → ready --json → claim → note → close → list --all → history. Embedded at top of README under the badges. Recorded with asciinema, rendered with svg-term-cli into a self-contained animated SVG that plays on page load. ~20s loop, no external dependencies.
