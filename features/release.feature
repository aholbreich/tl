@implemented
Feature: Release a claim on a task
  As an actor stepping away from a task
  I want to clear my claim so others can pick it up
  So that work does not stall behind an idle claim

  Background:
    Given an initialized TaskLedger repository

  Scenario: The claim's actor releases its own claim and the task returns to open
    Given a task "task-abc123" claimed by "claude-code:main" with an active lease
    When the agent runs `tl release task-abc123 --actor claude-code:main`
    Then "task-abc123" is not claimed
    And "task-abc123" has status "open"
    And an event "released" is recorded for "task-abc123"

  Scenario: Releasing a claim held by a different actor is rejected without force
    Given a task "task-abc123" claimed by "claude-code:frontend" with an active lease
    When the agent runs `tl release task-abc123 --actor claude-code:main`
    Then the command reports the claim is held by a different actor
    And "task-abc123" is still claimed by "claude-code:frontend"

  Scenario: Forcing release of a stale claim succeeds
    Given a task "task-abc123" with an expired claim by "claude-code:frontend"
    When the developer runs `tl release task-abc123 --actor human --force`
    Then "task-abc123" is not claimed
    And "task-abc123" has status "open"
