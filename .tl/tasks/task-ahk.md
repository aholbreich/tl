---
id: task-ahk
title: Report feature-level spec tags on the read commands
status: open
priority: medium
type: task
created_at: 2026-09-08T11:53:17Z
updated_at: 2026-09-08T22:22:54Z
created_by: claude
assignee: null
depends_on:
  - task-kd0
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
  - .decisions/0002-reading-referenced-files.md
---

## Description

## Problem

A reader asking "what is specified, and what does the author say about it?" currently has to assemble the answer from two sources: tl for ticket status, and the repository for what the spec claims about itself.

Once a spec reference is recognised (task-kd0), the tags on that feature are already there — in version control, reviewable in a diff, written next to the thing they describe. tl does not read them, so nobody can query them.

## Proposed solution

Under `tl list --spec-status`, report the feature-level tags of each referenced spec file alongside its existence.

```
tl list --spec-status

ID        Status  Spec                       Tags
task-ps0  done    features/list.feature      @implemented
task-7fi  open    features/tree.feature      (missing)
task-cys  open    features/agents.feature    @implemented
task-wke  open    -                          -
```

The vocabulary is the project's, not tl's. tl reports the tags it finds; it does not define what `@implemented` means, does not require it, and does not validate it against a list. A project using `@wip`, `@draft` or `@manual` gets those reported instead, with no code change.

## What a tag is and is not evidence of

A tag is a **claim the author made deliberately**, in version control. That is the same contract as every other field in the ledger, and it is worth surfacing.

It is not evidence that the scenarios pass. Only a test runner knows that, and teaching tl to run suites is squarely a non-goal — "AI agent execution itself", "long-running background workers", and by extension any build tooling. Ingesting a test report is a much larger commitment and must not be smuggled in through this ticket.

Nor is it evidence about *this story*. At file granularity a tag describes the whole feature file, which usually serves several stories; a story sharing a mature file with delivered work will show that file's tags while being unbuilt. That ceiling is decision 0002's accepted trade — see task-kd0.

## Design constraints

- **Feature-level tags only.** The lines immediately preceding `Feature:`. Not scenario-level: at file granularity there is no way to know which scenario belongs to this story, so reporting scenario tags would imply a precision tl does not have.
- **Same read as task-kd0.** Existence and tags come from one open of the file. Do not add a second pass.
- **Parse only the header.** Read as far as the `Feature:` line and stop. No full Gherkin parser, no new dependency.
- **Never on the default path.** Plain `tl list` must not become one file read per row.
- **Degrade quietly.** Missing, unreadable or tagless reports "unknown" or an empty list, never an error. `tl doctor` owns complaining about dead references; a listing must not fail because someone deleted a feature file.
- **Report, do not interpret.** No hardcoded meaning for any tag string, `@implemented` included.
- **Tags are a list.** A feature may carry several. Do not collapse to one.

## On the task-hwb prerequisite

task-hwb documents tl's own `@implemented` convention. It is worth doing, but it is not a technical prerequisite: since tl assigns no meaning to any tag, nothing here depends on that convention being written down. Keep the dependency for dogfood tidiness, or drop it if it blocks progress.

## Problem

A task overview that answers "what is specified, and how far along is it?" currently has to be assembled by hand from two sources: tl for ticket status, and the repository for what the spec says about itself.

Once a spec reference resolves to a file or a scenario (task-kd0), the tags sitting on it are already there, already in version control, already reviewable in a diff. tl does not read them, so nobody can query them.

## Proposed solution

Where a task carries a spec reference, report the tags found on the resolved unit — feature-level tags for a file reference, scenario-level tags for an anchored one.

```
tl list --spec-status

ID        Status  Spec                                          Tags
task-ps0  open    features/list.feature                         @implemented
task-7fi  open    features/tree.feature#Rendering a subtree     (untagged)
task-kd0  open    -                                             -
```

The vocabulary is the project's, not tl's. tl reports the tags it finds; it does not define what `@implemented` means, does not require it, does not validate it against a list. A project using `@wip`, `@manual` or `@ignore` gets those reported instead, with no code change.

## Tags and scenario resolution answer different questions

An earlier draft of this ticket argued for tags *instead of* scenario resolution. That was the wrong frame — they are orthogonal, and neither is what a test runner gives you:

| Signal | Answers | Does not answer |
|---|---|---|
| Scenario resolves (task-kd0) | does the spec for *this story* exist yet | whether it passes |
| Tags on that unit (this ticket) | what the author asserts about it | whether the assertion is true |
| Test run (out of scope) | whether it passes | — |

Resolution is granular but silent about intent; a tag is a deliberate claim but says nothing about which story it covers. Together they answer "this story's scenario exists and the author has marked it ready"; separately, neither does.

That tl cannot report pass/fail is a boundary, not a gap. Running a suite is squarely a non-goal — "AI agent execution itself", "long-running background workers", and by extension any build tooling. Ingesting a test report is a much larger commitment and should not be smuggled in through this ticket.

## Design constraints

- **Gated on decision 0002.** This assumes read commands may open files outside `.tl/`. If that lands as "doctor only" or "never", this ticket is rewritten or cancelled, not quietly implemented.
- **Never on the default path.** Plain `tl list` must not become one file read per row.
- **Degrade quietly.** Missing, unreadable or tagless resolves to "unknown", never an error. `tl doctor` owns complaining about dead references; a listing must not fail because someone deleted a feature file.
- **Parse only as far as needed.** Feature-level tags are the lines immediately preceding `Feature:`; scenario-level tags precede `Scenario:`. No full Gherkin parser, no new dependency.
- **Report, do not interpret.** No hardcoded meaning for any tag string, `@implemented` included.
- **Tags are a list, not a value.** A scenario may carry several. Do not collapse to one.

## On the task-hwb prerequisite

task-hwb documents tl's own `@implemented` convention. It is worth doing, but it is no longer a hard prerequisite for this ticket: since tl assigns no meaning to any tag, nothing here depends on that convention being written down. Keeping the dependency is a choice about dogfood tidiness, not a technical constraint — drop it if it blocks progress.

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

## Notes

- 2026-09-08T22:22:54Z [claude] note: Retitled: the --spec-status flag it named no longer exists. Spec state is now resolved unconditionally on list, ready and tree (decision 0002, amended), so this ticket's tags would appear on the ordinary read commands rather than behind a flag. Nothing else about its scope changes.
