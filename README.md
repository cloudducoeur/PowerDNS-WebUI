![](img/powerdns.png)

## Usage

```bash
NAME:
   PowerDNS-WebUI - A web UI for read PowerDNS zones

USAGE:
   PowerDNS-WebUI [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value        Path to the configuration file (TOML format) [$CONFIG_FILE]
   --powerdns-url value  PowerDNS API URL [$POWERDNS_URL]
   --api-key value       PowerDNS API key [$API_KEY]
   --server-id value     PowerDNS server ID [$SERVER_ID]
   --port value          Port to run the server on (default: "8080") [$PORT]
   --help, -h            show help
```

## Configuration

```bash
vi config.toml
...
powerdns_url = "http://localhost:8081"
api_key = "my_awesome_api_key"
server_id = "localhost"
```
