// Package pagerduty provides a client for the PagerDuty EventsAPI v2
package pagerduty

var endpoint = "https://events.pagerduty.com/v2/enqueue"

// EventAction define the type of event (can be Acknowledge, Trigger, Resolve)
type EventAction string

// Severity is the perceived severity of the status the event is describing with respect to the affected system (can be Critical, Errorn Warning, Info)
type Severity string

const (
	// EventActionTrigger will trigger a new event
	EventActionTrigger EventAction = "trigger"
	// EventActionAcknowledge will acknowledge the current event
	EventActionAcknowledge EventAction = "acknowledge"
	// EventActionResolve will resolve the current event
	EventActionResolve EventAction = "resolve"

	// SeverityCritical will set the event serverity to critical
	SeverityCritical Severity = "critical"
	// SeverityError will set the event severity to Error
	SeverityError Severity = "error"
	// SeverityWarning will set the event severity to Warning
	SeverityWarning Severity = "warning"
	// SeverityInfo will set the event serverity to info
	SeverityInfo Severity = "info"
)

type event struct {
	RoutingKey  string   `json:"routing_key"`
	EventAction string   `json:"event_action"`
	DeDupKey    string   `json:"dedup_key,omitempty"`
	Payload     payload  `json:"payload"`
	Images      []*Image `json:"images,omitempty"`
	Links       []*Link  `json:"links,omitempty"`
}

type payload struct {
	Source        string      `json:"source"`
	Summary       string      `json:"summary"`
	Severity      string      `json:"severity"`
	Component     string      `json:"component,omitempty"`
	Group         string      `json:"group,omitempty"`
	Class         string      `json:"class,omitempty"`
	CustomDetails interface{} `json:"custom_details,omitempty"`
}

// Image is a structure that represent an image linked to the event
type Image struct {
	Src  string `json:"src"`  // Src represent the source of the image being attached to the incident. This image must be served via HTTPS.
	Href string `json:"href"` // Href: Optional URL; makes the image a clickable link.
	Alt  string `json:"alt"`  // Alt: Optional alternative text for the image.
}

// Link is a structure that represent a link linked to the event
type Link struct {
	Href string `json:"href"` // URL of the link to be attached.
	Text string `json:"test"` // Plain text that describes the purpose of the link, and can be used as the link's text.
}

// Response represent an api response
type Response struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	DeDupKey string `json:"dedup_key"`
}

// EventOptions is the structure used to pass optionnal options
type EventOptions struct {
	DeDupKey      string      // DeDupKey: Deduplication key for correlating triggers and resolves
	Component     string      // Component: Component of the source machine that is responsible for the event
	Group         string      // Group: Logical grouping of components of a service
	Class         string      // Class: The class/type of the event
	CustomDetails interface{} // CustomDetails: Additional details about the event and affected system
	Images        []*Image    // Images: List of images to include
	Links         []*Link     // Links: List of links to include
}
