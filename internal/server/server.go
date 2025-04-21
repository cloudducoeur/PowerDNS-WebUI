package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cloudducoeur/PowerDNS-WebUI/internal/handlers"
)

// StartServer initializes and starts an HTTP server.
//
// Parameters:
// - listenAddress: The IP address or hostname on which the server will listen.
// - port: The port number on which the server will listen.
//
// The function sets up two HTTP handlers:
// 1. "/" - This is handled by the `ListZonesHandler` function from the `handlers` package.
// 2. "/static/" - This serves static files from the "static" directory, stripping the "/static/" prefix.
//
// The server address is constructed using the provided `listenAddress` and `port`.
// It logs the server's start address and begins listening for incoming HTTP requests.
//
// Returns:
// - An error if the server fails to start or encounters an issue during execution.
func StartServer(listenAddress, port string) error {
	http.HandleFunc("/", handlers.ListZonesHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	address := fmt.Sprintf("%s:%s", listenAddress, port)
	log.Printf("Server started on %s", address)
	return http.ListenAndServe(address, nil)
}
