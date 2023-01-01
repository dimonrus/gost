package system

import (
	"fmt"
	"github.com/dimonrus/gorest"
	"gost/app/base"
	"net/http"
)

// BuildInfo build information
type BuildInfo struct {
	// Tag
	Tag string `json:"tag"`
	// Commit hash
	Commit string `json:"commit"`
	// Release time
	Release string `json:"release"`
}

// Load via ldflags information here
var (
	Tag     = "unknown"
	Commit  = "unknown"
	Release = "unknown"
)

// ResponseHealthCheckHandler Health check response
//
// swagger:response ResponseHealthCheckHandler
type ResponseHealthCheckHandler struct {
	// In: body
	Body struct {
		// Message
		// Required: true
		// Example: Memory usage
		Message string `json:"message,omitempty"`
		// Build info
		// Required: true
		Data BuildInfo `json:"data"`
	}
}

// swagger:route GET /system/health System HealthCheckHandler
//
// System. Health check
//
// This will show information about health of service.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Responses:
//	  200: ResponseHealthCheckHandler
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	info := BuildInfo{
		Tag:     Tag,
		Commit:  Commit,
		Release: Release,
	}
	message := fmt.Sprintf("%s service is alive", base.App.GetConfig().Project.Name)
	gorest.Send(w, gorest.NewOkJsonResponse(message, info, nil))
}
