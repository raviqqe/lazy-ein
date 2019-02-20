Feature: Build
  Scenario: Build exectuables
    Given a file named "main.ein" with:
    """
    main : Number -> Number
    main x = 42
    """
    And I successfully run `ein build main.ein`
    When I run `ls a.out`
    Then the exit status should be 0

  Scenario: Build modules
    Given a file named "main.ein" with:
    """
    export { x }

    x : Number
    x = case 1 of 2 -> 3
    """
    And I successfully run `ein build main.ein`
    When I run `ls a.out`
    Then the exit status should not be 0

  Scenario: Build main modules importing submodules
    Given a file named "foo.ein" with:
    """
    export { x }

    x : Number
    x = 42
    """
    Given a file named "main.ein" with:
    """
    import "tmp/aruba/foo"

    main : Number -> Number
    main x = foo.x
    """
    And I successfully run `ein build main.ein`
    When I run `ls a.out`
    Then the exit status should be 0
