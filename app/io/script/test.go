// script file
package script

import (
	"github.com/dimonrus/gocli"
	"gost/app/base"
)

func init() {
	base.App.GetScripts()["test"] = func(args gocli.Arguments) {
		base.App.GetLogger().Infoln("cron is works, ENV =", base.App.GetENV())
	}
}
