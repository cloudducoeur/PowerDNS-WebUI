![](img/powerdns.png)

## Description

**PowerDNS-WebUI** is a web user interface for reading and searching DNS zones managed by a PowerDNS server via its REST API.

![](img/capture.png)

## Features

- Display DNS zones and their records.
- Advanced search by name, type, or content of DNS records.
- Modern user interface based on [Bulma](https://bulma.io/).

## Prerequisites

- Go 1.24 or higher.
- A PowerDNS server configured with the REST API enabled.

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/cloudducoeur/PowerDNS-WebUI.git
   cd PowerDNS-WebUI
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Build the application:
   ```bash
   make build
   ```

## Usage

### Running the Application

To run the application, use the following command:

```bash
./build/powerdns-webui --config config.toml
```

### Global Options

```bash
GLOBAL OPTIONS:
   --config value        Path to the configuration file (TOML format) [$CONFIG_FILE]
   --help, -h            show help
```

### Configuration Example

Create a `config.toml` file with the following content:

```toml
powerdns_url = "http://localhost:8081"
api_key = "my_awesome_api_key"
server_id = "localhost"
listen_address = "0.0.0.0"
port = "8080"
```

## Development

### Running in Development Mode

To run the application in development mode:

```bash
make run
```

### Tests

To run unit tests:

```bash
make test
```

## Contribution

Contributions are welcome! Please submit a pull request or open an issue to report a problem or suggest an improvement.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
