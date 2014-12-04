package be

import (
	"github.com/coocood/qbs"
)

func MigrateDB() error {
	logger.Print("Migrating DB Schema if necessary")
	migration, err := qbs.GetMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	migration.CreateTableIfNotExists(new(User))
	logger.Print("Migrated DB Schema")
	return nil
}
