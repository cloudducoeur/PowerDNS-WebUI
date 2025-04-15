package main

import (
	"log"
	"net/http"
)

// StartServer initializes and starts the HTTP server.
func StartServer(port string) error {
	http.HandleFunc("/", listZonesHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Server started on port %s", port)
	return http.ListenAndServe(":"+port, nil)
}
