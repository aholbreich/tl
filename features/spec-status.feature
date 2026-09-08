@implemented
Feature: Spec status for referenced specifications
  As a developer or product owner
  I want to see which tasks point at a specification and whether it exists
  So that I can tell specified work from unspecified work without opening files

  # A reference whose path ends in .feature is a spec reference. Recognition is
  # a string rule and costs nothing. Reading the file happens only under
  # --spec-status — see .decisions/0002-reading-referenced-files.md.

  Background:
    Given an initialized task ledger repository

  # -------------------------------------------------------------------------
  # Recognition — string rule, no file access.
  # -------------------------------------------------------------------------
  Scenario: A reference ending in .feature is marked as a spec in tl show
    Given a task "task-abc123" with references "features/login.feature" and "src/auth/login.go"
    When the developer runs `tl show task-abc123`
    Then the output contains "features/login.feature (spec)"
    And the output contains "  - src/auth/login.go"

  Scenario: A URL ending in .feature is not treated as a spec reference
    Given a task "task-abc123" with reference "https://example.com/specs/login.feature"
    When the developer runs `tl show task-abc123`
    Then the output does not contain "(spec)"

  # -------------------------------------------------------------------------
  # --spec-status — the only path that opens a referenced file.
  # -------------------------------------------------------------------------
  Scenario: The spec column reports an existing specification and its scenario count
    Given a task "task-abc123" with reference "features/login.feature"
    And the repository has a feature file "features/login.feature" with 3 scenarios
    When the developer runs `tl list --spec-status`
    Then the output lists "task-abc123" with spec "features/login.feature" and 3 scenarios

  Scenario: The spec column reports a referenced specification that does not exist
    Given a task "task-abc123" with reference "features/missing.feature"
    When the developer runs `tl list --spec-status`
    Then the output lists "task-abc123" with spec "features/missing.feature" marked missing

  Scenario: A task with no spec reference shows no spec
    Given a task "task-abc123" with reference "src/auth/login.go"
    When the developer runs `tl list --spec-status`
    Then the output lists "task-abc123" with no spec

  Scenario: A Scenario Outline counts as one scenario
    Given a task "task-abc123" with reference "features/login.feature"
    And the repository has a feature file "features/login.feature" with 1 scenario and 1 scenario outline of 3 examples
    When the developer runs `tl list --spec-status`
    Then the output lists "task-abc123" with spec "features/login.feature" and 2 scenarios

  Scenario: A task carrying two spec references reports both
    Given a task "task-abc123" with references "features/login.feature" and "features/logout.feature"
    And the repository has a feature file "features/login.feature" with 2 scenarios
    And the repository has a feature file "features/logout.feature" with 1 scenarios
    When the developer runs `tl list --spec-status`
    Then the output lists "task-abc123" with 2 specs

  Scenario: An unreadable specification degrades to unknown rather than failing the listing
    Given a task "task-abc123" with reference "features/login.feature"
    And the repository has an unreadable file "features/login.feature"
    When the developer runs `tl list --spec-status`
    Then the command exits with code 0
    And the output lists "task-abc123" with spec "features/login.feature" marked unknown

  # -------------------------------------------------------------------------
  # Default path stays pure — no reads without the flag.
  # -------------------------------------------------------------------------
  Scenario: Plain list does not show a spec column
    Given a task "task-abc123" with reference "features/login.feature"
    And the repository has a feature file "features/login.feature" with 3 scenarios
    When the developer runs `tl list`
    Then the output does not contain "Spec"

  # -------------------------------------------------------------------------
  # JSON — schema-stable: the key is always present.
  # -------------------------------------------------------------------------
  Scenario: JSON emits a null spec when the flag is not passed
    Given a task "task-abc123" with reference "features/login.feature"
    And the repository has a feature file "features/login.feature" with 3 scenarios
    When the developer runs `tl list --json`
    Then the JSON task "task-abc123" has a null "spec" field

  Scenario: JSON exposes the resolved spec when the flag is passed
    Given a task "task-abc123" with reference "features/login.feature"
    And the repository has a feature file "features/login.feature" with 3 scenarios
    When the developer runs `tl list --spec-status --json`
    Then the JSON task "task-abc123" has a spec entry for "features/login.feature" that exists with 3 scenarios

  Scenario: JSON emits an empty spec array for a task with no spec reference
    Given a task "task-abc123" with reference "src/auth/login.go"
    When the developer runs `tl list --spec-status --json`
    Then the JSON task "task-abc123" has an empty "spec" array

  Scenario: Ready supports the same flag
    Given a task "task-abc123" with reference "features/login.feature"
    And the repository has a feature file "features/login.feature" with 3 scenarios
    When the agent runs `tl ready --spec-status`
    Then the output lists "task-abc123" with spec "features/login.feature" and 3 scenarios
