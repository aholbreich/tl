Feature: Inspect the event history for a task
  As a reviewer or auditor
  I want to see every event recorded for a task
  So that I can reconstruct who did what and when

  Background:
    Given an initialized TaskLedger repository

  Scenario: History accumulates events across mutations in chronological order
    Given a ready task "task-abc123"
    When the agent runs `tl claim task-abc123 --actor claude-code:main`
    And the agent runs `tl note task-abc123 --actor claude-code:main --message "WIP"`
    And the agent runs `tl close task-abc123 --actor claude-code:main`
    And the developer runs `tl history task-abc123`
    Then the history output lists events in this order:
      | event       |
      | created     |
      | claimed     |
      | note_added  |
      | closed      |

  Scenario: Each history entry records an actor and a timestamp
    Given a task "task-abc123" claimed by "claude-code:main" with an active lease
    When the developer runs `tl history task-abc123`
    Then the history output for the "claimed" event shows actor "claude-code:main"
    And the history output for the "claimed" event has a timestamp

  Scenario: History as JSON returns the raw event objects for the task
    Given a ready task "task-abc123"
    When the developer runs `tl history task-abc123 --json`
    Then the JSON output is an array of event objects for "task-abc123"
    And each event object contains a type, a timestamp, and an actor

  Scenario: History for an unknown task is rejected
    When the developer runs `tl history task-xyz`
    Then the command exits with code 3
