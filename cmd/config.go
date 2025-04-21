package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	PowerDNSURL   string `toml:"powerdns_url"`
	APIKey        string `toml:"api_key"`
	ServerID      string `toml:"server_id"`
	ListenAddress string `toml:"listen_address"`
	Port          string `toml:"port"`
}

var config Config

// loadConfigFromFile loads the application configuration from a TOML file.
//
// Parameters:
// - filePath: The path to the configuration file to be loaded.
//
// Behavior:
// - Decodes the contents of the specified TOML file into the `config` variable.
// - If the file cannot be read or decoded, logs a fatal error and terminates the application.
//
// Dependencies:
// - Uses the `toml.DecodeFile` function from the BurntSushi/toml package to parse the file.
func loadConfigFromFile(filePath string) {
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}
}
