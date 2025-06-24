package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)

	h := newHub()

	mux := http.ServeMux{}
	mux.Handle("/", http.FileServer(http.Dir("client")))
	mux.HandleFunc("/playable", h.handlePlayable)
	mux.HandleFunc("/ws", h.serveWS)

	port := os.Getenv("PORT")
	server := http.Server{
		Addr:    ":" + port,
		Handler: &mux,
	}

	log.Printf("server listening on port: %v", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error listening to server: %v\n", err)
	}
}
