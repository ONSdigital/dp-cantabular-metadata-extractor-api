Feature: Health check

  Scenario: Check the health endpoint
    When I GET "/health"
    Then the HTTP status code should be "200"