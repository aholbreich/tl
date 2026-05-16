@implemented
Feature: Require an initialized ledger
  As a developer or agent
  I want commands that need ledger state to explain missing setup
  So that I can initialize TaskLedger without reading a raw path error

  Scenario Outline: Ledger commands report how to initialize an uninitialized repository
    Given the current directory has no TaskLedger ledger
    When the developer runs `tl <command>`
    Then the command exits with code 1
    And the output reports that TaskLedger is not initialized
    And the output suggests running "tl init"

    Examples:
      | command                               |
      | list                                  |
      | show task-abc123                      |
      | create "Add login validation"         |
