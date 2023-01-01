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

// ResponseError JSON Answer. Common error response
//
// swagger:response ResponseError
type ResponseError struct {
	// In: body
	Body struct {
		// Error
		// Required: true
		Error porterr.PortError `json:"error"`
	}
}

// ResponseMessage Common message response
//
// swagger:response ResponseMessage
type ResponseMessage struct {
	// In: body
	Body struct {
		// System message
		// Required: true
		// Example: Success
		Message string `json:"message,omitempty"`
	}
}

// Get routes
func GetRoutes() *mux.Router {
	middleWare := goweb.NewMiddlewareCollection(base.App.GetConfig().Web, base.App.Application, -1)

	routes := mux.NewRouter()
	// Main route
	MainRoute := routes.PathPrefix("/gost").Subrouter()

	// Api routes
	ApiRoute := MainRoute.PathPrefix("/api").Subrouter()

	// System sub route
	system.Init(ApiRoute)

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
