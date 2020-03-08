package system

import (
	"fmt"
	"github.com/dimonrus/gorest"
	"github.com/dimonrus/porterr"
	"gost/app/base"
	"net/http"
	"runtime"
)

// Health Check
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	gorest.Send(w, gorest.NewOkJsonResponse("Alive", nil, nil))
}

// Ping db
func HealthDB(w http.ResponseWriter, r *http.Request) {
	err := base.App.GetDB().Ping()
	if err != nil {
		e := porterr.New(porterr.PortErrorNetwork, "Database is not works")
		gorest.Send(w, gorest.NewErrorJsonResponse(e))
		return
	}
	gorest.Send(w, gorest.NewOkJsonResponse("Database is alive", nil, nil))
}

// Memory Usage
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
