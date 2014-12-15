package be

import (
	"testing"

	. "github.com/chai2010/assert"
	"github.com/coocood/qbs"
)

func TestUserAPI(t *testing.T) {
	DropAndCreateTestDB()

	testApi, err := NewTestAPI()
	AssertNil(t, err)
	defer testApi.Stop()

	db, err := qbs.GetQbs()
	AssertNil(t, err)
	defer db.Close()

	user, err := CreateUser("adrian@monk.example.com", "Adrian", "Monk", false, db)
	AssertNil(t, err)

	Assert403(t, "GET", testApi.URL()+"/user/")
	Assert403(t, "GET", testApi.URL()+"/user/"+user.UUID)
}

func TestUser(t *testing.T) {
	DropAndCreateTestDB()

	db, err := qbs.GetQbs()
	AssertNil(t, err)
	defer db.Close()

	user, err := CreateUser("adrian@monk.example.com", "Adrian", "Monk", false, db)
	AssertNil(t, err)
	AssertNotEqual(t, user.UUID, "")

	_, err = FindUser("not-a-uuid", db)
	AssertNotNil(t, err)

	user2, err := FindUser(user.UUID, db)
	AssertNil(t, err)
	AssertEqual(t, user2.UUID, user.UUID)
	AssertEqual(t, user2.Email, user.Email)

	user2.Email = "crosby@bing.example.com"
	err = UpdateUser(user2, db)
	AssertNil(t, err)
	AssertEqual(t, user2.UUID, user.UUID)
	user3, err := FindUser(user2.UUID, db)
	AssertNil(t, err)
	AssertEqual(t, user2.Email, user3.Email)

	// TODO
	/*
		Figure out why test DB isn't dropped
		Figure out API testing
		Test schema API
		Test authentication
		Test versioning enforcement
		Test User API CRUD
		Test Staff-only enforced on User API
	*/
}
