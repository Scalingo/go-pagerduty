package pagerduty

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	errgo "gopkg.in/errgo.v1"
)

// API describe the PagerDuty API Client
type API interface {
	// SendEvent will send the event to the api.
	// The field are:
	// * EventAction: the type of event
	// * string: the source of the event
	// * Severity: the perceived severity of the status event
	// * string: the summary of the event
	// * EventOptions: Optional parameters for the Event
	SendEvent(EventAction, string, Severity, string, EventOptions) (*Response, error)
}

// Client is the structure used to create a PagerDuty Client
type Client struct {
	apiKey string
}

// NewClient generate a new client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

// SendEvent send an event to the api.
func (c *Client) SendEvent(eventAction EventAction, source string, severity Severity, summary string, options EventOptions) (*Response, error) {

	payload := payload{
		Summary:       summary,
		Source:        source,
		Severity:      string(severity),
		Component:     options.Component,
		Group:         options.Group,
		Class:         options.Class,
		CustomDetails: options.CustomDetails,
	}

	event := event{
		RoutingKey:  c.apiKey,
		EventAction: string(eventAction),
		DeDupKey:    options.DeDupKey,
		Payload:     payload,
		Images:      options.Images,
		Links:       options.Links,
	}

	// Marshal data to JSON
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(&event)
	if err != nil {
		return nil, errgo.Notef(err, "Unable to encode json")
	}

	// Prepare request
	request, _ := http.NewRequest("POST", endpoint, buffer)
	request.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errgo.Notef(err, "Unable to send event")
	}
	defer resp.Body.Close()

	// If the status was'nt 2XX
	if resp.Status[0] != '2' {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New("Invalid return code: " + resp.Status + " " + string(body))
	}

	// Unmarshal response
	response := &Response{}
	err = json.NewDecoder(resp.Body).Decode(response)

	if err != nil {
		return nil, errgo.Notef(err, "Unable to read response")
	}

	return response, nil
}
