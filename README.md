# Jelly Metrics

**Jelly Metrics** is a lightweight Prometheus exporter for [Jellyfin](https://jellyfin.org/). It exposes key metrics
related to your Jellyfin instance, including:

* Media count by type.
* Connected clients by username.
* Actively playing streams by username.

I use this tool in my [Ultimate Jellyfin on macOS guide](https://github.com/Digital-Shane/jellyfin-on-macos?tab=readme-ov-file#monitoring).
A grafana dashboard using metrics from this exporter and prometheus scrape config is provided in that guide. 

---

## Table of Contents

1. [Features](#features)
2. [Getting Started](#getting-started)
   1. [Prerequisites](#prerequisites)
   2. [Installation](#installation)
   3. [Run the Exporter](#run-the-exporter)
3. [Configuration](#configuration)
4. [Available Metrics](#available-metrics)
   * [jellyfin_media_count](#jellyfin_media_count)
   * [jellyfin_connected_clients_count](#jellyfin_connected_clients_count)
   * [jellyfin_stream_count](#jellyfin_stream_count)
5. [Contributing](#contributing)
6. [Changelog](#changelog)
7. [License](#license)

---

## Features

- **Prometheus Integration**: Jelly Metrics exposes `/metrics` endpoint for easy scraping by Prometheus.
- **Minimal Cardinality**: The exporter aims to provide only essential metrics to prevent large data ingestion footprints.
- **Lightweight**: Uses Go’s built-in HTTP server and minimal external dependencies to run efficiently.

---

## Getting Started

### Prerequisites

- [Go](https://go.dev/) 1.23.6 or higher (if you intend to build from source)
- A running [Jellyfin](https://jellyfin.org/) server
- A valid Jellyfin API key

### Installation

```bash
go install github.com/Digital-Shane/jelly-metrics@latest
```

### Run the Exporter

```bash
export JELLYFIN_TOKEN="YOUR_JELLYFIN_API_KEY"
# Optional environment variables:
# export JELLYFIN_HOST="http://localhost:8096"
# export PORT="8097"

./jelly-metrics
```

---

## Configuration

| Environment Variable | Description                                                               | Default                      |
|----------------------|---------------------------------------------------------------------------|------------------------------|
| **JELLYFIN_TOKEN**   | **Required** – Jellyfin API token. Without this, the exporter cannot run. | *none* (must be provided)    |
| **JELLYFIN_HOST**    | URL of your Jellyfin server.                                              | `http://localhost:8096`      |
| **PORT**             | Port on which the jelly-metrics exporter listens.                         | `8097`                       |

---

## Available Metrics

Jelly Metrics defines three primary metrics:

### jellyfin_media_count

Gauge metric keyed by `type`. Reflects the number of different media types available on your Jellyfin server.

**Example:**
```
jellyfin_media_count{type="albums"} 120
jellyfin_media_count{type="movies"} 53
```

### jellyfin_connected_clients_count

Gauge metric keyed by `username`. Shows the count of connected (but not necessarily playing) client sessions per user.

**Example:**
```
jellyfin_connected_clients_count{username="alice"} 1
jellyfin_connected_clients_count{username="bob"} 2
```

### jellyfin_stream_count

Gauge metric keyed by `username`. Indicates the number of actively playing streams for each user.

**Example:**
```
jellyfin_stream_count{username="alice"} 1
jellyfin_stream_count{username="bob"} 0
```

---

## Contributing

Contributions are welcome! If you have any suggestions or encounter a bug, please open an [issue](https://github.com/your-org/jelly-metrics/issues)
or submit a pull request.

When contributing:

1. Fork the repository and create a new feature branch.
2. Make your changes in a well-structured commit history.
3. Include tests (when applicable).
4. Submit a pull request with a clear description of your changes.

---

## Changelog

A detailed list of changes and release notes is maintained in [CHANGELOG.md](./CHANGELOG.md). Refer to it for information
about new features, fixes, and updates in each version.

---

## License

Project uses MIT License, view the full details in the [LICENSE](./LICENSE) file.

---

# Happy monitoring!
