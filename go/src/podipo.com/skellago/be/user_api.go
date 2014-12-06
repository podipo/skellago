package be

import (
	"net/http"
	"strconv"
)

type UserResource struct {
}

func NewUserResource() *UserResource {
	return &UserResource{}
}

func (UserResource) Name() string  { return "user" }
func (UserResource) Path() string  { return "/user/{id:[0-9]+}" }
func (UserResource) Title() string { return "The user account record" }
func (UserResource) Description() string {
	return "Each account is associated with a User."
}

func (resource UserResource) Properties() []Property {
	properties := []Property{
		Property{
			Name:        "id",
			Description: "id",
			DataType:    "int",
		},
		Property{
			Name:        "email",
			Description: "email",
			DataType:    "string",
		},
		Property{
			Name:        "first-name",
			Description: "first name",
			DataType:    "string",
			Optional:    true,
		},
		Property{
			Name:        "last-name",
			Description: "last name",
			DataType:    "string",
			Optional:    true,
		},
	}
	return properties
}

func (resource UserResource) Get(request *APIRequest) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}
	id, _ := strconv.ParseInt(request.PathValues["id"], 10, 64)
	user, err := FindUser(id, request.DB)
	if err != nil {
		return 404, "No such user: " + strconv.FormatInt(id, 10), responseHeader
	}
	return 200, user, responseHeader
}
