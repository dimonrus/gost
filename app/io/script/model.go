package script

import (
	"errors"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/godb"
	"gost/app/base"
	"strings"
)

func init() {
	base.App.GetScripts()["model"] = func(args gocli.Arguments) {
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
		err := godb.MakeModel(base.App.GetDB(), "app/io/db/models", schema, table, "vendor/github.com/dimonrus/godb/model.tmpl", godb.DefaultSystemColumnsSoft)
		if err != nil {
			base.App.FatalError(err)
		}
	}
}
