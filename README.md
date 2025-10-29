# ts-viewer

TeamSpeak Server Viewer - A web-based viewer for TeamSpeak servers written in Go.

## Features

- **Stateless HTTP server** with embedded assets
- **Configurable via environment variables**
- **Server-side rendering** using Go templates
- **Docker support** with multi-stage builds
- **Health check endpoint** for monitoring

## Quick Start

### Prerequisites

- Go 1.22+ (for local development)
- Docker (optional, for containerized deployment)

### Running Locally

1. Clone the repository:
```bash
git clone https://github.com/MetaEMK/ts-viewer.git
cd ts-viewer
```

2. Run the server:
```bash
go run ./cmd/server
```

3. Open your browser and navigate to:
   - Main UI: http://localhost:8080
   - Health check: http://localhost:8080/healthz

### Using Docker

1. Build the Docker image:
```bash
docker build -t ts-viewer .
```

2. Run the container:
```bash
docker run -p 8080:8080 ts-viewer
```

3. Access the application at http://localhost:8080

### Configuration

The application is configured via a YAML configuration file (`config.yaml`). Copy the example file to get started:

```bash
cp config.yaml.example config.yaml
```

Then edit `config.yaml` to configure your TeamSpeak servers:

```yaml
# HTTP server configuration
http_addr: ":8080"
log_level: "info"

# TeamSpeak servers configuration
# Each server must have a unique name which will be used in the URL path
servers:
  # Access via: http://localhost:8080/ts-view/production
  production:
    host: "ts.example.com"
    port: 10011
    username: "serveradmin"  # Optional
    password: "secret123"    # Optional

  # Access via: http://localhost:8080/ts-view/dev
  dev:
    host: "192.168.1.100"
    port: 10011
```

#### Configuration Options

**Application Settings:**
- `http_addr`: HTTP server listen address (default: `:8080`)
- `log_level`: Logging level (default: `info`)

**Server Configuration:**
Each server in the `servers` map requires:
- `host`: TeamSpeak server hostname or IP address (required)
- `port`: ServerQuery port (default: `10011` if not specified)
- `username`: ServerQuery username (optional, for authenticated connections)
- `password`: ServerQuery password (optional, for authenticated connections)

#### Environment Variables

You can override the config file path using:
- `TS_CONFIG_FILE`: Path to the configuration file (default: `config.yaml`)

Example:
```bash
TS_CONFIG_FILE=/etc/ts-viewer/config.yaml ./ts-viewer
```
```

## Project Structure

```
ts-viewer/
├── cmd/
│   └── server/          # Main application entrypoint
├── internal/
│   ├── assets/          # Embedded web assets (templates, CSS)
│   ├── config/          # Configuration management
│   ├── server/          # HTTP server and handlers
│   └── tsviewer/        # TeamSpeak data models and providers
├── Dockerfile           # Multi-stage Docker build
└── README.md
```

## Development

### Building

```bash
go build -o ts-viewer ./cmd/server
./ts-viewer
```

### Testing

```bash
go test ./...
```

### Linting

```bash
go fmt ./...
go vet ./...
```

## Current Status

The application connects to real TeamSpeak servers via ServerQuery protocol. Current features:
- ✅ Serves a web UI showing a TeamSpeak-like channel tree
- ✅ Real-time TeamSpeak ServerQuery integration
- ✅ Support for multiple TeamSpeak servers via configuration
- ✅ Path-based server selection (`/ts-view/{server-name}`)
- ✅ YAML configuration file support
- ✅ Layered architecture (API, Business, Data layers)
- ✅ Health check endpoint
- ✅ Docker support

Access configured servers:
- Default dummy data: `http://localhost:8080/`
- Real server by name: `http://localhost:8080/ts-view/{server-name}`

Future enhancements will include:
- 🔲 Authentication support (username/password)
- 🔲 Live data updates (WebSocket)
- 🔲 Server selection UI
- 🔲 Additional features

## License

[Add your license here]