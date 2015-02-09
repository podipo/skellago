package main

import (
	"github.com/coocood/qbs"

	"example.com/api/cms"
)

func migrateDB() error {
	migration, err := qbs.GetMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	migration.CreateTableIfNotExists(new(cms.Log))
	migration.CreateTableIfNotExists(new(cms.Entry))
	migration.CreateTableIfNotExists(new(cms.Tag))
	return nil
}

func wipeDB() error {
	db, err := qbs.GetQbs()
	if err != nil {
		return err
	}
	tables := []string{"tag", "entry", "log"}
	for _, table := range tables {
		_, err = db.Exec("delete from " + table)
		if err != nil {
			return err
		}
	}
	return nil
}
