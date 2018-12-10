Feature: Build binaries
  Scenario: Build an executable
    Given a file named "main.jsonxx" with:
    """
    main : Number -> Number
    main x = 42
    """
    When I successfully run `jsonxx build main.jsonxx`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"
