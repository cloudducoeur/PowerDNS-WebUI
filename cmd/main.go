package main

import (
	"log"
	"os"

	"github.com/cloudducoeur/PowerDNS-WebUI/internal/handlers"
	"github.com/cloudducoeur/PowerDNS-WebUI/internal/server"
	"github.com/cloudducoeur/PowerDNS-WebUI/pkg/powerdns"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "PowerDNS-WebUI",
		Usage: "A web UI for reading PowerDNS zones",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Usage:    "Path to the configuration file (TOML format)",
				EnvVars:  []string{"CONFIG_FILE"},
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			configFile := c.String("config")
			if configFile != "" {
				loadConfigFromFile(configFile)
			}

			if config.PowerDNSURL == "" || config.APIKey == "" || config.ServerID == "" {
				log.Fatal("Missing required configuration: powerdns_url, api_key, and server_id must be provided")
			}

			// Initialize the PowerDNS client
			powerDNSClient := powerdns.NewPowerDNSClient(config.PowerDNSURL, config.APIKey, config.ServerID)

			// Pass the client to the handlers
			handlers.SetPowerDNSClient(powerDNSClient)

			// Start the server
			return server.StartServer(config.ListenAddress, config.Port)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
