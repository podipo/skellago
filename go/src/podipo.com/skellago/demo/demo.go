package main

/*
	Install demo data to show off the system
*/

import (
	"log"
	"os"

	"github.com/coocood/qbs"
	"podipo.com/skellago/be"
)

var logger = log.New(os.Stdout, "[demo] ", 0)

func main() {
	err := be.InitDB()
	if err != nil {
		logger.Fatal("Could not init the db", err)
		return
	}

	db, err := qbs.GetQbs()
	if err != nil {
		logger.Fatal("Could not get the db", err)
		return
	}
	defer db.Close()

	err = be.DeleteAllPasswords(db)
	if err != nil {
		logger.Fatal(err)
		return
	}
	err = be.DeleteAllUsers(db)
	if err != nil {
		logger.Fatal(err)
		return
	}

	_, err = createUser("alice@example.com", "Alice", "Smith", true, "1234", db)
	if err != nil {
		return
	}

	_, err = createUser("bob@example.com", "Bob", "Garvey", false, "1234", db)
	if err != nil {
		return
	}
}

func createUser(email string, firstName string, lastName string, staff bool, password string, db *qbs.Qbs) (*be.User, error) {
	user, err := be.CreateUser(email, firstName, lastName, staff, db)
	if err != nil {
		logger.Fatal("Could not create user", err)
		return nil, err
	}
	_, err = be.CreatePassword(password, user.Id, db)
	if err != nil {
		logger.Fatal("Could not create password", err)
		return nil, err
	}
	return user, nil
}
