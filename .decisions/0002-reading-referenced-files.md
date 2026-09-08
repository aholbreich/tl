# 0002. Reading referenced files from read commands

**Status:** Proposed (2026-09-08)

## Context

References are generic strings pointing at related artefacts: file paths,
URLs, ticket IDs, free text. tl has always treated them as opaque. The
question is whether a read command may *open* a referenced file and report
something about its contents — for example, that a `.feature` file contains
the scenario a story names, or which tags sit on it.

This is the precedent-setting question behind the spec-linking work
(task-kd0, task-ahk) and behind any future resolver for other formats. It
should be settled before those are built rather than discovered afterwards.

The boundary has already been crossed once. `tl doctor` stats referenced
paths today (`checkReferences`, internal/doctor/doctor.go:396) and reports
dead ones. So the question is not whether tl may look outside `.tl/`, but
whether that licence extends from the diagnostic command to the everyday
read commands.

### Why it is not obviously fine

PRD §3 thesis point 7 is "Small and predictable — no daemon, no hidden
database, no automatic remote push". Reading arbitrary repo files at display
time weakens two properties:

- **Predictability.** Output would depend on working-tree state outside the
  ledger. The same ledger at the same commit renders differently with a
  dirty tree.
- **Cost and failure modes.** Every listed task potentially becomes a file
  read. Missing, unreadable or malformed files need defined behaviour rather
  than an error that breaks the whole listing.

### Why it may be fine anyway

- **Repo-local and deterministic.** A file read in the same working tree is
  not a network call, a subprocess or a daemon. It is far closer to what
  `doctor` already does than to anything on the non-goals list.
- **The alternative is worse.** The other way to know a spec's state is to
  store it in the ledger — derived state that goes stale silently the moment
  someone edits the feature file. That contradicts "human-readable" and
  creates exactly the drift tl exists to prevent.
- **It is the one thing tl can uniquely do.** tl sits in the same repository,
  at the same commit, as the specs. A hosted tracker can link to a spec but
  can never tell you whether it matches. This is not competing with Jira on
  feature richness (PRD §4); it is the structural advantage of being
  repo-local.

## Decision

**Proposed:** read commands may open referenced files, under four limits.

1. **The reference is the trigger.** tl performs no project-level detection —
   no `features/` directory convention, no build-tool sniffing, no mode flag.
   A file is read only because a task points at it. A project that
   references no specs sees no reads and no output change, automatically.

2. **Two cost tiers.** Existence is a `stat` and may be unconditional, since
   `doctor` already pays that cost. Parsing file *contents* happens only for
   references that carry a scenario anchor, or when explicitly requested by
   a flag. Reads are bounded by the number of tasks carrying spec
   references, never by ledger size.

3. **Automatic in human output, explicit in JSON.** `show`, `list` and `tree`
   may enrich their human-facing output silently. `--json` keys are always
   present and null when there is nothing to say, so a consumer's schema
   never changes because someone added a feature file.

4. **Structure, never semantics.** tl reports what it can observe — the file
   exists, it holds N scenarios, the named scenario resolves, these tags are
   present. tl assigns no meaning to any tag and never concludes that work
   is "done". Interpreting the vocabulary belongs to the consumer.

### Consequences

**Accepted costs:**

- Output depends on working-tree state, not the ledger alone. This is a
  stated property: the divergence between ledger and spec is the signal, not
  a defect.
- Read commands acquire a failure mode they did not have. Every read must
  degrade to "unknown" rather than erroring; a listing must never fail
  because a referenced file was deleted.
- tl learns one spec format's grammar. The interface is named for reference
  resolution, not Gherkin, so a second resolver does not require a redesign —
  but the first one still has to be written.

**Avoided costs:**

- No cached or stored spec state in frontmatter that can disagree with the
  file.
- No project configuration, no detection heuristic, no enable flag.
- No test execution, no build tooling, no report ingestion (PRD §4).

## Alternatives considered

1. **Never — recognition only.** Consumers join spec state themselves from
   `--json` plus their own test runner. Rejected: it leaves the one join tl
   is uniquely placed to make undone, and the per-task subprocess cost is
   what stops such tooling being written at all.

2. **`doctor` only.** Extend the existing check; read commands stay pure.
   The conservative fallback, and `doctor` is already the "look at the
   world" command. Rejected as the primary answer because a diagnostic
   command is the wrong place for everyday queries, but it remains the
   safe retreat if limit 2 proves too costly in practice.

3. **Opt-in flag on read commands only.** All reads behind an explicit
   `--spec-status`. Rejected as too weak: existence is already paid for by
   `doctor` at the same cost, and hiding it behind a flag means the default
   view stays uninformative for no saving. Retained for the parsing tier.

4. **Always, everywhere.** Spec status wherever a spec reference appears,
   including bulk JSON by default. Rejected: it makes every listing pay for
   a feature most projects will not use, and it changes JSON shape based on
   working-tree state.

## Open questions

- Does a *missing* referenced file mean "no spec" or "broken link"? `doctor`
  calls it a warning today; an enriched listing must render it as one or the
  other, and they carry opposite meanings.
- What does `tl tree` render for a node whose spec is unresolved, and should
  that propagate to an ancestor whose own spec resolves?
- Which resolver is second? Not to build it — to check the interface is not
  accidentally Gherkin-shaped.

## References

- PRD `docs/PRD.md` §3 (Core thesis: small and predictable)
- PRD `docs/PRD.md` §4 (Non-goals: not a Jira / Linear replacement; no
  background workers; no agent execution)
- Decision `0001` (encode dimensions as convention, not schema fields)
- `internal/doctor/doctor.go:396` — the existing precedent for reading
  outside `.tl/`
- task-thh, task-kd0, task-ahk — the work gated on this decision
