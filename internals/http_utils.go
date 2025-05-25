package internals

import (
	"encoding/json"
	"net/http"
	"os"
)

func GetEnvStringOrDefault(key, defaultvalue string) string {
	v, found := os.LookupEnv(key)
	if !found {
		return defaultvalue
	}
	return v
}

// Helper function to send JSON responses
func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	encoder.Encode(data)
}

// Helper function to send error responses
func sendError(w http.ResponseWriter, message string, statusCode int) {
	response := map[string]any{
		"Success": false,
		"Error":   message,
	}
	sendJSON(w, response, statusCode)
}
