package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func New() *mux.Router {
	r := mux.NewRouter()

	// middleware
	r.Use(loggingMiddlware)
	r.Use(defaultHeadersMiddleware)

	// user
	r.HandleFunc("/api/v1/user", createUser).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/user", getUser).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/user", updateUser).Methods(http.MethodPatch)

	return r
}
