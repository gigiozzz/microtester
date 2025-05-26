# 🚀 Microtester

A lightweight REST API server built with Go's standard library, designed specifically for container deployment. MicroTester provides essential endpoints for testing, checking, and debugging HTTP requests and REST API deployments in containerized environments.

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Multi--stage-blue.svg)](https://docs.docker.com/build/building/multi-stage/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-green.svg)](https://kubernetes.io/)
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
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: microtester
spec:
  replicas: 1
  selector:
    matchLabels:
      app: microtester
  template:
    metadata:
      labels:
        app: microtester
    spec:
      containers:
      - name: microtester
        image: gigiozzz/microtester:latest
        env:
        - name: PORT
          value: "8080"
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 10
        readinessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 5          
        resources:
          requests:
            memory: "16Mi"
            cpu: "10m"
          limits:
            memory: "64Mi"
            cpu: "100m"
---
apiVersion: v1
kind: Service
metadata:
  name: microtester-service
spec:
  selector:
    app: microtester
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

## 📋 Examples

### Basic Health Check
```bash
curl http://localhost:8080/healthz
# Response: {"Data": { "status": "ok", "timestamp": 1748295254 }, "Success": true}
```

### Test with Payload
```bash
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'
```

### Environment Debugging
```bash
curl http://localhost:8080/api/environments
# Returns all environment variables for debugging
```

### Delay Testing
```bash
curl http://localhost:8080/api/timeout?timeout=15
# Waits 15 seconds before responding
```

### Testing Different HTTP Methods

```bash
# POST request (will return method not allowed, but you'll see the request details in logs)
curl -X POST "http://localhost:3000/api/v1/debug/request?method=post"

# The debug endpoint only accepts GET, but server logs will show all request details
```

### Testing in Kubernetes

```bash
# Port forward to access locally
kubectl port-forward service/microtester-service 3000:80

# Test the service
curl "http://localhost:3000/api/v1/debug/request?k8s=test&cluster=local"
```
### Request Debugging
```http
GET /api/v1/debug/request
```
Inspects and returns detailed information about the HTTP request including all parameters, headers, and metadata.

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | Server port |



## 📋 API Endpoints

### Health & Status
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Basic health check endpoint |


### Testing & Debugging
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/echo` | ALL | Echo back request details (headers, body, method) |
| `/test/get` | GET | Test GET requests with query parameters |
| `/test/post` | POST | Test POST requests with JSON payload |
| `/test/delay/{seconds}` | GET | Simulate response delays for testing |

### Network & Connectivity
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/ping/{host}` | GET | Test connectivity to external hosts |
| `/dns/{hostname}` | GET | DNS resolution testing |
| `/proxy/{url}` | GET | Proxy requests to test service mesh |



## 🔍 Use Cases

### API Development Testing
- **Parameter validation**: Test how your API handles different parameter combinations
- **Query string debugging**: Inspect complex URL encoding scenarios
- **Header analysis**: Verify custom headers are properly transmitted

### Integration Testing
- **Webhook testing**: Use as a target for webhook deliveries
- **Load balancer verification**: Check request routing and header forwarding
- **Proxy testing**: Verify proxy configurations and header modifications

### Debugging Workflows
- **Request inspection**: See exactly what your client is sending
- **Parameter parsing**: Understand how complex query strings are interpreted
- **Header debugging**: Analyze authentication headers and custom values

## 🏗️ Building from Source

```bash
git clone https://github.com/gigiozzz/microtester.git
cd microtester
docker build -t microtester .
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