Feature: Build binaries
  Scenario Outline: Build executables
    Given a file named "main.ein" with:
    """
    f : Number -> Number
    f x = x

    g : Number -> Number -> Number
    g x y = y

    h : Number -> Number -> Number -> Number
    h x y z = z

    main : Number -> Number
    <bind>
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"
    Examples:
      | bind                         |
      | main x = 42                  |
      | main x = x                   |
      | main x = let y = x in y      |
      | main x = let y = x in x      |
      | main x = f x                 |
      | main x = g 13 x              |
      | main x = h 13 13 x           |
      | main x = f (f x)             |
      | main x = f (f (f x))         |
      | main x = g (f x) (f x)       |
      | main = f                     |
      | main = g 13                  |
      | main = h 13 13               |
      | main x = 40 + 2              |
      | main x = 21 + 7 * 3          |
      | main x = 7 + 12 / 3 * 10 - 5 |
      | main x = f 40 + 2            |
      | main x = case 1 of 1 -> 42   |

  Scenario: Use default alternatives in case expressions
    Given a file named "main.ein" with:
    """
    main : Number -> Number
    main x =
      case 1 of
        2 -> 13
        x -> 41 + x
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"

  Scenario: Use nested case expressions
    Given a file named "main.ein" with:
    """
    main : Number -> Number
    main x =
      case 1 of
        2 -> 13
        x -> case 2 of
               3 -> 13
               x -> 40 + x
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"

  Scenario: Emit build errors
    Given a file named "main.ein" with:
    """
    main : Number -> Number
    main x =
    """
    When I run `ein build main.ein`
    Then the exit status should not be 0
    And the stderr should contain exactly:
    """
    SyntaxError: unexpected end of source
    main.ein:2:9:	main x =
    """

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
