package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	URL      = "https://srv.msk01.gigacorp.local"
	INTERVAL = 2 * time.Second
	TIMEOUT  = 5 * time.Second
)

var errorCount int

func main() {
	client := &http.Client{
		Timeout: TIMEOUT,
	}

	for {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, URL+"/_stats", nil)
		if err != nil {
			fmt.Printf("Request creation error: %v\n", err)
			incrementError()
			time.Sleep(INTERVAL)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Request error: %v\n", err)
			incrementError()
		} else {
			ok := handleResponse(resp)
			if !ok {
				incrementError()
			} else {
				errorCount = 0 // reset on success
			}
		}

		time.Sleep(INTERVAL)
	}
}

func handleResponse(resp *http.Response) bool {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	string_data := strings.Split(strings.TrimSpace(string(bodyBytes)), ",")
	data := make([]int, len(string_data))
	for i, s := range string_data {
		var value int
		if _, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &value); err == nil {
			data[i] = value
		} else {
			return false // format error
		}
	}

	// Print stats if needed
	fmt.Printf("Server stats: %v\n", data)
	return true
}

func incrementError() {
	errorCount++
	if errorCount >= 3 {
		fmt.Println("Unable to fetch server stats. Please check the server.")
	}
}
