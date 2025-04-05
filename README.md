# Projects in Go

Repository with small projects written in Go. Below is a short description of each.

## Projects

### Load balancer

This load balancer is designed to manage HTTP/HTTPS traffic at the application layer (L7),distributing requests based on application-level data (URLs, headers, cookies, etc...) among multiple backend servers to ensure reliability and scalability.

#### Goals

- LB send traffic to two or more servers;
- Health check the servers;
- Handle a server going offline (failing a health check); and
- Handle a server coming back online (passing a health check).

### Tools

A collection of command-line tools including utility programs. Each tool is organized within its own module.

---

Happy Coding!