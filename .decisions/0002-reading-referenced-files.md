# 0002. Reading referenced files from read commands

**Status:** Accepted (2026-09-08)

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

Read commands may open referenced files, **only behind an explicit flag**,
and only at file granularity.

1. **The reference is the trigger.** tl performs no project-level detection —
   no `features/` directory convention, no build-tool sniffing, no mode flag.
   A file is read only because a task points at it. A project that
   references no specs sees no reads and no output change, automatically.

2. **Explicit flag, never the default path.** Plain `tl list`, `tl ready`,
   `tl show` and `tl tree` read nothing outside `.tl/`. Reads happen when
   the reader asks: `tl list --spec-status`. The default-path guarantee — a
   listing costs one directory of Markdown and nothing else — is preserved
   exactly.

3. **File granularity only.** A spec reference names a file. tl does not
   parse a scenario name out of the reference, and does not look for a
   ledger identifier inside the spec. Both were considered and rejected
   (alternatives 5 and 6); each buys story-level precision at a cost the
   project is not willing to pay.

4. **Schema-stable JSON.** The `spec` key is always present in `--json`
   output and null when the flag was not passed, so a consumer's schema
   never depends on which flags were used or on whether a project writes
   Gherkin.

5. **Structure, never semantics.** tl reports what it can observe — the file
   exists, it holds N scenarios, these feature-level tags are present. tl
   assigns no meaning to any tag and never concludes that work is done.
   Interpreting the vocabulary belongs to the consumer.

### What this deliberately does not deliver

A feature file describes a capability; a story is usually a slice of one.
At file granularity tl cannot say whether *this story's* behaviour is
specified — only that a spec file the story points at exists and carries
certain tags. A story sharing a mature feature file with delivered work will
show that file's tags while being unbuilt.

That is the accepted ceiling of this decision, not an oversight. Delivery
state remains the ledger's own status; the spec column sits beside it, and
the reader draws the join.

### Consequences

**Accepted costs:**

- Spec information is coarse. See above.
- Output under the flag depends on working-tree state, not the ledger alone.
  The divergence between ledger and spec is the signal, not a defect.
- Read commands acquire a failure mode they did not have. Under the flag,
  every read must degrade to "unknown" rather than erroring; a listing must
  never fail because a referenced file was deleted.
- tl learns one spec format's grammar. The interface is named for reference
  resolution, not Gherkin, so a second resolver does not require a redesign —
  but the first one still has to be written.

**Avoided costs:**

- No cached or stored spec state in frontmatter that can disagree with the
  file.
- No project configuration, no detection heuristic, no enable flag.
- No new reference syntax, so no change to how `tl doctor` validates
  references and no risk of `--fix` destroying links.
- No ledger identifiers written into specification files.
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

3. **Always, everywhere.** Spec status wherever a spec reference appears,
   including bulk JSON by default. Rejected: it makes every listing pay for
   a feature most projects will not use, and it changes JSON shape based on
   working-tree state.

4. **Automatic in human output, explicit in JSON.** Enrich `list`, `show`
   and `tree` silently since they are already for people, and gate only the
   machine-readable path. Rejected: it makes plain `tl list` cost one file
   read per row and makes the same ledger render differently on two
   machines. The default path stays pure.

5. **Scenario anchors in the reference** —
   `--ref "features/login.feature#Signing in"`. Buys story-level
   granularity with no change to spec files. Rejected: the link is an exact
   match against a human-written title, so renaming a scenario — ordinary
   editing, not a mistake — breaks it silently. It also requires new
   reference syntax, which `tl doctor` would treat as a dead path and
   `--fix` would delete. Addable later as a second resolver if file
   granularity proves too coarse in practice.

6. **Ledger identifiers tagged on scenarios** — `@task-kd0` above the
   scenario, the way teams tag `@JIRA-1234`. Robust to renames, standard
   Gherkin, and needs no new reference syntax. Rejected on artefact
   lifetime: the specification is the durable artefact and the ledger is
   working state, so ticket identifiers embedded in specs outlive their
   meaning and degrade the spec for later readers.

## Open questions

- Does a *missing* referenced file mean "no spec" or "broken link"? `doctor`
  calls it a warning today; an enriched listing must render it as one or the
  other, and they carry opposite meanings.
- If a scenario count is displayed, does a `Scenario Outline` count as one
  block or as one per `Examples` row? Cosmetic; count blocks and document it.
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
