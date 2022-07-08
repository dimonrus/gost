package system

import (
	"fmt"
	"github.com/dimonrus/gorest"
	"github.com/dimonrus/porterr"
	"gost/app/base"
	"io/ioutil"
	"net/http"
	"runtime"
)

// swagger:route GET /system/health system HealthCheckHandler
//
// System. Health check
//
// This will show information about health of service.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Responses:
//       200: ResponseMessage
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	gorest.Send(w, gorest.NewOkJsonResponse(
		fmt.Sprintf("%s service is alive", base.App.GetConfig().Project.Name), nil, nil),
	)
}

// swagger:route GET /system/health-db system HealthDB
//
// System. Health check for db
//
// This will show information about health of db.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Responses:
//       200: ResponseMessage
//       500: ResponseError
func HealthDB(w http.ResponseWriter, r *http.Request) {
	err := base.App.GetDB().Ping()
	if err != nil {
		e := porterr.New(porterr.PortErrorNetwork, "Database is not works")
		gorest.Send(w, gorest.NewErrorJsonResponse(e))
		return
	}
	gorest.Send(w, gorest.NewOkJsonResponse("Database is alive", nil, nil))
}

// swagger:route GET /system/memory system MemoryUsage
//
// System. Memory usage
//
// This will show information about memory of service.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Responses:
//       200: ResponseMemoryUsage
func MemoryUsage(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	report := make(map[string]string)
	report["allocated"] = fmt.Sprintf("%v KB", m.Alloc/1024)
	report["total_allocated"] = fmt.Sprintf("%v KB", m.TotalAlloc/1024)
	report["system"] = fmt.Sprintf("%v KB", m.Sys/1024)
	report["garbage_collectors"] = fmt.Sprintf("%v", m.NumGC)
	gorest.Send(w, gorest.NewOkJsonResponse("Memory usage", report, nil))
}

// Swagger route
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
