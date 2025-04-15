package main

import (
	"log"
	"os"

	"github.com/cloudducoeur/PowerDNS-WebUI/pkg/powerdns"
	"github.com/urfave/cli/v2"
)

var powerDNSClient *powerdns.PowerDNSClient

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
			&cli.StringFlag{
				Name:     "powerdns-url",
				Usage:    "PowerDNS API URL",
				EnvVars:  []string{"POWERDNS_URL"},
				Required: false,
			},
			&cli.StringFlag{
				Name:     "api-key",
				Usage:    "PowerDNS API key",
				EnvVars:  []string{"API_KEY"},
				Required: false,
			},
			&cli.StringFlag{
				Name:     "server-id",
				Usage:    "PowerDNS server ID",
				EnvVars:  []string{"SERVER_ID"},
				Required: false,
			},
			&cli.StringFlag{
				Name:     "port",
				Usage:    "Port to run the server on",
				Value:    "8080",
				EnvVars:  []string{"PORT"},
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			configFile := c.String("config")
			if configFile != "" {
				loadConfigFromFile(configFile)
			} else {
				config.PowerDNSURL = c.String("powerdns-url")
				config.APIKey = c.String("api-key")
				config.ServerID = c.String("server-id")
			}

			if config.PowerDNSURL == "" || config.APIKey == "" || config.ServerID == "" {
				log.Fatal("Missing required configuration: powerdns-url, api-key, and server-id must be provided")
			}

			powerDNSClient = powerdns.NewPowerDNSClient(config.PowerDNSURL, config.APIKey, config.ServerID)

			port := c.String("port")
			return StartServer(port)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
