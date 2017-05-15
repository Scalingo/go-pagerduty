package pagerduty

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	errgo "gopkg.in/errgo.v1"
)

var endpoint = "https://events.pagerduty.com/v2/enqueue"

const (
	// EventActionTrigger will trigger a new event
	EventActionTrigger = "trigger"
	// EventActionAcknowledge will acknowledge the current event
	EventActionAcknowledge = "acknowledge"
	// EventActionResolve will resolve the current event
	EventActionResolve = "resolve"

	// SeverityCritical will set the event serverity to critical
	SeverityCritical = "critical"
	// SeverityError will set the event severity to Error
	SeverityError = "error"
	// SeverityWarning will set the event severity to Warning
	SeverityWarning = "warning"
	// SeverityInfo will set the event serverity to info
	SeverityInfo = "info"
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
	DeDupKey      string
	Component     string
	Group         string
	Class         string
	CustomDetails interface{}
	Images        []*Image
	Links         []*Link
}

// SendEvent Send the event to PagerDuty
func SendEvent(key, eventAction, source, severity, summary string, options EventOptions) (*Response, error) {
	payload := payload{
		Summary:       summary,
		Source:        source,
		Severity:      severity,
		Component:     options.Component,
		Group:         options.Group,
		Class:         options.Class,
		CustomDetails: options.CustomDetails,
	}

	event := event{
		RoutingKey:  key,
		EventAction: eventAction,
		DeDupKey:    options.DeDupKey,
		Payload:     payload,
		Images:      options.Images,
		Links:       options.Links,
	}

	buffer := new(bytes.Buffer)

	err := json.NewEncoder(buffer).Encode(&event)

	if err != nil {
		return nil, errgo.Notef(err, "Unable to encode json")
	}

	request, _ := http.NewRequest("POST", endpoint, buffer)
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errgo.Notef(err, "Unable to send event")
	}
	defer resp.Body.Close()

	if resp.Status[0] != '2' {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New("Invalid return code: " + resp.Status + " " + string(body))
	}

	response := &Response{}
	err = json.NewDecoder(resp.Body).Decode(response)

	if err != nil {
		return nil, errgo.Notef(err, "Unable to read response")
	}

	return response, nil
}
