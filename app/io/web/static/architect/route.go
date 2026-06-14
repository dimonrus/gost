package architect

import (
	"gost/app/base"
	"net/http"

	"github.com/gorilla/mux"
)

// Init register system sub routes
func Init(appRoute *mux.Router, middleware ...mux.MiddlewareFunc) {
	// Architect route
	Architect := appRoute.PathPrefix("/architect").Subrouter()
	// Fileserver
	fileServer := http.FileServer(http.Dir(base.App.GetAbsolutePath("resource/architect/")))
	// Architect app
	Architect.HandleFunc("", ArchitectHandler).Methods(http.MethodGet)
	// File routes
	appRoute.Handle("/{file}", http.StripPrefix("/gost", fileServer))
	// Use middleware
	Architect.Use(middleware...)
}
