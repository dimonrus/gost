package system

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Init register system sub routes
func Init(apiRoute *mux.Router, middleware ...mux.MiddlewareFunc) {
	// SystemRoute
	SystemRoute := apiRoute.PathPrefix("/system").Subrouter()
	// App health
	SystemRoute.HandleFunc("/health", HealthCheckHandler).Methods(http.MethodGet)
	// DB health
	SystemRoute.HandleFunc("/health-db", HealthDB).Methods(http.MethodGet)
	// Memory usage
	SystemRoute.HandleFunc("/memory", MemoryUsage).Methods(http.MethodGet)
	// Swagger route
	SystemRoute.HandleFunc("/swagger", Swagger).Methods(http.MethodGet)
	// Use middleware
	apiRoute.Use(middleware...)
}
