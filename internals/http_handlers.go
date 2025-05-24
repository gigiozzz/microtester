package internals

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	LISTEN_ADDR_KEY            string = "LISTEN_ADDR"
	LISTEN_ADDR_DEFAULT        string = ":8080"
	HELATH_PATH_KEY            string = "HEALTH_PATH"
	HELATH_PATH_DEFAULT        string = "/healthz"
	CUSTOM_STATUS_PATH_DEFAULT string = "/api/custom-status"
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
