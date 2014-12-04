package be

/*
	The Resource and API functionality.
*/

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"
	PATCH  = "PATCH"
)

/*
	An API resource which handles requests to a given URL endpoint
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
	Get(url.Values, http.Header) (int, interface{}, http.Header)
}
type PostSupported interface {
	Post(url.Values, http.Header) (int, interface{}, http.Header)
}
type PutSupported interface {
	Put(url.Values, http.Header) (int, interface{}, http.Header)
}
type DeleteSupported interface {
	Delete(url.Values, http.Header) (int, interface{}, http.Header)
}
type HeadSupported interface {
	Head(url.Values, http.Header) (int, interface{}, http.Header)
}
type PatchSupported interface {
	Patch(url.Values, http.Header) (int, interface{}, http.Header)
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
	schemaResource := NewSchemaResource(api)
	api.AddResource(schemaResource)
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
		var methodHandler func(url.Values, http.Header) (int, interface{}, http.Header)
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
		code, data, header := methodHandler(request.Form, request.Header)
		content, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
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
