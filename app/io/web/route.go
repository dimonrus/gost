package web

import (
	"github.com/dimonrus/gorest"
	"github.com/dimonrus/goweb"
	"github.com/dimonrus/porterr"
	"github.com/gorilla/mux"
	"gost/app/base"
	"gost/app/io/web/api/system"
	"net/http"
)

// Get routes
func GetRoutes() *mux.Router {

	middleWare := goweb.NewMiddlewareCollection(base.App.GetConfig().Web, base.App.Application, -1)

	routes := mux.NewRouter()
	// Main route
	MainRoute := routes.PathPrefix("/gost").Subrouter()

	// Api routes
	ApiRoute := MainRoute.PathPrefix("/api").Subrouter()

	// SystemRoute
	SystemRoute := ApiRoute.PathPrefix("/system").Subrouter()
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
	SystemRoute.HandleFunc("/health", system.HealthCheckHandler).Methods(http.MethodGet)
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
	SystemRoute.HandleFunc("/health-db", system.HealthDB).Methods(http.MethodGet)
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
	SystemRoute.HandleFunc("/memory", system.MemoryUsage).Methods(http.MethodGet)
	// Swagger
	SystemRoute.HandleFunc("/swagger", system.Swagger).Methods(http.MethodGet)

	// Setup middleware
	routes.Use(middleWare.LoggingMiddleware)
	routes.Use(swaggerCorsMiddleware)
	routes.NotFoundHandler = http.HandlerFunc(middleWare.NotFoundHandler)
	routes.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)
	return routes
}

// CORS headers for swagger
func corsSetupHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "https://petstore.swagger.io/")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,UPDATE,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
}

// CORS Headers middleware
func swaggerCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsSetupHeaders(w)
		next.ServeHTTP(w, r)
	})
}

// Not allowed except swagger
func notAllowed(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == http.MethodOptions {
		corsSetupHeaders(w)
		return
	}
	e := porterr.New(porterr.PortErrorHandler, "Requested method is not allowed. Check it").HTTP(http.StatusMethodNotAllowed)
	gorest.Send(w, gorest.NewErrorJsonResponse(e))
}
