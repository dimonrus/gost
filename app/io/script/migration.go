// script file
package script

import (
	"errors"
	"github.com/dimonrus/gocli"
	"gost/app/base"
	_ "gost/app/io/db/migrations/data"
	_ "gost/app/io/db/migrations/schema"
)

func init() {
	base.App.GetScripts()["migration"] = func(args gocli.Arguments) {
		base.App.GetLogger().Info("Start migrations...")
		// Set db form migration
		base.App.GetMigration().DBO = base.App.GetDB()
		// Get type of migration
		class := args["class"].GetString()
		if class == "" {
			base.App.FatalError(errors.New("no migration class specified"))
		}
		// Get filename
		file := args["file"].GetString()
		if file != "" {
			err := base.App.GetMigration().CreateMigrationFile(class, file)
			if err != nil {
				base.App.FatalError(errors.New("create migration file error: " + err.Error()))
			}
			return
		}
		// Init type migration
		err := base.App.GetMigration().InitMigration(class)
		if err != nil {
			base.App.FatalError(errors.New("init migration error: " + err.Error()))
		}

		// Execute migrations
		err = base.App.GetMigration().Upgrade(class)
		if err != nil {
			base.App.FatalError(errors.New("execute migration error: " + err.Error()))
		}
	}
}
