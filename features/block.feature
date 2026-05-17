@implemented
Feature: Mark a task blocked by an external condition
  As an actor unable to make progress
  I want to mark a task blocked with a reason
  So that other actors know not to pick it up until the blocker is cleared

  Background:
    Given an initialized TaskLedger repository

  Scenario: An open task can be marked blocked with a reason
    Given a task "task-abc123" with status "open"
    When the developer runs `tl block task-abc123 -m "Waiting on upstream library release."`
    Then "task-abc123" has status "blocked"
    And "task-abc123" has a note containing "Waiting on upstream library release."
    And an event "blocked" is recorded for "task-abc123"

  Scenario: Blocking a claimed task releases the claim
    Given a task "task-abc123" claimed by "claude-code:main" with an active lease
    When the agent runs `tl block task-abc123 --actor claude-code:main -m "CI infra down."`
    Then "task-abc123" has status "blocked"
    And "task-abc123" is not claimed

  Scenario: A blocked task is not in the ready queue
    Given a task "task-abc123" with status "blocked"
    When the agent runs `tl ready`
    Then the ready output does not contain "task-abc123"

  Scenario: Unblocking returns the task to open
    Given a task "task-abc123" with status "blocked"
    When the developer runs `tl unblock task-abc123`
    Then "task-abc123" has status "open"
    And an event "unblocked" is recorded for "task-abc123"

  Scenario: After unblocking the task is ready again
    Given a task "task-abc123" with status "blocked" and no dependencies
    When the developer runs `tl unblock task-abc123`
    And the agent runs `tl ready`
    Then the ready output contains "task-abc123"

  Scenario: Blocking without a reason is rejected
    Given a task "task-abc123" with status "open"
    When the developer runs `tl block task-abc123`
    Then the command exits with code 2
    And the output reports that a reason is required

  Scenario: Unblocking a task that is not blocked is rejected
    Given a task "task-abc123" with status "open"
    When the developer runs `tl unblock task-abc123`
    Then the command reports the task is not blocked
    And "task-abc123" still has status "open"
