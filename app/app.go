package app

import (
	"fileServer/app/handle"
	"github.com/gorilla/mux"
	"net/http"
)

// Handler .
func Handler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", handle.Home).Methods("GET")
	router.HandleFunc("/ping", handle.Ping)
	router.PathPrefix("/files/").Handler(http.StripPrefix("/files/", handle.FileUpload()))

	return router
}
