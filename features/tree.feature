@implemented
Feature: Rendering the dependency graph
  As a developer or product owner
  I want to see how a task decomposes into the tasks it depends on
  So that I can tell which slice is blocking without reading each task in turn

  # A task's children are the tasks it depends on, so a parent waits for the
  # slices beneath it. A root is a task nothing else depends on.

  Background:
    Given an initialized task ledger repository

  # -------------------------------------------------------------------------
  # Shape — children are dependencies, roots are what nothing depends on.
  # -------------------------------------------------------------------------
  Scenario: A task is rendered above the tasks it depends on
    Given a task "task-par001" titled "Add RSS sources"
    And a task "task-chi001" titled "Subscribe to a URL"
    And "task-par001" depends on "task-chi001"
    When the developer runs `tl tree task-par001`
    Then the output shows "task-chi001" as a child of "task-par001"

  Scenario: Depth is rendered for a chain of dependencies
    Given a task "task-par001" titled "Add RSS sources"
    And a task "task-chi001" titled "Subscribe to a URL"
    And a task "task-gra001" titled "Fetch and parse feeds"
    And "task-par001" depends on "task-chi001"
    And "task-chi001" depends on "task-gra001"
    When the developer runs `tl tree task-par001`
    Then the output shows "task-gra001" at depth 2

  Scenario: Without an argument every root is rendered
    Given a task "task-par001" titled "Add RSS sources"
    And a task "task-chi001" titled "Subscribe to a URL"
    And a task "task-oth001" titled "Unrelated work"
    And "task-par001" depends on "task-chi001"
    When the developer runs `tl tree`
    Then the output shows "task-par001" at depth 0
    And the output shows "task-oth001" at depth 0
    And the output shows "task-chi001" at depth 1

  Scenario: A task nothing depends on and that depends on nothing is its own root
    Given a task "task-lon001" titled "Standalone work"
    When the developer runs `tl tree`
    Then the output shows "task-lon001" at depth 0

  Scenario: A shared dependency is rendered under each parent
    Given a task "task-aaa001" titled "First parent"
    And a task "task-bbb001" titled "Second parent"
    And a task "task-sha001" titled "Shared slice"
    And "task-aaa001" depends on "task-sha001"
    And "task-bbb001" depends on "task-sha001"
    When the developer runs `tl tree`
    Then the output shows "task-sha001" 2 times

  # -------------------------------------------------------------------------
  # Status and titles travel with each node.
  # -------------------------------------------------------------------------
  Scenario: Each node shows its identifier, title and status
    Given a task "task-par001" titled "Add RSS sources"
    When the developer runs `tl tree`
    Then the output contains "task-par001"
    And the output contains "Add RSS sources"
    And the output contains "open"

  # -------------------------------------------------------------------------
  # Priority and colour follow the tl list conventions.
  # -------------------------------------------------------------------------
  Scenario: Each node shows its priority
    Given the following tasks exist:
      | id          | status | priority | title              |
      | task-abc123 | open   | high     | High priority task |
    When the developer runs `tl tree`
    Then the tree row for "task-abc123" contains "high"

  Scenario: Rendering with forced colour highlights priority values
    Given the following tasks exist:
      | id          | status | priority | title                |
      | task-abc123 | open   | high     | High priority task   |
      | task-def456 | open   | medium   | Medium priority task |
      | task-ghi789 | open   | low      | Low priority task    |
    When the developer runs `tl --color=always tree`
    Then the output colorizes "high" with "red"
    And the output colorizes "medium" with "yellow"
    And the output colorizes "low" with "blue"

  Scenario: Closed rows are dimmed when revealed with --all
    Given a task "task-par001" titled "Add RSS sources"
    And a task "task-don001" with status "done"
    And "task-par001" depends on "task-don001"
    When the developer runs `tl --color=always tree --all`
    Then the output colorizes the line for "task-don001" with "dim"

  # -------------------------------------------------------------------------
  # Spec status — the same flag list and ready carry.
  # -------------------------------------------------------------------------
  Scenario: The tree reports the spec of a referenced feature file
    Given a task "task-abc123" with reference "features/login.feature"
    And the repository has a feature file "features/login.feature" with 3 scenarios
    When the developer runs `tl tree --spec-status`
    Then the tree row for "task-abc123" contains "features/login.feature (3)"

  Scenario: A dependency carries its own spec in the tree
    Given a task "task-par001" titled "Add RSS sources"
    And a task "task-chi001" with reference "features/login.feature"
    And "task-par001" depends on "task-chi001"
    And the repository has a feature file "features/login.feature" with 2 scenarios
    When the developer runs `tl tree task-par001 --spec-status`
    Then the tree row for "task-chi001" contains "features/login.feature (2)"

  Scenario: Plain tree does not resolve specs
    Given a task "task-abc123" with reference "features/login.feature"
    And the repository has a feature file "features/login.feature" with 3 scenarios
    When the developer runs `tl tree`
    Then the output does not contain "features/login.feature"

  # -------------------------------------------------------------------------
  # Closed tasks — hidden by default, like tl list.
  # -------------------------------------------------------------------------
  Scenario: A closed dependency is hidden by default
    Given a task "task-par001" titled "Add RSS sources"
    And a task "task-don001" with status "done"
    And "task-par001" depends on "task-don001"
    When the developer runs `tl tree task-par001`
    Then the output does not contain "task-don001"

  Scenario: --all reveals closed dependencies
    Given a task "task-par001" titled "Add RSS sources"
    And a task "task-don001" with status "done"
    And "task-par001" depends on "task-don001"
    When the developer runs `tl tree task-par001 --all`
    Then the output shows "task-don001" as a child of "task-par001"

  # -------------------------------------------------------------------------
  # Cycles must render rather than hang. tl doctor owns diagnosing them.
  # -------------------------------------------------------------------------
  Scenario: A dependency cycle renders once and is marked
    Given a task "task-cyc001" titled "First of a cycle"
    And a task "task-cyc002" titled "Second of a cycle"
    And "task-cyc001" depends on "task-cyc002"
    And "task-cyc002" depends on "task-cyc001"
    When the developer runs `tl tree task-cyc001`
    Then the command exits with code 0
    And the output marks "task-cyc001" as a cycle

  # -------------------------------------------------------------------------
  # Errors.
  # -------------------------------------------------------------------------
  Scenario: Rendering an unknown task reports it as not found
    When the developer runs `tl tree task-nope99`
    Then the command exits with code 3

  # -------------------------------------------------------------------------
  # JSON — nested children so consumers do not re-derive the graph.
  # -------------------------------------------------------------------------
  Scenario: JSON nests each task's dependencies as children
    Given a task "task-par001" titled "Add RSS sources"
    And a task "task-chi001" titled "Subscribe to a URL"
    And "task-par001" depends on "task-chi001"
    When the developer runs `tl tree task-par001 --json`
    Then the JSON tree root "task-par001" has a child "task-chi001"

  Scenario: JSON emits an empty children array for a leaf
    Given a task "task-lon001" titled "Standalone work"
    When the developer runs `tl tree task-lon001 --json`
    Then the JSON tree root "task-lon001" has no children
