package main

import (
	"log"
	"medods/app"
	"net/http"
)

func main() {
	router := app.Initialize()
	log.Println("Server run is 8000 port...")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}
