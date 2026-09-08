@implemented
Feature: Resolve a pending task with an answer
  As a human responding to an agent's question
  I want to provide an answer that returns the task to the ready queue
  So that the agent can resume work with the needed input

  Background:
    Given an initialized task ledger repository

  Scenario: Resolving a pending task appends the answer and returns the task to open
    Given a task "task-abc123" with status "pending_human"
    And "task-abc123" has the question "Which auth provider first?"
    When the developer runs `tl resolve task-abc123 --answer "Use GitHub OAuth first."`
    Then "task-abc123" has status "open"
    And "task-abc123" has a note containing "Use GitHub OAuth first."
    And an event "pending_resolved" is recorded for "task-abc123"

  Scenario: An explicit actor receives credit for the answer instead of the environment actor
    Given a task "task-auth" with status "pending_human"
    And "task-auth" has the question "Which auth provider first?"
    And environment variable "TL_ACTOR" is "env-agent"
    When the developer runs `tl resolve task-auth --answer "Use GitHub OAuth first." --actor aho`
    Then an event "pending_resolved" is recorded for "task-auth" by "aho"
    And "task-auth" has a canonical "resolved" note from "aho" with message "Use GitHub OAuth first."

  Scenario: Resolving a task that is not pending_human is rejected
    Given a task "task-abc123" with status "open"
    When the developer runs `tl resolve task-abc123 --answer "moot"`
    Then the command reports the task is not pending_human
    And "task-abc123" still has status "open"
