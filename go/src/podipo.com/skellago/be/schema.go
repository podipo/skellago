package be

import (
	"net/http"
	"net/url"

	"github.com/coocood/qbs"
)

type Schema struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Properties  []Property `json:"properties"`
}

type Property struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DataType    string `json:"data-type"` // string, int, float, array, object
	Optional    bool   `json:"optional"`
}

type SchemaResource struct {
	api *API
}

func NewSchemaResource(api *API) *SchemaResource {
	return &SchemaResource{
		api: api,
	}
}

func (SchemaResource) Name() string  { return "schema" }
func (SchemaResource) Path() string  { return "/schema" }
func (SchemaResource) Title() string { return "The JSON Schema for this API" }
func (SchemaResource) Description() string {
	return "Use this JSON schema to implement your front end API wrapper."
}

func (resource SchemaResource) Properties() []Property {
	properties := []Property{
		Property{
			Name:        "api",
			Description: "Information about this web API",
			DataType:    "object",
		},
	}
	return properties
}

func (tr SchemaResource) Get(vars map[string]string, vals url.Values, requestHeader http.Header, db *qbs.Qbs) (int, interface{}, http.Header) {
	header := map[string][]string{}

	endpoints := make([]Schema, len(tr.api.resources))
	for i, resource := range tr.api.resources {
		endpoints[i] = schemaFromResource(resource)
	}

	data := make(map[string]interface{})
	data["api"] = map[string]string{
		"version": VERSION,
	}
	data["endpoints"] = endpoints
	return 200, data, header
}

func schemaFromResource(resource Resource) Schema {
	schema := Schema{
		Name:        resource.Name(),
		Path:        resource.Path(),
		Title:       resource.Title(),
		Description: resource.Description(),
		Properties:  resource.Properties(),
	}
	return schema
}