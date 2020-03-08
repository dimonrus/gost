package system

import (
	"github.com/dimonrus/gorest"
	"github.com/dimonrus/porterr"
	"gost/app/base"
	"io/ioutil"
	"net/http"
)

func Swagger(w http.ResponseWriter, r *http.Request) {
	var e porterr.IError
	if file, err := ioutil.ReadFile(base.App.GetAbsolutePath("resource/swagger.json")); err == nil {
		_, err := w.Write(file)
		if err == nil {
			return
		}
		e = porterr.New(porterr.PortErrorIO, "file cannot be write "+err.Error())
	} else {
		e = porterr.New(porterr.PortErrorIO, "file cannot be read "+err.Error())
	}
	gorest.Send(w, gorest.NewErrorJsonResponse(e))
}
