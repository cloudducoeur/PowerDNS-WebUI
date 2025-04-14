package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/cloudducoeur/PowerDNS-WebUI/pkg/powerdns"
)

var powerDNSClient *powerdns.PowerDNSClient

func main() {
	configFile := flag.String("config", "", "Path to the configuration file (TOML format)")
	powerDNSURL := flag.String("powerdns-url", "", "PowerDNS API URL")
	apiKey := flag.String("api-key", "", "PowerDNS API key")
	serverID := flag.String("server-id", "", "PowerDNS server ID")
	port := flag.String("port", "8080", "Port to run the server on")

	flag.Parse()

	if *configFile != "" {
		loadConfigFromFile(*configFile)
	} else {
		config.PowerDNSURL = *powerDNSURL
		config.APIKey = *apiKey
		config.ServerID = *serverID
	}

	if config.PowerDNSURL == "" || config.APIKey == "" || config.ServerID == "" {
		log.Fatal("Missing required configuration: powerdns-url, api-key, and server-id must be provided")
	}

	powerDNSClient = powerdns.NewPowerDNSClient(config.PowerDNSURL, config.APIKey, config.ServerID)

	http.HandleFunc("/", listZonesHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Server started on port %s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func loadConfigFromFile(filePath string) {
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}
}
