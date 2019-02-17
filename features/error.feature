Feature: Error
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
    tmp/aruba/main.ein:2:9:	main x =
    """

  Scenario: Emit build errors
    Given a file named "main.ein" with:
    """
    main : Number -> Number
    main x = case 1 of 2 -> 3
    """
    And I run `ein build main.ein`
    When I run `sh -c ./a.out`
    Then the exit status should not be 0
    And the stderr from "sh -c ./a.out" should contain exactly "Match error!"
