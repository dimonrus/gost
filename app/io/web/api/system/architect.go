package system

import (
	"gost/app/base"
	"net/http"
)

// ResponseArchitectHandler Architect response
//
// swagger:response ResponseArchitectHandler
type ResponseArchitectHandler struct{}

// swagger:route GET /system/architect System ArchitectHandler
//
// System. Architect your app
//
// This help to you build perfect world.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Responses:
//	  200: ResponseArchitectHandler
func ArchitectHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, base.App.GetAbsolutePath("resource/architect/index.html"))
	return
}
