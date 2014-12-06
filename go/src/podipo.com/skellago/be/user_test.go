package be

import (
	"testing"

	. "github.com/chai2010/assert"
	"github.com/coocood/qbs"
)

func TestUserAPI(t *testing.T) {
	DropAndCreateTestDB()

	db, err := qbs.GetQbs()
	AssertNil(t, err)
	defer db.Close()

	user, err := CreateUser("adrian@monk.example.com", "Adrian", "Monk", false, db)
	AssertNil(t, err)
	AssertNotEqual(t, user.Id, 0)

	_, err = FindUser(1000000, db)
	AssertNotNil(t, err)

	user2, err := FindUser(user.Id, db)
	AssertNil(t, err)
	AssertEqual(t, user2.Id, user.Id)
	AssertEqual(t, user2.Email, user.Email)

	user2.Email = "crosby@bing.example.com"
	err = UpdateUser(user2, db)
	AssertNil(t, err)
	AssertEqual(t, user2.Id, user.Id)
	user3, err := FindUser(user2.Id, db)
	AssertNil(t, err)
	AssertEqual(t, user2.Email, user3.Email)
}
