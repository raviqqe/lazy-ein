Feature: Function
  Scenario Outline: Call functions
    Given a file named "main.ein" with:
    """
    f : Number -> Number
    f x = x

    g : Number -> Number -> Number
    g x y = y

    h : Number -> Number -> Number -> Number
    h x y z = z

    i : (Number -> Number) -> Number
    i f = f 42

    main : Number -> [Number]
    main x = [<expression>]
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"
    Examples:
      | expression     |
      | 42             |
      | x              |
      | let y = x in y |
      | let y = x in x |
      | f x            |
      | g 13 x         |
      | (g 13) x       |
      | h 13 13 x      |
      | ((h 13) 13) x  |
      | f (f x)        |
      | f (f (f x))    |
      | g (f x) (f x)  |
      | i f            |

  Scenario Outline: Partial applications
    Given a file named "main.ein" with:
    """
    f : Number -> [Number]
    f x = [x]

    g : Number -> Number -> [Number]
    g x y = [y]

    h : Number -> Number -> Number -> [Number]
    h x y z = [z]

    main : Number -> [Number]
    <bind>
    """
    When I successfully run `ein build main.ein`
    And I successfully run `sh -c ./a.out`
    Then the stdout from "sh -c ./a.out" should contain exactly "42"
    Examples:
      | bind                         |
      | main = f                     |
      | main = g 13                  |
      | main = h 13 13               |
