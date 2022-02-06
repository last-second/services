package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func createUser(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, "OK")
}

func getUser(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, "OK")
}

func New() *mux.Router {
	r := mux.NewRouter()

	// middleware
	r.Use(loggingMiddlware)

	// user
	r.HandleFunc("/api/v1/user", createUser).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/user", getUser).Methods(http.MethodGet)

	return r
}
