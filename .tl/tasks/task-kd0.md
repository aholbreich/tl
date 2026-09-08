---
id: task-kd0
title: Recognize .feature references as specs, and surface them under --spec-status
status: open
priority: medium
type: task
created_at: 2026-09-08T11:52:38Z
updated_at: 2026-09-08T13:34:08Z
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
  - .decisions/0001-multi-agent-coordination-via-tags.md
  - .decisions/0002-reading-referenced-files.md
---

## Description

## Problem

tl already has everything needed to connect a task to the specification that defines it, but nothing joins them up.

Already true today:
- References are generic strings and explicitly support this. `features/references.feature` uses `--ref features/login.feature` in its mixed-kinds scenario.
- `tl doctor` already validates path-shaped references and flags dead ones (checkReferences, internal/doctor/doctor.go:396), skipping URLs and bare ticket IDs.

What is missing is that tl never treats such a reference as anything but an opaque string. A reader cannot ask "which stories have a spec?" or "which spec links are dead?".

## Proposed solution — convention, not syntax

A reference whose path ends in `.feature` is a spec reference. That is the whole rule.

- No new flag on `tl create` or `tl refine`. `--ref features/login.feature` already works and already means this.
- No new frontmatter field, no migration. Existing ledgers gain the behaviour the moment they use the convention.
- No project-level detection: no `features/` directory assumed, no build-tool sniffing, no mode. The reference is the trigger, so a project that writes no Gherkin has no matching references and sees no change anywhere. A monorepo with specs in six places works without configuration.

Surface it in three places:

1. **`tl show`** groups the References block so a spec is visually distinct from code paths, URLs and ticket IDs. String-only; no file access.
2. **`tl list --spec-status`** adds a spec column: the referenced file, whether it exists, and what the header says about it. This is the only path that opens a file, per decision 0002.
3. **`--json`** exposes a `spec` key, always present and null when the flag was not passed, so a consumer's schema never depends on which flags were used.

## Granularity — the accepted ceiling

A feature file describes a capability; a story is usually a slice of one. So this reports "this story points at a spec file that exists and says X", not "this story's behaviour is specified". A story sharing a mature feature file with delivered work will show that file's tags while being unbuilt.

That is decision 0002's accepted trade, not an oversight. The two ways past it — scenario anchors in the reference, and ledger identifiers tagged on scenarios — were both considered and rejected there. Do not reintroduce either in this ticket.

## Design constraints

- **Optional and additive.** Inert unless a project adopts the convention. Nothing rejected, nothing required, no existing output changes shape.
- **Derived, never stored.** Spec-ness is computed from the reference string. It must never become a frontmatter field that can disagree with the reference list.
- **Never on the default path.** Plain `tl list` must not become one file read per row. Only `--spec-status` reads.
- **Not Gherkin-specific in the plumbing.** The suffix set should be one small list so another spec format can be accommodated later without redesign. Start with `.feature` only, and name the concept reference resolution rather than anything Gherkin-shaped.
- **Degrade quietly.** A missing, unreadable or malformed file yields "unknown", never an error that breaks a listing. `tl doctor` owns complaining about dead references.
- **More than one spec reference is legal.** Uncommon, but a story may point at two files. Do not assume a single value; do not design the display around the plural case either.

## Prerequisite

Decision 0002 — accepted. This ticket implements it.

## Use case context

Reported from rssb, which follows a BDD-first workflow: every feature starts as a `.feature` file before implementation exists. Tickets and specs are close to one-to-one there, which is the case where file granularity is most informative.

## Problem

tl already has everything needed to connect a task to the specification that defines it, but nothing joins them up.

Already true today:
- References are generic strings and explicitly support this. `features/references.feature` uses `--ref features/login.feature` in its mixed-kinds scenario.
- `tl doctor` already validates path-shaped references and flags dead ones (checkReferences, internal/doctor/doctor.go:396), skipping URLs and bare ticket IDs.

What is missing is that tl never treats such a reference as anything but an opaque string. A reader cannot ask "which tasks have a spec?" or "which spec does this task implement?".

## Why a file-level link is not enough

The naive rule — "a reference ending in `.feature` is this task's spec" — breaks on granularity that is general to BDD projects, not specific to any one ledger. A feature file describes a capability; a story is a slice of one. The relationship is one file to many stories.

In tl's own ledger three open, unbuilt tickets (agents --remove, agents --output, Cursor rules support) all reference `features/agents.feature`, a file of fourteen scenarios describing work that is largely already delivered. Nothing about that file can tell you the state of any one of those three stories, because the file is not about any one of them.

So a file-level link answers "is there a spec nearby?" — useful, but weaker than it looks, and actively misleading if rendered as delivery state.

## Proposed solution — the reference is the detector

Two parts.

**1. Recognition needs no project detection.** A reference whose path part ends in `.feature` is a spec reference. That is the whole rule. There is no `features/` directory convention to assume, no build-tool sniffing, no project-level mode. A project that writes no Gherkin simply has no matching references and sees no change anywhere. A monorepo with specs in six places works without configuration. Per-task, not per-project: in a ledger where three stories of two hundred are spec'd, only those three light up.

**2. A reference may name a scenario, not just a file.**

```
tl create "Strip managed blocks from agent files" \
  --ref "features/agents.feature#Removing the managed block from AGENTS.md"
```

This makes the link story-granular, and it makes the link *checkable*: whether the named scenario is present in the file is a fact about the file, requiring no tag and no convention. It also matches the BDD-first order this project already prescribes — the scenario is written before the code, so the reference can be attached at refinement time and resolves from red to green without ever being edited.

## What tl reports, and what it must not decide

tl parses structure and reports it. tl assigns no meaning.

| tl reports | tl does not decide |
|---|---|
| the referenced file exists, or is missing | whether the spec is "done" |
| it holds N scenarios | whether N is enough |
| the named scenario resolves, or is absent | whether it passes |
| the tags present, verbatim, as data | what any tag means |

A team filtering on `@wip`, `@manual` or `@ignore` gets the same machinery as one using `@implemented`, with no change to tl. Interpreting the vocabulary is the consumer's job. Reporting the tags is the follow-up ticket; this one owns resolution.

## Cost, and where the work happens

Two tiers, because they cost very different amounts:

- **Existence** — one `stat` per referenced path. `tl doctor` already does exactly this. Cheap enough to be unconditional.
- **Scenario resolution** — parse the file header and scenario names. Needed only for references that carry an anchor, or when explicitly requested. Bounded by the number of tasks carrying spec references, not by ledger size.

This suggests the split: automatic in human-facing output (`show`, `list`, `tree` — already for people), explicit and schema-stable in `--json`, where the key is always present and null when there is nothing to say, so a consumer never breaks because someone added a feature file.

## Design constraints

- **Optional and additive.** Inert unless a project adopts it. Nothing is rejected, nothing required, no existing output changes shape.
- **Derived, never stored.** Spec-ness is computed from the reference string and the file on disk. It must never become a frontmatter field that can disagree with them.
- **Name the concept generically.** This is *reference resolution*, with Gherkin as the first resolver. OpenAPI paths or ADR status fields should slot in later without a redesign. Implement only Gherkin now; do not name the interface after it.
- **No full Gherkin parser and no new dependency.** Feature-level tags are the lines above `Feature:`; scenario names are `Scenario:` / `Scenario Outline:` lines. Read that much.
- **Degrade quietly.** A missing, unreadable or malformed file yields "unknown", never an error that breaks a listing. `tl doctor` owns complaining about dead references.
- **More than one spec reference is legal.** A story may span scenarios in two files. Do not assume a single value.

## Prerequisites

- task-thh — `tl doctor --fix` deletes anchored references today. Must land first or adopting the convention destroys data.
- Decision 0002 — whether read commands may open files outside `.tl/` at all, and how far. If that lands as "doctor only" or "never", the anchor syntax and part 1 still stand; only the resolution half is affected.

## Use case context

Reported from rssb, which follows a BDD-first workflow: every feature starts as a `.feature` file before implementation exists. Tickets and specs are close to one-to-one there, which is what makes the file-level rule look sufficient until a ledger with coarser feature files hits it.

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


