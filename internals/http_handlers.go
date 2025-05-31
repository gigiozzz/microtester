package internals

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
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

	// Pretty print to console as well
	fmt.Println("\n🔍 REQUEST DEBUG INFO")
	fmt.Println("=====================")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL.String())
	fmt.Printf("Path: %s\n", r.URL.Path)
	fmt.Printf("Raw Query: %s\n", r.URL.RawQuery)
	fmt.Printf("Remote Address: %s\n", r.RemoteAddr)

	fmt.Println("\n📋 Query Parameters:")
	if totalParams == 0 {
		fmt.Println("  (no parameters)")
	} else {
		for key, values := range query {
			if len(values) == 1 {
				fmt.Printf("  %-15s = %s\n", key, values[0])
			} else {
				fmt.Printf("  %-15s = %v (array with %d values)\n", key, values, len(values))
			}
		}
	}

	fmt.Println("\n📡 Headers:")
	for name, values := range r.Header {
		if len(values) == 1 {
			fmt.Printf("  %-20s = %s\n", name, values[0])
		} else {
			fmt.Printf("  %-20s = %v\n", name, values)
		}
	}
	fmt.Println("=====================")

	response := map[string]any{
		"Success": true,
		"Data": map[string]interface{}{
			"request_info": requestInfo,
			"timestamp":    time.Now().Unix(),
		},
	}
	sendJSON(w, response, http.StatusOK)
}

func TimeoutHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	timeout := query.Get("timeout")
	sleepTime, err := strconv.Atoi(timeout)
	if err != nil {

	}
	start := time.Now()
	time.Sleep(time.Duration(sleepTime) * time.Second)
	end := time.Now()

	response := map[string]any{
		"Success": true,
		"Data": map[string]interface{}{
			"status":    "ok",
			"duration":  fmt.Sprintf("%.3fs", end.Sub(start).Seconds()),
			"timestamp": time.Now().Unix(),
		},
	}
	sendJSON(w, response, http.StatusOK)
}

func DnsResolverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//	oldEnv := os.Getenv("GODEBUG")
	//	os.Setenv("GODEBUG", "netdns=go")
	query := r.URL.Query()
	hostname := query.Get("hostname")
	ips, err := net.LookupIP(hostname)

	success := true
	status := "ok"
	statusCode := http.StatusOK
	key := "ips"
	value := ""
	if err != nil {
		success = false
		status = "fail"
		statusCode = http.StatusInternalServerError
		key = "error"
		value = err.Error()
	} else {
		sep := ""
		fmt.Println("\n🔍 DNS DEBUG INFO")
		fmt.Println("=====================")
		fmt.Printf("Hostname: %s\n", hostname)
		for id, ip := range ips {
			fmt.Printf("ip: %s\n", ip)
			if len(ip) != net.IPv4len {
				continue
			}
			if id != 0 {
				sep = " "
			}
			value = value + sep + ip.String()
		}
	}

	response := map[string]any{
		"Success": success,
		"Data": map[string]any{
			"status":    status,
			"timestamp": time.Now().Unix(),
			"hostname":  hostname,
			key:         value,
		},
	}
	//	os.Setenv("GODEBUG", oldEnv)
	sendJSON(w, response, statusCode)

}

func HttpConnectionResolverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	url := query.Get("url")
	success := true
	status := "ok"
	errorMsg := ""
	statusCode := http.StatusOK

	_, httpStatus, err := httpGetWithTimeout(url, 25)
	if err != nil {
		success = false
		errorMsg = err.Error()
		status = "fail"
		statusCode = http.StatusInternalServerError
	} else {
		statusCode = httpStatus
	}

	response := map[string]any{
		"Success": success,
		"Data": map[string]any{
			"status":     status,
			"timestamp":  time.Now().Unix(),
			"url":        url,
			"statusCode": statusCode,
			"errorMsg":   errorMsg,
		},
	}
	sendJSON(w, response, statusCode)

}

func httpGetWithTimeout(url string, timeout int) (string, int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	return string(body), resp.StatusCode, nil
}
