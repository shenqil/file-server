package app

import (
	"fileServer/handle"
	"github.com/gorilla/mux"
	"net/http"
)

// Handler .
func Handler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", handle.Home).Methods("GET")
	router.PathPrefix("/upload/").Handler(http.StripPrefix("/upload/", handle.FileUpload()))
	router.Use(mux.CORSMethodMiddleware(router))

	return router
}
