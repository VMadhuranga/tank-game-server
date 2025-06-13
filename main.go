package main

import (
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile)

	mux := http.ServeMux{}
	mux.Handle("/", http.FileServer(http.Dir("client")))

	server := http.Server{
		Addr:    ":" + "8080",
		Handler: &mux,
	}

	log.Printf("server listening on port: %v", "8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error listening to server: %v", err)
	}
}
