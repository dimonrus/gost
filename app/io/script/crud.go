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
		err := gomodel.GenerateCrud("app/core", "app/client", "gost", schema, table, "", base.App.GetDB())
		if err != nil {
			base.App.FatalError(err)
		}
	}
}
