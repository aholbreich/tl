@implemented
Feature: Actor identity resolution for claims
  As an agent picking up work
  I want my identity to be resolved without needing --actor every time
  So that claims are frictionless while still preventing collisions

  Background:
    Given an initialized task ledger repository

  Scenario: CLI --actor takes highest priority
    Given environment variable "TL_ACTOR" is "env-agent"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123 --actor cli-agent`
    Then "task-abc123" is claimed by "cli-agent"

  Scenario: TL_ACTOR env var when --actor is absent
    Given environment variable "TL_ACTOR" is "pi:main"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123`
    Then "task-abc123" is claimed by "pi:main"

  Scenario: ACTOR_NAME as second env fallback
    Given environment variable "ACTOR_NAME" is "fallback-agent"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123`
    Then "task-abc123" is claimed by "fallback-agent"

  Scenario: BEADS_ACTOR as third env fallback
    Given environment variable "BEADS_ACTOR" is "beads-agent"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123`
    Then "task-abc123" is claimed by "beads-agent"

  Scenario: TL_ACTOR takes precedence over ACTOR_NAME and BEADS_ACTOR
    Given environment variable "TL_ACTOR" is "primary"
    And environment variable "ACTOR_NAME" is "secondary"
    And environment variable "BEADS_ACTOR" is "tertiary"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123`
    Then "task-abc123" is claimed by "primary"

  Scenario: Auto-detect agent name when no explicit actor is set
    Given the detected agent is "claude"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123`
    Then "task-abc123" is claimed by "claude"

  # -------------------------------------------------------------------------
  # Auto-detection from harness environment markers. Only markers a harness
  # sets reliably are trusted: a false positive attributes one agent's claims
  # to another.
  # -------------------------------------------------------------------------
  Scenario: Aider is detected from AIDER_MODEL
    Given environment variable "AIDER_MODEL" is "gpt-4o"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123`
    Then "task-abc123" is claimed by "aider"

  Scenario: Windsurf is detected from CODEIUM_API_KEY
    Given environment variable "CODEIUM_API_KEY" is "sk-example"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123`
    Then "task-abc123" is claimed by "windsurf"

  Scenario: Pi is detected from PI_CODING_AGENT
    Given environment variable "PI_CODING_AGENT" is "true"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123`
    Then "task-abc123" is claimed by "pi"

  Scenario: Pi is detected from PI_AGENT_ID
    Given environment variable "PI_AGENT_ID" is "agent-7"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123`
    Then "task-abc123" is claimed by "pi"

  Scenario: PI_CODING_AGENT set to false is not a pi marker
    Given environment variable "PI_CODING_AGENT" is "false"
    And environment variable "AIDER_MODEL" is "gpt-4o"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123`
    Then "task-abc123" is claimed by "aider"

  Scenario: Claude Code takes precedence over a later marker
    Given environment variable "CLAUDE_CODE_SESSION_ID" is "0f9c1a2b"
    And environment variable "AIDER_MODEL" is "gpt-4o"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123`
    Then "task-abc123" is claimed by "claude-code"

  Scenario: An explicit actor overrides a detected harness
    Given environment variable "AIDER_MODEL" is "gpt-4o"
    And a ready task "task-abc123"
    When the agent runs `tl claim task-abc123 --actor human`
    Then "task-abc123" is claimed by "human"

  Scenario: Reject claim from a different actor
    Given a task "task-abc123" claimed by "claude-code:frontend" with an active lease
    And environment variable "TL_ACTOR" is "claude-code:main"
    When the agent runs `tl claim task-abc123`
    Then the command exits with code 5
    And "task-abc123" is still claimed by "claude-code:frontend"

  Scenario: Same actor can renew its own claim
    Given a task "task-abc123" claimed by "claude-code:main" with an active lease
    And environment variable "TL_ACTOR" is "claude-code:main"
    When the agent runs `tl claim task-abc123`
    Then the claim expiry for "task-abc123" is extended
