package main

/*
	Install demo data to show off the system

	Yeah!

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

	err = be.DeleteAllUsers(db)

	alice, err := be.CreateUser("alice@example.com", "Alice", "Smith", true, db)
	if err != nil {
		logger.Fatal("Could not create user", err)
		return
	}
	logger.Print("Created", alice)

	bob, err := be.CreateUser("bob@example.com", "Bob", "Garvey", false, db)
	if err != nil {
		logger.Fatal("Could not create user", err)
		return
	}
	logger.Print("Created", bob)
}
