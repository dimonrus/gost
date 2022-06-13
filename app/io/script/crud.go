// script file
package script

import (
	"errors"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/gomodel"
	"gost/app/base"
	"strings"
)

func init() {
	base.App.GetScripts()["crud"] = func(args gocli.Arguments) {
		schema := "public"
		table := args["file"].GetString()
		if table == "" {
			base.App.FatalError(errors.New("table name is empty"))
		}
		if strings.Contains(table, ".") {
			names := strings.Split(table, ".")
			schema = names[0]
			table = names[1]
		}
		db := base.App.GetDB()
		db.Debug = false
		crud := gomodel.NewCRUDGenerator("app/core", "app/client", "app/io/web/api", "gost")
		err := crud.Generate(base.App.GetDB(), schema, table, "v1", uint8(args["num"].GetInt()))
		if err != nil {
			base.App.FatalError(err)
		}
	}
}
