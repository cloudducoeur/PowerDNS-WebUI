package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cloudducoeur/PowerDNS-WebUI/pkg/powerdns"
)

var powerDNSClient *powerdns.PowerDNSClient

func main() {
	loadConfig()

	powerDNSClient = powerdns.NewPowerDNSClient(config.PowerDNSURL, config.APIKey, config.ServerID)

	http.HandleFunc("/", listZonesHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Starting the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
