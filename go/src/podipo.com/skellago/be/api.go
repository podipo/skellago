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

	"github.com/coocood/qbs"
	"github.com/goincremental/negroni-sessions"
	"github.com/gorilla/mux"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"
	PATCH  = "PATCH"

	// Session constants
	AuthCookieName string = "SkellaAuth"
	UserUUIDKey    string = "user-uuid"

	// List URL parameters
	OffsetKey string = "offset"
	LimitKey  string = "limit"
)

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
	A data structure used when returning a list from an API resource
*/
type APIList struct {
	Offset  int         `json:"offset"`
	Limit   int         `json:"limit"`
	Objects interface{} `json:"objects"`
}

/*
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
	Data for a request to an API endpoint
*/
type APIRequest struct {
	PathValues map[string]string
	Values     url.Values
	Header     http.Header
	Body       io.ReadCloser
	DB         *qbs.Qbs
	Session    sessions.Session
	User       *User
}

/*
	An API resource which handles APIRequests to a given URL endpoint
*/
type Resource interface {
	Name() string        // The name used by the mux
	Path() string        // The path used by the mux
	Title() string       // A human readable name which should be one or two words
	Description() string // A long form, human readable description in Markdown
	Properties() []Property
}

/*
	These *Supported interfaces are used by Resources to indicate whether a given HTTP method is supported
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
	The API collects a tree of Resources, manages the mux, and adds the schema resource
*/
type API struct {
	Mux       *mux.Router
	Path      string
	resources []Resource
}

func NewAPI(path string) *API {
	api := &API{
		Mux:       mux.NewRouter(),
		Path:      path,
		resources: make([]Resource, 0),
	}
	api.AddResource(NewSchemaResource(api))
	api.AddResource(NewCurrentUserResource())
	api.AddResource(NewUsersResource())
	api.AddResource(NewUserResource())
	return api
}

func (api *API) AddResource(resource Resource) {
	api.resources = append(api.resources, resource)
	api.Mux.HandleFunc(api.Path+resource.Path(), api.createHandlerFunc(resource)).Name(resource.Name())
}

/*
	Generate the http.HandlerFunc for a given Resource
*/
func (api *API) createHandlerFunc(resource Resource) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
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
			return
		}

		db, err := qbs.GetQbs()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
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
		}

		// Fetch the User from the session
		if session != nil {
			s_uuid := session.Get(UserUUIDKey)
			if s_uuid != nil {
				uuid, _ := s_uuid.(string)
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
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("Skella-Version", VERSION)
		for name, values := range header {
			for _, value := range values {
				rw.Header().Add(name, value)
			}
		}
		rw.WriteHeader(code)
		rw.Write(content)
	}
}
