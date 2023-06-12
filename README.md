# Jellyfin Exporter

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![License](https://img.shields.io/badge/license-MIT-green)

Jellyfin Exporter is a Prometheus exporter for Jellyfin metrics. It fetches data from the Jellyfin API and exposes it in a format that Prometheus can consume.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features

- Fetch Jellyfin server metrics
- Expose data for Prometheus
- Metrics include: active streams, active users, library statistics, etc.

## Prerequisites

- Go version 1.18 or higher
- Jellyfin server up and running
- Prometheus server up and running

## Installation

### From Source

Clone the repository and build using Go:

```sh
git clone https://github.com/rafaelvieiras/jellyfin-exporter.git
cd jellyfin-exporter
go build
```

### Docker

Build the Docker image and run it:

```sh
docker build -t jellyfin-exporter .
docker run -d -p 8080:8080 --name jellyfin-exporter jellyfin-exporter
```

## Usage

1. Set the Jellyfin API URL and API Token as environment variables:

```sh
export JELLYFIN_API_URL="http://your.jellyfin.server:8096"
export JELLYFIN_TOKEN="your-api-token"
```

2. Start Jellyfin Exporter:

```sh
./jellyfin-exporter
```

3. Point your Prometheus server to `http://localhost:8080/metrics` for scraping.

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to get started.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```
