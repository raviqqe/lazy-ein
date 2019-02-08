Feature: List
  Scenario: Define list variables
    Given a file named "main.ein" with:
    """
    l : [Number]
    l = [42]

    main : Number -> Number
    main x = 42
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"

  Scenario: Use list case expressions
    Given a file named "main.ein" with:
    """
    l : [Number]
    l = [42]

    main : Number -> Number
    main x = case l of [y] -> y
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"
