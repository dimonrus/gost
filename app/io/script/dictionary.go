// script file
package script

import (
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/gomodel"
	"gost/app/base"
)

func init() {
	base.App.GetScripts()["dictionary"] = func(args gocli.Arguments) {
		var err error
		path := args["file"].GetString()
		if path != "" {
			err = gomodel.GenerateDictionaryMapping(path, base.App.GetDB())
		}
		if err != nil {
			base.App.FatalError(err)
		}
	}
}
