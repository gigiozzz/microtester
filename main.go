package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gigiozzz/microtester/internals"
)

const (
	LISTEN_ADDR_KEY               string = "LISTEN_ADDR"
	LISTEN_ADDR_DEFAULT           string = ":8080"
	HELATH_PATH_KEY               string = "HEALTH_PATH"
	HELATH_PATH_DEFAULT           string = "/healthz"
	CUSTOM_METHOD_AND_STATUS_PATH string = "/api/custom-method-and-status"
	ENVIRONMENTS_PATH             string = "/api/environments"
	DEBUG_REQUEST_PATH            string = "/api/debug-request"
	TIMEOUT_PATH                  string = "/api/timeout"
	DNS_PATH                      string = "/api/dns"
	TEST_CONNECTION_PATH          string = "/api/http-connection"
)

func main() {
	// Setup routes
	mux := http.NewServeMux()
	handler := internals.LoggingMiddleware(mux)

	// Health check
	mux.HandleFunc(internals.GetEnvStringOrDefault(HELATH_PATH_KEY, HELATH_PATH_DEFAULT), internals.HealthHandler)

	mux.HandleFunc(CUSTOM_METHOD_AND_STATUS_PATH, internals.CustomStatusHandler)
	mux.HandleFunc(CUSTOM_METHOD_AND_STATUS_PATH+"/*", internals.CustomStatusHandler)
	mux.HandleFunc(ENVIRONMENTS_PATH, internals.EnvListHandler)
	mux.HandleFunc(DEBUG_REQUEST_PATH, internals.DebugRequestHandler)
	mux.HandleFunc(TIMEOUT_PATH, internals.TimeoutHandler)
	mux.HandleFunc(DNS_PATH, internals.DnsResolverHandler)
	mux.HandleFunc(TEST_CONNECTION_PATH, internals.HttpConnectionResolverHandler)

	// endpoint for OOM error

	// github workflow
	/*

			   [![GitHub release](https://img.shields.io/github/release/gigiozzz/microtester.svg)](https://github.com/gigiozzz/microtester/releases)
			   [![GitHub issues](https://img.shields.io/github/issues/gigiozzz/microtester)](https://github.com/gigiozzz/microtester/issues)
			   [![GitHub stars](https://img.shields.io/github/stars/gigiozzz/microtester)](https://github.com/gigiozzz/microtester/stargazers)
			   [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
			   [![Build Status](https://img.shields.io/github/actions/workflow/status/gigiozzz/microtester/build.yml?branch=main)](https://github.com/gigiozzz/microtester/actions)



















		### Testing Different HTTP Methods

		ToDo
		#```bash
		#POST request (will return method not allowed, but you'll see the request details in logs)
		#curl -X POST "http://localhost:3000/api/v1/debug/request?method=post"
		#The debug endpoint only accepts GET, but server logs will show all request details
		#```

		### Network & Connectivity
		| Endpoint | Method | Description |
		|----------|--------|-------------|
		| `/ping/{host}` | GET | Test connectivity to external hosts |

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


	*/

	server := &http.Server{
		Addr:    internals.GetEnvStringOrDefault(LISTEN_ADDR_KEY, LISTEN_ADDR_DEFAULT),
		Handler: handler,
	}

	quit := make(chan os.Signal, 1)
	// Register the channel to receive specific signals
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}

}
