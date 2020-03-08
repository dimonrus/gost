package core

import (
	"gost/app/base"
	"testing"
)

func TestCommon(t *testing.T) {
	base.App.GetLogger().Infoln("Tests are work's fine")
}
