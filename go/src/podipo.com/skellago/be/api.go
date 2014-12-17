package be

/*
	The Resource and API functionality.
*/

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/coocood/qbs"
	"github.com/goincremental/negroni-sessions"
	"github.com/gorilla/mux"
)

// GET and the other HTTP methods
const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"
	PATCH  = "PATCH"
)

// AuthCookieName and UserUUID are used by the session mechanism
const (
	AuthCookieName string = "SkellaAuth"
	UserUUIDKey    string = "user-uuid"
)

// OffsetKey and LimitKey are list URL parameters
const (
	OffsetKey string = "offset"
	LimitKey  string = "limit"
)

// AcceptHeaderPrefix should be followed by the version in the Accept header of requests
const (
	AcceptHeaderPrefix = "application/vnd.api+json; version="
)

// APIListProperties are the properties any resource which is a list
var APIListProperties = []Property{
	Property{
		Name:        "offset",
		Description: "The index into the result set for the first item in this list",
		DataType:    "int",
	},
	Property{
		Name:        "limit",
		Description: "The maximum number of results returned",
		DataType:    "int",
	},
	Property{
		Name:        "objects",
		Description: "The array of objects in this list",
		DataType:    "array",
	},
}

/*
APIList is a data structure used when returning a list from an API resource
*/
type APIList struct {
	Offset  int         `json:"offset"`
	Limit   int         `json:"limit"`
	Objects interface{} `json:"objects"`
}

/*
GetOffsetAndLimit find the range values from a request's url.Values
By default, return 0, 100
If limit and offset are set in the values, return those
*/
func GetOffsetAndLimit(values url.Values) (offset int, limit int) {
	offsetVal, err := strconv.Atoi(values.Get(OffsetKey))
	if err == nil {
		offset = offsetVal
	} else {
		offset = 0
	}
	limitVal, err := strconv.Atoi(values.Get(LimitKey))
	if err == nil {
		limit = limitVal
	} else {
		limit = 100
	}
	return
}

/*
APIRequest is data for a request to an API endpoint
*/
type APIRequest struct {
	PathValues map[string]string
	Values     url.Values
	Header     http.Header
	Body       io.ReadCloser
	DB         *qbs.Qbs
	Session    sessions.Session
	User       *User
	Version    string
}

/*
Resource is and API resource which handles APIRequests to a given URL endpoint
*/
type Resource interface {
	Name() string        // The name used by the mux
	Path() string        // The path used by the mux
	Title() string       // A human readable name which should be one or two words
	Description() string // A long form, human readable description in Markdown
	Properties() []Property
}

/*
GetSupported and the other interfaces are used by Resources to indicate whether a given HTTP method is supported
*/
type GetSupported interface {
	Get(request *APIRequest) (int, interface{}, http.Header)
}
type PostSupported interface {
	Post(request *APIRequest) (int, interface{}, http.Header)
}
type PutSupported interface {
	Put(request *APIRequest) (int, interface{}, http.Header)
}
type DeleteSupported interface {
	Delete(request *APIRequest) (int, interface{}, http.Header)
}
type HeadSupported interface {
	Head(request *APIRequest) (int, interface{}, http.Header)
}
type PatchSupported interface {
	Patch(request *APIRequest) (int, interface{}, http.Header)
}

/*
API collects a tree of Resources, manages the mux, and adds the schema resource
*/
type API struct {
	Mux       *mux.Router
	Path      string
	Version   string
	resources []Resource
}

func NewAPI(path string, version string) *API {
	api := &API{
		Mux:       mux.NewRouter(),
		Path:      path,
		Version:   version,
		resources: make([]Resource, 0),
	}
	api.AddResource(NewSchemaResource(api), false)
	api.AddResource(NewCurrentUserResource(), true)
	api.AddResource(NewUsersResource(), true)
	api.AddResource(NewUserResource(), true)
	return api
}

func (api *API) AddResource(resource Resource, versioned bool) {
	api.resources = append(api.resources, resource)
	api.Mux.HandleFunc(api.Path+resource.Path(), api.createHandlerFunc(resource, versioned)).Name(resource.Name())
}

func (api *API) acceptableAcceptHeader(acceptTypes []string) bool {
	if len(acceptTypes) == 0 {
		return false
	}
	acceptTypes = strings.Split(acceptTypes[0], ",")
	for _, acceptType := range acceptTypes {
		if acceptType == AcceptHeaderPrefix+api.Version {
			return true
		}
	}
	return false
}

/*
	Generate the http.HandlerFunc for a given Resource
*/
func (api *API) createHandlerFunc(resource Resource, versioned bool) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		if versioned && !api.acceptableAcceptHeader(request.Header["Accept"]) {
			rw.WriteHeader(http.StatusBadRequest)
			errorString, _ := json.MarshalIndent(IncorrectVersionError, "", "")
			rw.Write(errorString)
			return
		}
		if request.ParseForm() != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		var methodHandler func(*APIRequest) (int, interface{}, http.Header)
		switch request.Method {
		case GET:
			if resource, ok := resource.(GetSupported); ok {
				methodHandler = resource.Get
			}
		case POST:
			if resource, ok := resource.(PostSupported); ok {
				methodHandler = resource.Post
			}
		case PUT:
			if resource, ok := resource.(PutSupported); ok {
				methodHandler = resource.Put
			}
		case DELETE:
			if resource, ok := resource.(DeleteSupported); ok {
				methodHandler = resource.Delete
			}
		case HEAD:
			if resource, ok := resource.(HeadSupported); ok {
				methodHandler = resource.Head
			}
		case PATCH:
			if resource, ok := resource.(PatchSupported); ok {
				methodHandler = resource.Patch
			}
		}
		if methodHandler == nil {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			errorString, _ := json.MarshalIndent(MethodNotAllowedError, "", "")
			rw.Write(errorString)
			return
		}

		db, err := qbs.GetQbs()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			jError := APIError{
				Id:      "db_error",
				Message: "Database error: " + err.Error(),
			}
			errorString, _ := json.MarshalIndent(jError, "", "")
			rw.Write(errorString)
			return
		}
		defer db.Close()

		session := sessions.GetSession(request)

		apiRequest := &APIRequest{
			PathValues: mux.Vars(request),
			Values:     request.Form,
			Header:     request.Header,
			Body:       request.Body,
			DB:         db,
			Session:    sessions.GetSession(request),
			Version:    api.Version,
		}

		// Fetch the User from the session
		if session != nil {
			sUUID := session.Get(UserUUIDKey)
			if sUUID != nil {
				uuid, _ := sUUID.(string)
				user, err := FindUser(uuid, db)
				if err == nil {
					apiRequest.User = user
				}
			}
		}

		code, data, header := methodHandler(apiRequest)
		content, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			jError := APIError{
				Id:      "json_serialization_error",
				Message: "JSON serialization error: " + err.Error(),
			}
			errorString, _ := json.MarshalIndent(jError, "", "")
			rw.Write(errorString)
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("API-Version", api.Version)
		rw.Header().Add("Request-Id", UUID()) // Useful for tracking requests across the front and back end
		for name, values := range header {
			for _, value := range values {
				rw.Header().Add(name, value)
			}
		}

		// Check whether the client's If-None-Match and the response header's ETag match
		if rw.Header().Get("Etag") != "" && rw.Header().Get("Etag") == request.Header.Get("If-None-Match") {
			rw.WriteHeader(http.StatusNotModified)
			return
		}

		rw.WriteHeader(code)
		rw.Write(content)
	}
}
