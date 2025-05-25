package internals

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Health check handler
func HealthHandler(w http.ResponseWriter, r *http.Request) {
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

func CustomStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	success := false
	successStr := "fail"
	status := http.StatusInternalServerError
	statusStr := query.Get("status")
	if statusStr != "" {
		statusInt, err := strconv.Atoi(statusStr)
		if err == nil && http.StatusText(statusInt) != "" {
			status = statusInt
			success = true
			successStr = "ok"
		} else {
			log.Printf("error converting '%s' err:'%s' or not valid number %d", statusStr, err, statusInt)
		}
	} else {
		status = http.StatusOK
		success = true
		successStr = "ok"
	}

	response := map[string]any{
		"Success": success,
		"Data": map[string]interface{}{
			"status":    successStr,
			"timestamp": time.Now().Unix(),
		},
	}

	sendJSON(w, response, status)
}

func EnvListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	envs := geEnvVars()

	response := map[string]any{
		"Success": true,
		"Data": map[string]any{
			"status":       "ok",
			"timestamp":    time.Now().Unix(),
			"environments": envs,
		},
	}

	sendJSON(w, response, http.StatusOK)

}

func geEnvVars() map[string]string {
	envMap := make(map[string]string)
	//log.Printf("environ '%s'", os.Environ())
	for _, env := range os.Environ() {
		//log.Printf("env '%s'", env)
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			key, value := pair[0], pair[1]

			envMap[key] = value
		}
	}

	return envMap
}

func DebugRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Create detailed request info
	requestInfo := map[string]interface{}{
		"method":           r.Method,
		"url":              r.URL.String(),
		"path":             r.URL.Path,
		"raw_query":        r.URL.RawQuery,
		"host":             r.Host,
		"remote_addr":      r.RemoteAddr,
		"user_agent":       r.UserAgent(),
		"headers":          make(map[string]interface{}),
		"query_parameters": make(map[string]interface{}),
		//"query_stats":      make(map[string]interface{}),
	}

	// Pretty print headers
	headers := make(map[string]interface{})
	for name, values := range r.Header {
		if len(values) == 1 {
			headers[name] = values[0]
		} else {
			headers[name] = values
		}
	}
	requestInfo["headers"] = headers

	// Parse all query parameters
	query := r.URL.Query()

	// Pretty print query parameters with detailed info
	queryParams := make(map[string]interface{})
	totalParams := 0

	for key, values := range query {
		totalParams++
		if len(values) == 1 {
			queryParams[key] = map[string]interface{}{
				"value": values[0],
				"type":  "single",
				"count": 1,
			}
		} else {
			queryParams[key] = map[string]interface{}{
				"values": values,
				"type":   "multiple",
				"count":  len(values),
			}
		}
	}

	requestInfo["query_parameters"] = queryParams

	response := map[string]any{
		"Success": true,
		"Data": map[string]interface{}{
			"request_info": requestInfo,
			"timestamp":    time.Now().Unix(),
		},
	}
	sendJSON(w, response, http.StatusOK)
}
