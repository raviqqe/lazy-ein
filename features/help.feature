Feature: Help
  Scenario: Run a command without any subcommand
    When I run `jsonxx`
    Then the exit status should not be 0
    And the stderr should contain "Usage"

  Scenario: Run a build subcommand without any filename
    When I run `jsonxx build`
    Then the exit status should not be 0
    And the stderr should contain "Usage"
