# go-pagerduty is a client for the PagerDuty Event API (v2)

This package provide a simple client for the PagerDuty Event API.

There is only one method: `SendEvent` this method will send an event to the API.

## How to use it:

Example:

```golang
  apiKey := "MyApiKey" // Your api key (this key is displayed when you add a new EventAPI integration)

  eventAction := pagerduty.EventActionTrigger // trigger a new event
  // This can take the follwing values:
  // * pagerduty.EventActionTrigger
  // * pagerduty.EventActionAcknowledge
  // * pagerduty.EventActionResolve

  source := "My Awesome Service" // The name of the service

  severity := pagerduty.ServerityError // The event severity
  // This can take the follwing values:
  // * pagerduty.ServerityCritical
  // * pagerduty.ServerityError
  // * pagerduty.ServerityWarning
  // * pagerduty.ServerityInfo

  summary := "The servers were replaced by cats"

  details := "Cats are taking over our IT infrastructure"

  component := "Infrastructure"

  group := "hardware-issues"

  resp err := pagerduty.SendEvent(key, eventAction, source, severity, summary, pagerduty.EventOptions{
    Source: source,
    Component: component,
    Group: group,
  })
```

More informations are availabe in [the official documentation](https://v2.developer.pagerduty.com/docs/send-an-event-events-api-v2).
