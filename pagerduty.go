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
	EventActionTrigger     = "trigger"
	EventActionAcknowledge = "acknowledge"
	EventActionResolve     = "resolve"

	SeverityCritical = "critical"
	SeverityError    = "error"
	SeverityWarning  = "warning"
	SeverityInfo     = "info"
)

type Event struct {
	RoutingKey  string   `json:"routing_key"`
	EventAction string   `json:"event_action"`
	DeDupKey    string   `json:"dedup_key,omitempty"`
	Payload     Payload  `json:"payload"`
	Images      []*Image `json:"images,omitempty"`
	Links       []*Link  `json:"links,omitempty"`
}

type Image struct {
	Src  string `json:"src"`  // The source of the image being attached to the incident. This image must be served via HTTPS.
	Href string `json:"href"` // Optional URL; makes the image a clickable link.
	Alt  string `json:"alt"`  // Optional alternative text for the image.
}

type Link struct {
	Href string `json:"href"` // URL of the link to be attached.
	Text string `json:"test"` // Plain text that describes the purpose of the link, and can be used as the link's text.
}

type Payload struct {
	Source        string      `json:"source"`
	Summary       string      `json:"summary"`
	Severity      string      `json:"severity"`
	Component     string      `json:"component,omitempty"`
	Group         string      `json:"group,omitempty"`
	Class         string      `json:"class,omitempty"`
	CustomDetails interface{} `json:"custom_details,omitempty"`
}

type Response struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	DeDupKey string `json:"dedup_key"`
}

type EventOptions struct {
	DeDupKey      string
	Component     string
	Group         string
	Class         string
	CustomDetails interface{}
	Images        []*Image
	Links         []*Link
}

func SendEvent(key, event_action, source, severity, summary string, options EventOptions) (*Response, error) {
	payload := Payload{
		Summary:       summary,
		Source:        source,
		Severity:      severity,
		Component:     options.Component,
		Group:         options.Group,
		Class:         options.Class,
		CustomDetails: options.CustomDetails,
	}

	event := Event{
		RoutingKey:  key,
		EventAction: event_action,
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
