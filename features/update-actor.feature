@implemented
Feature: Actor attribution for task updates
  As a developer or agent
  I want task updates to use the standard actor identity rules
  So that the audit trail consistently identifies who changed a task

  Background:
    Given an initialized task ledger repository

  Scenario Outline: Updates use the identity fallback chain when no actor flag is supplied
    Given a task "task-auth" with status "pending_human"
    And environment variable "TL_ACTOR" is "<tl-actor>"
    And environment variable "ACTOR_NAME" is "<actor-name>"
    And environment variable "BEADS_ACTOR" is "<beads-actor>"
    And the detected agent is "detected-agent"
    When the developer runs `tl <command> task-auth <change>`
    Then an event "<event>" is recorded for "task-auth" by "<actor>"

    Examples:
      | command | change          | tl-actor | actor-name | beads-actor | event            | actor          |
      | resolve | --answer GitHub | primary  | secondary  | tertiary    | pending_resolved | primary        |
      | resolve | --answer GitHub |          | secondary  | tertiary    | pending_resolved | secondary      |
      | resolve | --answer GitHub |          |            | tertiary    | pending_resolved | tertiary       |
      | resolve | --answer GitHub |          |            |             | pending_resolved | detected-agent |
      | refine  | --priority high | primary  | secondary  | tertiary    | refined          | primary        |
      | refine  | --priority high |          | secondary  | tertiary    | refined          | secondary      |
      | refine  | --priority high |          |            | tertiary    | refined          | tertiary       |
      | refine  | --priority high |          |            |             | refined          | detected-agent |

  Scenario: Mixed field and reference changes attribute every event to the explicit actor
    Given a task "task-auth" with reference "old.md"
    And environment variable "TL_ACTOR" is "env-agent"
    When the developer runs `tl refine task-auth --priority high --add-ref auth.md --remove-ref old.md --actor aho`
    Then an event "refined" is recorded for "task-auth" by "aho"
    And an event "reference_added" is recorded for "task-auth" by "aho"
    And an event "reference_removed" is recorded for "task-auth" by "aho"

  Scenario: A reference-only update uses the environment actor
    Given a task "task-auth" with no references
    And environment variable "TL_ACTOR" is "env-agent"
    When the developer runs `tl refine task-auth --add-ref docs/auth.md`
    Then an event "reference_added" is recorded for "task-auth" by "env-agent"
    And no event "refined" is recorded for "task-auth"

  Scenario: An actor flag alone does not count as a field change
    Given a task "task-auth" exists
    When the developer runs `tl refine task-auth --actor aho`
    Then the command exits with code 2
    And no event "refined" is recorded for "task-auth"

  Scenario: An actor flag does not turn an unchanged reference into an update
    Given a task "task-auth" with reference "docs/auth.md"
    When the developer runs `tl refine task-auth --add-ref docs/auth.md --actor aho`
    Then the command exits with code 0
    And no event "reference_added" is recorded for "task-auth"
    And no event "refined" is recorded for "task-auth"
