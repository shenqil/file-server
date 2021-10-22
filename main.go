package main

import (
	"fileServer/app"
	"fileServer/app/config"
	"log"
	"net/http"
)

func main() {

	config.MustLoad("./configs/config.toml")

	port := ":" + config.C.HTTP.Port

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
