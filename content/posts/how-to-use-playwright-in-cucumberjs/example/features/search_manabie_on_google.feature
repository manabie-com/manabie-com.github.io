Feature: Google search

    Scenario: Bob try to find Manabie on Google
        Given Bob go to Google website
        When Bob search Manabie
        Then Manabie appears on result list