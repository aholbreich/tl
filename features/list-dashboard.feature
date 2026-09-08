@implemented
Feature: Task dashboard
  As a developer
  I want a regeneratable Markdown overview of active work
  So that a wiki page can show task context without opening each task

  Background:
    Given an initialized task ledger repository

  Scenario: The dashboard groups compact task blocks by status
    Given the following tasks exist:
      | id       | status      | priority | claimed by | title              |
      | task-api | open        | high     |            | Document API       |
      | task-cli | in_progress | medium   | codex      | Improve CLI output |
      | task-old | done        | low      |            | Retire old client  |
    And "task-api" has a description "Describe the public endpoints."
    When the developer runs `tl list --dashboard`
    Then the dashboard Markdown is exactly:
      """
      # Task dashboard

      ## In progress

      ### task-cli: Improve CLI output

      - Status: in\_progress
      - Priority: medium
      - Type: task
      - Claimant: codex

      ## Open

      ### task-api: Document API

      - Status: open
      - Priority: high
      - Type: task
      - Claimant: -

      Description: Describe the public endpoints.
      """

  Scenario: The dashboard includes task references
    Given a task "task-api" with references "features/list.feature" and "https://example.org/api"
    When the developer runs `tl list --dashboard`
    Then the dashboard contains the line "- features/list.feature"
    And the dashboard contains the line "- https://example.org/api"

  Scenario: Dashboard filters compose to select a scoped overview
    Given the following tasks exist:
      | id       | status | priority | type    | tags | title            |
      | task-api | open   | high     | feature | api  | Document API     |
      | task-cli | open   | low      | feature | api  | Improve CLI      |
      | task-fix | open   | high     | bug     | api  | Fix API response |
      | task-doc | open   | high     | feature | docs | Publish guide    |
    When the developer runs `tl list --dashboard --type feature --priority h --tag api --status open`
    Then the output lists "task-api"
    And the output does not list "task-cli"
    And the output does not list "task-fix"
    And the output does not list "task-doc"

  Scenario: Closed tasks can be included explicitly
    Given the following tasks exist:
      | id       | status    | title             |
      | task-api | done      | Document API      |
      | task-cli | cancelled | Retire old client |
    When the developer runs `tl list --dashboard --all`
    Then the dashboard contains the line "## Done"
    And the dashboard contains the line "## Cancelled"
    And the output lists "task-api"
    And the output lists "task-cli"

  Scenario: A closed status filter overrides the default active view
    Given the following tasks exist:
      | id       | status | title        |
      | task-api | done   | Document API |
      | task-cli | open   | Improve CLI  |
    When the developer runs `tl list --dashboard --status done`
    Then the output lists "task-api"
    And the output does not list "task-cli"

  Scenario: The dashboard can show only the resolved actor's work
    Given environment variable "TL_ACTOR" is "codex"
    And the following tasks exist:
      | id       | status      | claimed by | title        |
      | task-api | in_progress | codex      | Document API |
      | task-cli | in_progress | pi         | Improve CLI  |
    When the developer runs `tl list --dashboard --mine`
    Then the output lists "task-api"
    And the output does not list "task-cli"

  Scenario: JSON takes precedence over Markdown
    Given the following tasks exist:
      | id       | status | title        |
      | task-api | open   | Document API |
      | task-cli | done   | Improve CLI  |
    When the developer runs `tl list --dashboard --json`
    Then the JSON output is an array of 1 tasks
    And the JSON output contains a task with identifier "task-api"

  Scenario: An empty dashboard explains that no tasks match
    Given no tasks exist
    When the developer runs `tl list --dashboard`
    Then the dashboard Markdown is exactly:
      """
      # Task dashboard

      No tasks match.
      """

  Scenario: Forced terminal color does not contaminate Markdown
    Given the following tasks exist:
      | id       | priority | title        |
      | task-api | high     | Document API |
    When the developer runs `tl --color=always list --dashboard`
    Then the command exits with code 0
    And the output does not contain ANSI color

  Scenario: The dashboard can be scoped to another actor
    Given the following tasks exist:
      | id       | status      | claimed by | title        |
      | task-api | in_progress | codex      | Document API |
      | task-cli | in_progress | pi         | Improve CLI  |
    When the developer runs `tl list --dashboard --claimed-by pi`
    Then the output lists "task-cli"
    And the output does not list "task-api"

  Scenario: Tasks without a type belong to the default task type
    Given the following tasks exist:
      | id       | type | title        |
      | task-api |      | Document API |
      | task-cli | bug  | Improve CLI  |
    When the developer runs `tl list --dashboard --type task`
    Then the output lists "task-api"
    And the output does not list "task-cli"

  Scenario: Invalid priority filters are rejected
    When the developer runs `tl list --dashboard --priority urgent`
    Then the command exits with code 2
