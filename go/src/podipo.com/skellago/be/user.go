package be

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/coocood/qbs"
)

type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
}

var dummyUsers = []User{
	User{
		Id:        1,
		Email:     "adrian@monk.example.com",
		FirstName: "Adrian",
		LastName:  "Monk",
	},
	User{
		Id:        2,
		Email:     "sharona@monk.example.com",
		FirstName: "Sharona",
		LastName:  "Fleming",
	},
}

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

func (resource UserResource) Get(vars map[string]string, values url.Values, requestHeader http.Header, db *qbs.Qbs) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}
	id, _ := strconv.Atoi(vars["id"])
	if id < 1 || id > len(dummyUsers) {
		return 404, "No such user: " + strconv.Itoa(id), responseHeader
	}
	return 200, dummyUsers[id-1], responseHeader
}
