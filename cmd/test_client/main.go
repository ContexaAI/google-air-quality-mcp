package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	// 1. Connect to SSE endpoint
	resp, err := http.Get("http://localhost:8080/mcp")
	if err != nil {
		fmt.Printf("Failed to connect to SSE: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Connected to SSE")

	// 2. Read events to find the endpoint
	scanner := bufio.NewScanner(resp.Body)
	var postEndpoint string

	// We expect an 'endpoint' event
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("SSE Line: %s\n", line)
		if strings.HasPrefix(line, "event: endpoint") {
			// Next line should be data
			if scanner.Scan() {
				dataLine := scanner.Text()
				fmt.Printf("SSE Data: %s\n", dataLine)
				if strings.HasPrefix(dataLine, "data: ") {
					postEndpoint = strings.TrimPrefix(dataLine, "data: ")
					// The endpoint might be relative or absolute.
					// Usually it's a relative path like "/mcp?sessionId=..."
					// Let's assume it's relative if it starts with /
					if strings.HasPrefix(postEndpoint, "/") {
						postEndpoint = "http://localhost:8080" + postEndpoint
					}
					break
				}
			}
		}
	}

	if postEndpoint == "" {
		fmt.Println("Failed to find POST endpoint in SSE stream")
		return
	}

	fmt.Printf("Found POST endpoint: %s\n", postEndpoint)

	// 3. Send resources/read request
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "resources/read",
		"params": map[string]interface{}{
			"uri": "airquality://current/37.7749,-122.4194",
		},
		"id": 1,
	}

	jsonData, _ := json.Marshal(reqBody)
	postReq, err := http.NewRequest("POST", postEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to create POST request: %v\n", err)
		return
	}
	postReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	postResp, err := client.Do(postReq)
	if err != nil {
		fmt.Printf("Failed to send POST request: %v\n", err)
		return
	}
	defer postResp.Body.Close()

	body, _ := io.ReadAll(postResp.Body)
	fmt.Printf("Response Status: %s\n", postResp.Status)
	fmt.Printf("Response Body: %s\n", string(body))
}
