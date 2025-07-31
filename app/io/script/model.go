package script

import (
	"errors"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/gomodel"
	"gost/app/base"
	"strings"
)

func init() {
	base.App.GetScripts()["model"] = func(args gocli.ArgumentMap) {
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
		_, _, err := gomodel.MakeModel(base.App.GetDB(), "app/io/db/models", schema, table, "", gomodel.DefaultSystemColumnsSoft)
		if err != nil {
			base.App.FatalError(err)
		}
	}
}
