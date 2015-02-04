package main

import (
	"log"
	"os"

	"github.com/coocood/qbs"

	"example.com/api/cms"
	"podipo.com/skellago/be"
)

var logger = log.New(os.Stdout, "[example-demo] ", 0)

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

	err = cms.DeleteAllTags(db)
	if err != nil {
		logger.Fatal("Could not delete tags", err)
		return
	}
	err = cms.DeleteAllEntries(db)
	if err != nil {
		logger.Fatal("Could not delete entries", err)
		return
	}
	err = cms.DeleteAllLogs(db)
	if err != nil {
		logger.Fatal("Could not delete logs", err)
		return
	}

	_, err = cms.CreateLog("Blargh", "blargh", db)

}
