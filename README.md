# Go HTTP Server with Graceful Shutdown

This is a simple HTTP server written in Go that supports graceful shutdown, configurable port, and optional request delay through query parameters.

## Features

- Basic HTTP server responding to `GET /`
- Optional delay using `sleep` query parameter (e.g., `/?sleep=5`)
- Graceful shutdown on `SIGINT` or `SIGTERM`
- Configurable port via `PORT` environment variable
- Configurable shutdown grace period via `GRACE_PERIOD_DURATION` (default 30 seconds)
- Can skip signal handling using `NO_SIGNALS` env var

## Requirements

- Go 1.21 or later (uses `cmp.Or` introduced in Go 1.21)

## Usage

### Run the server

```bash
go run main.go
```

By default, it listens on port 8080. You can override this by setting the PORT environment variable:

```bash
PORT=9090 go run main.go
```

### Optional Parameters

- sleep: Causes the server to sleep for the given number of seconds before responding. Useful for simulating slow responses.

Example:

```bash
curl "http://localhost:8080/?sleep=3"
```

### Graceful Shutdown

To shut down the server gracefully, send a SIGINT or SIGTERM signal (e.g., Ctrl+C or Kubernetes pod termination). The server will wait for GRACE_PERIOD_DURATION seconds (default 30) before shutting down.

You can configure the grace period via:

```bash
GRACE_PERIOD_DURATION=10 go run main.go
```

### Disabling Signal Handling
You can configure to force shutdown via:

```bash
NO_SIGNALS=true go run main.go
```
