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

func loadConfigFromFile(filePath string) {
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}
}
