Feature: Number
  Scenario Outline: Use operators
    Given a file named "main.ein" with:
    """
    f : Number -> Number
    f x = x

    main : Number -> [Number]
    main x = [<expression>]
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"
    Examples:
      | expression          |
      | 42                  |
      | 40 + 2              |
      | 21 + 7 * 3          |
      | 7 + 12 / 3 * 10 - 5 |
      | f 40 + 2            |

  Scenario: Use case expressions
    Given a file named "main.ein" with:
    """
    main : Number -> [Number]
    main x = case 1 of 1 -> [42]
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"

  Scenario: Use default alternatives in case expressions
    Given a file named "main.ein" with:
    """
    main : Number -> [Number]
    main x =
      case 1 of
        2 -> [13]
        x -> [41 + x]
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"

  Scenario: Use nested case expressions
    Given a file named "main.ein" with:
    """
    main : Number -> [Number]
    main x =
      case 1 of
        2 -> [13]
        x -> case 2 of
               3 -> [13]
               x -> [40 + x]
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"

  Scenario: Sum numbers up to 100
    Given a file named "main.ein" with:
    """
    sum : Number -> Number
    sum x =
      case x of
        0 -> 0
        x -> x + sum (x - 1)

    main : Number -> [Number]
    main x = [sum 100]
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "5050"
