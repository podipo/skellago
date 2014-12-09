package be

import (
	"encoding/json"
	"net/http"
)

type LoginData struct {
	Email    string "json:email"
	Password string "json:password"
}

type CurrentUserResource struct {
}

func NewCurrentUserResource() *CurrentUserResource {
	return &CurrentUserResource{}
}

func (CurrentUserResource) Name() string  { return "current-user" }
func (CurrentUserResource) Path() string  { return "/user/current" }
func (CurrentUserResource) Title() string { return "The logged in User" }
func (CurrentUserResource) Description() string {
	return "The User in the requesting session."
}

func (resource CurrentUserResource) Properties() []Property {
	properties := []Property{
		Property{
			Name:        "uuid",
			Description: "uuid",
			DataType:    "string",
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

func (resource CurrentUserResource) Get(request *APIRequest) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}
	if request.User == nil {
		return 404, "Not logged in", responseHeader
	}
	return 200, request.User, responseHeader
}

func (resource CurrentUserResource) Delete(request *APIRequest) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}
	if request.User == nil {
		return 200, "Ok", responseHeader
	}
	request.Session.Delete(UserUUIDKey)
	logger.Print("Logout: ", request.User.Email)
	return 200, "Ok", responseHeader
}

func (resource CurrentUserResource) Post(request *APIRequest) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}
	var loginData LoginData
	err := json.NewDecoder(request.Body).Decode(&loginData)
	if err != nil {
		return 400, "Error: " + err.Error(), responseHeader
	}
	if loginData.Email == "" || loginData.Password == "" {
		return 400, "Incorrect login", responseHeader
	}
	user, err := FindUserByEmail(loginData.Email, request.DB)
	if err != nil {
		return 400, "No such user", responseHeader
	}
	if PasswordMatches(user.Id, loginData.Password, request.DB) == false {
		return 400, "Incorrect password", responseHeader
	}
	request.Session.Set(UserUUIDKey, user.UUID)
	logger.Print("Login: ", user.Email)
	return 200, user, responseHeader
}

type UserResource struct {
}

func NewUserResource() *UserResource {
	return &UserResource{}
}

func (UserResource) Name() string  { return "user" }
func (UserResource) Path() string  { return "/user/{uuid:UUID[0-9,a-z,-]+}" }
func (UserResource) Title() string { return "The user account record" }
func (UserResource) Description() string {
	return "Each account is associated with a User."
}

func (resource UserResource) Properties() []Property {
	properties := []Property{
		Property{
			Name:        "uuid",
			Description: "uuid",
			DataType:    "string",
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
	uuid, _ := request.PathValues["uuid"]
	user, err := FindUser(uuid, request.DB)
	if err != nil {
		return 404, "No such user: " + uuid, responseHeader
	}
	return 200, user, responseHeader
}

type UsersResource struct {
}

func NewUsersResource() *UsersResource {
	return &UsersResource{}
}

func (UsersResource) Name() string  { return "users" }
func (UsersResource) Path() string  { return "/user/" }
func (UsersResource) Title() string { return "A list of users" }
func (UsersResource) Description() string {
	return "A list of users."
}

func (resource UsersResource) Properties() []Property {
	properties := []Property{
		Property{
			Name:        "resources",
			Description: "the list",
			DataType:    "list",
		},
		// TODO: Add pagination information
	}
	return properties
}

func (resource UsersResource) Get(request *APIRequest) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}
	// TODO limit and offset from vars
	offset := 0
	limit := 100

	users, err := FindUsers(offset, limit, request.DB)
	if err != nil {
		return 500, "Could not list users", responseHeader
	}
	list := &APIList{
		Offset:  offset,
		Limit:   limit,
		Objects: users,
	}
	return 200, list, responseHeader
}
