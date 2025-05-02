# Projects in Go

This repository contains small projects written in Go to explore different areas of backend development.

## Projects

### Load Balancer

This load balancer is designed to manage HTTP/HTTPS traffic at the application layer (L7), distributing incoming requests based on application-level data (URLs, headers, cookies, etc...) between backend servers using the simple scheduling algorithm round robin.

#### Features

- ✅ Distributes traffic across multiple backend servers (round-robin)
- ✅ Performs periodic health checks on each backend (`/health`)
- ✅ Automatically removes backends that fail health checks
- ✅ Automatically re-adds healthy backends once they recover
- ✅ Supports concurrent requests
- ✅ Configurable health check interval via CLI flag (`--healthcheck-interval`)
- ✅ Written in pure Go — no external dependencies
- ✅ Includes test script to validate routing and fault tolerance

#### Architecture Overview

```scss
Client
  ↓
Load Balancer (:9090)
  ↙       ↓       ↘
Backend1 Backend2 Backend3
(:8080)  (:8081)  (:8082)
```

#### Usage

```bash
# Build and run
go build -o lb main.go
./lb --healthcheck-interval=10s -backends=http://localhost:8080,http://localhost:8081,http://localhost:8082
```

#### Health Check Logic

Every <interval> seconds (default: 10s), the load balancer sends a GET /health request to each backend.
If the response is HTTP 200 OK, the server is marked as healthy. Otherwise, it is temporarily removed from rotation. Once it passes again, it is re-added.

#### Testing

A helper script to test is available in test/test.sh:

```bash
./test/test.sh
```

This script:

- Builds and starts 3 mock backends
- Starts the load balancer
- Sends 6 requests to verify round-robin behavior
- Simulates a failing backend and ensures traffic is only routed to healthy ones
- Simulates recovery and verifies the backend is added back
- Tests concurrent request handling using `curl --parallel`

### Tools (WIP)

A collection of command-line tools including utility programs. Each tool is organized within its own module.

---

Happy Coding!