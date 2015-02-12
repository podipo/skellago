package main

import (
	"log"
	"os"
	"time"

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

	log, err := cms.CreateLog("Blargh", "blargh", db)
	if err != nil {
		logger.Fatal("Could not create a log", err)
		return
	}
	log.Publish = true
	log.Tagline = "Stuff and thangs."
	err = cms.UpdateLog(log, db)
	if err != nil {
		logger.Fatal("Could not publish the log", err)
		return
	}

	entry1, err := cms.CreateEntry(log, "First post", "first", "This is the first post.", db)
	if err != nil {
		logger.Fatal("Could not create an entry", err)
		return
	}
	entry1.Publish = true
	entry1.Issued = time.Now()
	err = cms.UpdateEntry(entry1, db)
	if err != nil {
		logger.Fatal("Could not update the entry", err)
		return
	}

	entry2, err := cms.CreateEntry(log, "Second post", "second", "This is the second post.", db)
	if err != nil {
		logger.Fatal("Could not create an entry", err)
		return
	}
	entry2.Publish = true
	entry2.Issued = time.Now()
	err = cms.UpdateEntry(entry2, db)
	if err != nil {
		logger.Fatal("Could not update the entry", err)
		return
	}

	log2, err := cms.CreateLog("Flapdoodle and Flippers", "flapdoodle-and-flippers", db)
	if err != nil {
		logger.Fatal("Could not create a log", err)
		return
	}
	log2.Publish = true
	log2.Tagline = "All things F."
	err = cms.UpdateLog(log2, db)
	if err != nil {
		logger.Fatal("Could not publish the log", err)
		return
	}

	entry3, err := cms.CreateEntry(log2, "Flivvers", "flivvers", "For forest forgetfulness.", db)
	if err != nil {
		logger.Fatal("Could not create an entry", err)
		return
	}
	entry3.Publish = true
	entry3.Issued = time.Now()
	err = cms.UpdateEntry(entry3, db)
	if err != nil {
		logger.Fatal("Could not update the entry", err)
		return
	}
}
