package be

import (
	"strconv"
	"testing"

	. "github.com/chai2010/assert"
	"github.com/coocood/qbs"
)

func TestSchemaAPI(t *testing.T) {
	CreateAndInitDB()
	db, err := qbs.GetQbs()
	AssertNil(t, err)
	defer func() {
		WipeDB()
		db.Close()
	}()

	testApi, err := NewTestAPI()
	AssertNil(t, err)
	defer testApi.Stop()

	user, err := CreateUser("bronner@soap.example.com", "Dr", "Bronner", false, db)
	AssertNil(t, err)
	_, err = CreatePassword("1234", user.Id, db)
	AssertNil(t, err)
	staff, err := CreateUser("mr-clean@soap.example.com", "Mr", "Clean", true, db)
	AssertNil(t, err)
	_, err = CreatePassword("1234", staff.Id, db)
	AssertNil(t, err)

	AssertGetString(t, testApi.URL()+"/schema")

	userClient, err := NewClient(testApi.URL())
	AssertNil(t, err, "Could not create a client")
	AssertEqual(t, TEST_VERSION, userClient.Schema.API.Version)
	userClient, err = NewClient(testApi.URL())
	AssertNil(t, err, "Could not create another client")
	AssertEqual(t, TEST_VERSION, userClient.Schema.API.Version)
	Assert(t, len(userClient.Schema.Endpoints) >= 3, "Expected at least three endpoints: "+strconv.Itoa(len(userClient.Schema.Endpoints)))
}
