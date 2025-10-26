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

The application is configured via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | HTTP server listen address |
| `TS_SERVER_URL` | `""` | TeamSpeak server URL (not used yet) |
| `TS_API_TOKEN` | `""` | TeamSpeak API token (not used yet) |
| `TS_VIRTUAL_SERVER_ID` | `""` | Virtual server ID (not used yet) |
| `LOG_LEVEL` | `info` | Logging level |

Example with custom configuration:
```bash
HTTP_ADDR=:3000 LOG_LEVEL=debug go run ./cmd/server
```

Or with Docker:
```bash
docker run -p 3000:3000 -e HTTP_ADDR=:3000 ts-viewer
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

This is an initial implementation using dummy/static data. The application currently:
- ✅ Serves a web UI showing a TeamSpeak-like channel tree
- ✅ Displays dummy server data (channels and clients)
- ✅ Provides health check endpoint
- ✅ Supports configuration via environment variables
- ✅ Includes Docker support

Future enhancements will include:
- 🔲 Real TeamSpeak API integration
- 🔲 Live data updates
- 🔲 Additional features

## License

[Add your license here]