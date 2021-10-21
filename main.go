package main

import (
	"fileServer/app"
	"log"
	"net/http"
	"os"
)

const defaultPort = ":3000"

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = defaultPort
	}

	handler := app.Handler()

	finish := make(chan bool)
	go func() {
		if err := http.ListenAndServe(port, handler); err != nil {
			panic(err)
		}
	}()
	log.Println("Listening to http://localhost" + port)
	<-finish
}