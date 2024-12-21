# Vertex

Vertex is a high-performance, in-memory key-value store, built in Go, inspired by Redis. Designed for speed and efficiency, Vertex can be used as a database, cache, or message broker.

## Features

- **In-Memory Storage**: Fast read and write operations with support for various data structures.
- **Persistence**: Optional data persistence through snapshotting (RDB) and append-only logs (AOF).
- **Replication**: Master-slave replication for high availability and load balancing.
- **Pub/Sub**: Publish/subscribe messaging capabilities for real-time communication.
- **Eviction Policies**: Flexible memory management with configurable eviction strategies.
- **Single-Threaded Performance**: Optimized for high throughput with atomic operations.
- **Extensibility**: Support for custom modules to extend functionality.

## Getting Started

### Prerequisites

- Go 1.23 or later
- Docker (optional, for containerization)

### Installation

1. Clone the repository:

```sh
git clone https://github.com/AbdessamadEnabih/Vertex.git
cd Vertex
```

1. Build the project:

```sh
make build
```

2. Run the server:

```sh
make run-server
```

3. Run the CLI:

```sh
make run-cli
```

### Running Tests

To run the tests, use the following command:

```sh
make test
```

### Docker

1. Build the Docker image:

```sh
make docker-build
```

2. Run the Docker container:

```sh
make docker-run
```

### Configuration

Configuration files are located in the `config` directory. You can set environment variables to override default configurations.
