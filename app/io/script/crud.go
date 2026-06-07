// script file
package script

import (
	"errors"
	"gost/app/base"
	"io/fs"
	"strings"

	"github.com/dimonrus/gocli"
	"github.com/dimonrus/gomodel"
)

const (
	// FileMod permission for scripts
	FileMod fs.FileMode = 0750
)

func init() {
	base.App.GetScripts()["crud"] = func(args gocli.ArgumentMap) {
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
		crudNumber := args["num"]
		err := crud.Generate(base.App.GetDB(), schema, table, "v1", gomodel.CrudNumber(crudNumber.GetInt()&0xff))
		if err != nil {
			base.App.FatalError(err)
		}
	}
}
