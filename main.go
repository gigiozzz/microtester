package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	// Setup routes
	mux := http.NewServeMux()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	})

	// Health check
	mux.HandleFunc("/healthz", healthHandler)

	log.Fatal(http.ListenAndServe(":3000", handler))
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
