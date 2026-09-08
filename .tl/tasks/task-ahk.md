---
id: task-ahk
title: Surface spec status from the feature file's own Gherkin tags
status: open
priority: medium
type: task
created_at: 2026-09-08T11:53:17Z
updated_at: 2026-09-08T11:53:17Z
created_by: claude
assignee: null
depends_on:
  - task-kd0
  - task-1re
  - task-hwb
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags: []
references:
  - docs/gherkin-guidelines.md
  - features/references.feature
  - task-kd0
  - task-1re
---

## Description

## Problem

A task overview that answers "what is specified, and is it built yet?" currently has to be assembled by hand from two sources: tl for ticket status, and the test runner for whether the scenarios pass.

The insight is that tl own repository already solved half of this, informally: all 29 feature files carry a feature-level `@implemented` tag. That is the entire tag vocabulary in use. The author already records spec status — in the Gherkin, where it belongs, next to the thing it describes.

tl does not read it, so nobody can query it.

## Proposed solution

Where a task carries a spec reference (task-kd0), read the referenced file feature-level tags and report them.

```
tl list --spec-status

ID        Status  Spec                            Scenarios
task-ps0  open    features/list.feature           @implemented
task-7fi  open    features/tree.feature           (untagged)
task-kd0  open    -                               -
```

The vocabulary is the project, not tl. tl reports the tags it finds; it does not define what `@implemented` means, does not require it, and does not validate it against a list. A project using `@wip` or `@draft` gets those reported instead, with no code change.

## Why tags rather than parsing scenarios

Counting scenarios is tempting but tells you less than it appears. A file with seven scenarios says nothing about whether they pass — only the test runner knows that, and teaching tl to run test suites is squarely a non-goal ("AI agent execution itself", "long-running background workers", and by extension any build tooling).

A tag is a claim the author makes deliberately, in version control, reviewable in a diff. That is the same contract as every other field in the ledger. Scenario counts could be added later as a cheap extra; the tag is the part that carries meaning.

## Design constraints

- **Gated on task-1re.** This ticket assumes tl may read files outside `.tl/` on a read command. If that decision lands as "doctor only" or "never", this ticket must be rewritten or cancelled rather than quietly implemented.
- **Never on the default path.** Plain `tl list` must not become one file read per row.
- **Degrade quietly.** A missing, unreadable, or tagless file reports "unknown", not an error. `tl doctor` already owns complaining about dead references; a listing must not fail because someone deleted a feature file.
- **Parse only the header.** Feature-level tags are the lines immediately preceding `Feature:`. Read that far and stop — no full Gherkin parser, no new dependency.
- **Report, do not interpret.** No hardcoded meaning for any tag string.

## Prerequisite

The `@implemented` convention is currently undocumented — see the guidelines ticket. It should be written down before tl builds behaviour on top of it.
