package main

import (
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile)

	h := newHub()

	mux := http.ServeMux{}
	mux.Handle("/", http.FileServer(http.Dir("client")))
	mux.HandleFunc("/playable", h.handlePlayable)
	mux.HandleFunc("/ws", h.serveWS)

	server := http.Server{
		Addr:    ":" + "8080",
		Handler: &mux,
	}

	log.Printf("server listening on port: %v", "8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error listening to server: %v\n", err)
	}
}
