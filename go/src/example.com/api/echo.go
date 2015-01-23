package main

import (
	"net/http"
	"time"

	"podipo.com/skellago/be"
)

/*
EchoProperties describe this resource and are exported in the schema JSON
*/
var EchoProperties = []be.Property{
	be.Property{
		Name:        "text",
		Description: "The echo'ed text",
		DataType:    "string",
	},
	be.Property{
		Name:        "time",
		Description: "A timestamp, just for fun",
		DataType:    "date-time",
	},
}

// This is the struct that will be serialized to JSON and returned
// This is the struct that will be serialized to JSON and returned
// This is the struct that will be serialized to JSON and returned
// This is the struct that will be serialized to JSON and returned
type EchoResponse struct {
	Text string    `json:"text"`
	Time time.Time `json:"time"`
}

/*
EchoResource is a simple example resource
It simply returns what it was passed
*/
type EchoResource struct{}

func NewEchoResource() *EchoResource {
	return &EchoResource{}
}

func (EchoResource) Name() string  { return "echo" }  // The slug for the resource
func (EchoResource) Path() string  { return "/echo" } // The path (under /api/.../ for the resource)
func (EchoResource) Title() string { return "Echo" }  // A human readable name
func (EchoResource) Description() string {
	return "A simple example API resource." // A human readable description used in API docs
}

/*
Return the list which describes the Resource
*/
func (resource EchoResource) Properties() []be.Property {
	return EchoProperties
}

/*
Get is called when HTTP hits this resource
Create a Resource func for each HTTP method the resource supports
*/
func (resource EchoResource) Get(request *be.APIRequest) (int, interface{}, http.Header) {

	// The response header is a map of HTTP headers for the response
	responseHeader := map[string][]string{}

	// This is an example of setting a header
	responseHeader["X-Example"] = []string{"Howdy, folks"}

	// If the user is authenticated, request.User will be a be.User, otherwise nil
	if request.User == nil {
		return 404, be.NotLoggedInError, responseHeader
	}

	// Check for a URL parameter named 'text'
	text := request.Raw.FormValue("text")

	// This will be serialized to JSON and sent in the response
	responseData := EchoResponse{
		Text: text,
		Time: time.Now(),
	}

	// Return the HTTP status, an object that can be serialized to JSON, and the HTTP headers
	return 200, responseData, responseHeader
}
