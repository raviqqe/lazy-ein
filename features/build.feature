Feature: Build binaries
  Scenario: Build an executable
    Given a file named "main.ein" with:
    """
    main : Number -> Number
    main x = 42
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"

  Scenario: Build an executable of an identity function
    Given a file named "main.ein" with:
    """
    main : Number -> Number
    main x = x
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"
