# 🚀 Microtester

A lightweight REST API server built with Go's standard library, designed specifically for container deployment. MicroTester provides essential endpoints for testing, checking, and debugging HTTP requests and REST API deployments in containerized environments.

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-green.svg)](https://kubernetes.io/)
[![Docker Image Size](https://img.shields.io/docker/image-size/gigiozzz/microtester/latest)](https://hub.docker.com/r/gigiozzz/microtester)
[![Docker Pulls](https://img.shields.io/docker/pulls/gigiozzz/microtester)](https://hub.docker.com/r/gigiozzz/microtester)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## ✨ Features

- **🔍 Request Debugging**: Comprehensive HTTP request inspection and parameter analysis
- **⚡ Lightweight**: Minimal footprint with Go's standard library
- **🐳 Docker Ready**: Multi-stage Dockerfile for production-optimized containers (~5MB)
- **☸️ Kubernetes Native**: Graceful shutdown, health checks, and configurable ports
- **📊 Detailed Logging**: Request/response logging with precise timing information
- **🛡️ Security First**: Non-root container user and minimal attack surface
- **⚙️ Environment Configurable**: Port and settings via environment variables

## 🚀 Quick Start

### Docker
```bash
docker run -p 8080:8080 gigiozzz/microtester
```

### Kubernetes
```bash
kubectl -n your-namespace apply -f https://raw.githubusercontent.com/gigiozzz/microtester/refs/heads/main/manifest.yaml
```

## 📋 Examples

### Basic Health Check
```bash
curl http://localhost:8080/healthz
# Response: {"Data": { "status": "ok", "timestamp": 1748295254 }, "Success": true}
```

### HTTP Request Debugging
```bash
curl http://localhost:8080/debug-request
# The response body and the server logs will show detailed information about the HTTP request including all parameters, headers, and metadata.
```

### Environment Debugging
```bash
curl http://localhost:8080/api/environments
# Returns all environment variables for debugging
```

### DNS Debugging
```bash
curl http://localhost:8080/api/dns?hostname={hostname}
# Returns all ipv4 record found as dns resolve result for debugging
```

### Delay Testing
```bash
curl http://localhost:8080/api/timeout?timeout=15
# Waits 15 seconds before responding
```

### Testing Different HTTP Methods and Statuses with custom sub-path
```bash
curl "http://localhost:8080/api/custom-method-and-status?status=405"
# or with custom sub-path
curl "http://localhost:8080/api/custom-method-and-status/any-path-you-need?status=405"
# Returns 405 http status
```

### Testing HTTP GET connection to an url
```bash
curl "http://localhost:8080/api/http-connection?url=https://www.google.com/search?q=microtester"
# Returns 200 or 500 as http status with details in the body
```

### Testing in Kubernetes
```bash
# Port forward to access locally
kubectl port-forward service/microtester-service 8080:80

# Test the service
curl "http://localhost:8080/api/api/environments"
```


## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LISTEN_ADDR` | `:8080` | Server address and port toi listen to |
| `HEALTH_PATH` | `/healthz` | The path to use for a Kubernetes liveness endpoint |


## 📋 API Endpoints

### Health & Status
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/healthz` | GET | Basic health check endpoint |

### Testing & Debugging
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/environments` | GET | Whill show all environment variables |
| `/api/debug-request` | All | Echo back request details (headers, body, method) |
| `/api/dns?hostname={hostname}` | GET | DNS resolution testing  |
| `/custom-method-and-status?status={status}` | All | Test HTTP requests with any method and with custom return status equal to the query parameter `status` value, default status 200 |
| `/custom-method-and-status/{any-path-you-need}-?status={status}` | All | Test HTTP requests with any method and with custom return status equal to the query parameter `status` value, default status 200 |
| `/api/timeout?timeout={seconds}` | GET | Simulate response delays for testing |
| `/api/http-connection?url={url}` | GET | Execute an HTTP GET connection to an url |

## 🏗️ Building from Source

```bash
git clone https://github.com/gigiozzz/microtester.git
cd microtester
make build
```

## 🤝 Contributing
Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📄 License
This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Docker Hub](https://hub.docker.com/r/gigiozzz/microtester)
- [Issues](https://github.com/gigiozzz/microtester/issues)
- [Releases](https://github.com/gigiozzz/microtester/releases)

---

**Made with ❤️ by gigiozzz for the container and Kubernetes community**