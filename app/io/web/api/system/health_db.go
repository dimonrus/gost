package system

import (
	"github.com/dimonrus/gorest"
	"github.com/dimonrus/porterr"
	"gost/app/base"
	"net/http"
)

// swagger:route GET /system/health-db System HealthDB
//
// System. Health check for db
//
// This will show information about the health of db.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Responses:
//	  200: ResponseMessage
//	  500: ResponseError
func HealthDB(w http.ResponseWriter, r *http.Request) {
	err := base.App.GetDB().Ping()
	if err != nil {
		e := porterr.New(porterr.PortErrorNetwork, "Database is not works")
		gorest.Send(w, gorest.NewErrorJsonResponse(e))
		return
	}
	gorest.Send(w, gorest.NewOkJsonResponse("Database is alive", nil, nil))
}
