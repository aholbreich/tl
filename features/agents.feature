@implemented
Feature: Print recommended agent instructions
  As a developer setting up an agent-friendly repository
  I want to see the recommended AGENTS.md snippet
  So that I can copy it into my own AGENTS.md without the tool editing files for me

  Background:
    Given an initialized TaskLedger repository

  Scenario: Running agents prints the recommended AGENTS.md snippet to stdout
    When the developer runs `tl agents`
    Then the output contains a "TaskLedger Workflow" heading
    And the output describes the ready, claim, show, note, and close steps
    And the output formats task commands as Markdown code spans

  Scenario: Running agents does not modify any existing AGENTS.md
    Given the file "AGENTS.md" exists with content "# My Project"
    When the developer runs `tl agents`
    Then the file "AGENTS.md" still has content "# My Project"
