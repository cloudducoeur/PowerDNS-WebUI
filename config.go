package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	PowerDNSURL string `json:"powerdns_url"`
	APIKey      string `json:"api_key"`
	ServerID    string `json:"server_id"`
}

var config Config

func loadConfig() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %v", err)
	}
}
