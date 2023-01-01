package system

import (
	"fmt"
	"github.com/dimonrus/gorest"
	"net/http"
	"runtime"
)

// MemoryInfo Memory usage information
type MemoryInfo struct {
	// Allocated memory
	// Required: true
	// Example: 100 KB
	Allocated string `json:"allocated"`
	// Total allocated
	// Required: true
	// Example: 300 KB
	TotalAllocated string `json:"totalAllocated"`
	// System allocated memory
	// Required: true
	// Example: 200 KB
	System string `json:"system"`
	// Count of GC cycle
	// Required: true
	// Example: 1
	GarbageCollectors string `json:"garbageCollectors"`
}

// ResponseMemoryUsage Memory usage response
//
// swagger:response ResponseMemoryUsage
type ResponseMemoryUsage struct {
	// In: body
	Body struct {
		// Message
		// Required: true
		// Example: Memory usage
		Message string `json:"message,omitempty"`
		// Memory usage data
		// Required: true
		Data MemoryInfo `json:"data"`
	}
}

// swagger:route GET /system/memory System MemoryUsage
//
// System. Memory usage
//
// This will show information about memory of service.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Responses:
//	  200: ResponseMemoryUsage
func MemoryUsage(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	info := MemoryInfo{
		Allocated:         fmt.Sprintf("%v KB", m.Alloc/1024),
		TotalAllocated:    fmt.Sprintf("%v KB", m.TotalAlloc/1024),
		System:            fmt.Sprintf("%v KB", m.Sys/1024),
		GarbageCollectors: fmt.Sprintf("%v", m.NumGC),
	}
	gorest.Send(w, gorest.NewOkJsonResponse("Memory usage", info, nil))
}
