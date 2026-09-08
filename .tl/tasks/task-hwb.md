---
id: task-hwb
title: Document the @implemented tag convention in the Gherkin guidelines
status: open
priority: low
type: task
created_at: 2026-09-08T11:53:31Z
updated_at: 2026-09-08T11:53:31Z
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
  - docs/gherkin-guidelines.md
  - features/references.feature
---

## Description

## Problem

Every one of the 29 files under features/ opens with `@implemented`, and it is the only tag in the repository. It clearly means "this specification is backed by working code".

docs/gherkin-guidelines.md — which describes itself as a context contract for humans and LLMs, and covers files, formatting, vocabulary, feature blocks, backgrounds, scenarios, steps, tables, anti-patterns and a checklist — never mentions tags at all.

So a contributor or agent reading the guidelines has no idea the convention exists, and one reading the feature files has no idea what the tag obliges them to do. New files may or may not get it, and nothing says when it should be added or removed.

## Proposed solution

Add a short "Tags" section to docs/gherkin-guidelines.md in the existing MUST/SHOULD/MAY voice, stating:

- What `@implemented` asserts.
- That it is feature-level, on the line directly above `Feature:`.
- When it is added — with the implementation, not with the spec, since BDD-first means the file exists first and is deliberately untagged while red.
- When it is removed — if the behaviour is dropped or the spec outgrows the code again.
- That the tag vocabulary is open: projects may add their own, and tl does not define them.

Add the tag to the recommended template at the foot of the document, or show both the tagged and untagged forms, so the two states are visible.

## Why it matters beyond tidiness

task-ahk proposes having tl report these tags. Building tooling on an undocumented convention is how conventions rot: the first contributor who omits the tag is not wrong, because nothing told them. Write it down first.

## Design constraint

Describe the existing practice. Do not invent a wider tag taxonomy in the same change — if more tags are wanted, that is a separate discussion with its own ticket.
