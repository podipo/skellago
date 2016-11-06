package cms

import (
	"github.com/coocood/qbs"
)

func MigrateDB() error {
	migration, err := qbs.GetMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	migration.CreateTableIfNotExists(new(Log))
	migration.CreateTableIfNotExists(new(Entry))
	migration.CreateTableIfNotExists(new(Tag))
	return nil
}

func WipeDB() error {
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
