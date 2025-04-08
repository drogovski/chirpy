package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	const port = "8080"
	myHandler := http.NewServeMux()

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      myHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
