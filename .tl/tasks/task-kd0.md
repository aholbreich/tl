---
id: task-kd0
title: Recognize .feature references as specs, by convention
status: open
priority: medium
type: task
created_at: 2026-09-08T11:52:38Z
updated_at: 2026-09-08T11:52:38Z
created_by: claude
assignee: null
depends_on:
  - task-ps0
claim:
  actor: null
  claimed_at: null
  expires_at: null
  heartbeat_at: null
tags: []
references:
  - features/references.feature
  - internal/doctor/doctor.go
  - docs/gherkin-guidelines.md
  - cmd/show.go
  - task-ps0
---

## Description

## Problem

tl already has everything needed to connect a task to the specification that defines it, but nothing joins them up.

Already true today:
- References are generic strings and explicitly support this. features/references.feature literally uses `--ref features/login.feature` in its mixed-kinds scenario.
- `tl doctor` already validates references that look like repo-relative paths and flags dead ones (checkReferences, internal/doctor/doctor.go:396), skipping URLs and bare ticket IDs.
- tl own repository practises this: 29 feature files, one per behaviour area.

What is missing is that tl never treats such a reference as anything but an opaque string. A reader cannot ask "which tasks have a spec?" or "which spec does this task implement?".

## Proposed solution — convention, not syntax

A reference whose path ends in `.feature` is a spec reference. That is the whole rule.

- No new flag on `tl create` or `tl refine`. `--ref features/login.feature` already works and already means this.
- No new frontmatter field. No migration. Existing ledgers gain the behaviour for free the moment they use the convention.
- Projects that do not write Gherkin are unaffected — they simply have no refs matching the pattern.

Surface it in three places:

1. `tl show` groups the References block so the spec is visually distinct from code paths, URLs and ticket IDs.
2. `tl list --spec` / `--no-spec` filters to tasks that do or do not carry one — the query that makes "which features are specified?" answerable.
3. JSON exposes it as a derived convenience (for example a `spec` string alongside the untouched `references` array), so consumers do not re-implement the suffix rule.

## Design constraints

- **Optional and additive.** The convention is inert unless a project adopts it. Nothing is rejected, nothing is required, no existing output changes shape.
- **Derived, never stored.** Spec-ness is computed from the reference string. It must not become a second source of truth in the frontmatter that can disagree with the reference list.
- **Not Gherkin-specific in the plumbing.** The suffix set should be one small list, so a project using a different spec format can be accommodated later without redesign. Start with `.feature` only.
- **More than one spec ref is legal.** A task may implement scenarios in two files. Do not assume a single value.

## Open question, deliberately deferred

Whether tl should also *read* the referenced file to report its contents is a separate and much bigger question — see the decision ticket this depends on. This ticket stops at recognising the reference, which requires no file access beyond what `tl doctor` already does.

## Use case context

Reported from rssb, which follows a BDD-first workflow: every feature starts as a `.feature` file before implementation exists. Tickets and specs are one-to-one there, but the ledger cannot express or query that relationship.
