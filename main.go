package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	LISTEN_ADDR_KEY     string = "LISTEN_ADDR"
	LISTEN_ADDR_DEFAULT string = ":8080"
	HELATH_PATH_KEY     string = "HEALTH_PATH"
	HELATH_PATH_DEFAULT string = "/healthz"
)

func main() {
	// Setup routes
	mux := http.NewServeMux()
	handler := loggingMiddleware(mux)
	/*
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux.ServeHTTP(w, r)
		})
	*/
	// Health check
	mux.HandleFunc(getEnvStringOrDefault(HELATH_PATH_KEY, HELATH_PATH_DEFAULT), healthHandler)

	server := &http.Server{
		Addr:    getEnvStringOrDefault(LISTEN_ADDR_KEY, LISTEN_ADDR_DEFAULT),
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

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]any{
		"Success": true,
		"Data": map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		},
	}

	sendJSON(w, response, http.StatusOK)
}

func getEnvStringOrDefault(key, defaultvalue string) string {
	v, found := os.LookupEnv(key)
	if !found {
		return defaultvalue
	}
	return v
}

// ResponseWriter wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Middleware for logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Wrap the response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code
		}
		next.ServeHTTP(wrapped, r)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
	})
}

// Helper function to send JSON responses
func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Helper function to send error responses
func sendError(w http.ResponseWriter, message string, statusCode int) {
	response := map[string]any{
		"Success": false,
		"Error":   message,
	}
	sendJSON(w, response, statusCode)
}
