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
	LISTEN_ADDR_KEY     string = "LISTEN_ADDR"
	LISTEN_ADDR_DEFAULT string = ":8080"
	HELATH_PATH_KEY     string = "HEALTH_PATH"
	HELATH_PATH_DEFAULT string = "/healthz"
	CUSTOM_STATUS_PATH  string = "/api/custom-status"
	ENVIRONMENTS_PATH   string = "/api/environments"
	DEBUG_REQUEST_PATH  string = "/api/debug-request"
)

func main() {
	// Setup routes
	mux := http.NewServeMux()
	handler := internals.LoggingMiddleware(mux)
	/*
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux.ServeHTTP(w, r)
		})
	*/
	// Health check
	mux.HandleFunc(internals.GetEnvStringOrDefault(HELATH_PATH_KEY, HELATH_PATH_DEFAULT), internals.HealthHandler)

	mux.HandleFunc(CUSTOM_STATUS_PATH, internals.CustomStatusHandler)
	mux.HandleFunc(ENVIRONMENTS_PATH, internals.EnvListHandler)
	mux.HandleFunc(DEBUG_REQUEST_PATH, internals.DebugRequestHandler)

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
