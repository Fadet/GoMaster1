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
	url      = "srv.msk01.gigacorp.local"
	interval = 0 * time.Second
	timeout  = 5 * time.Second
)

var errorCount int

func main() {
	client := &http.Client{
		Timeout: timeout,
	}

	for {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url+"/stats", nil)
		if err != nil {
			fmt.Printf("Request creation error: %v\n", err)
			incrementError()
			time.Sleep(interval)
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
				errorCount = 0
			}
		}

		time.Sleep(interval)
	}
}

func handleResponse(resp *http.Response) bool {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK || resp.Header.Get("Content-Type") != "text/plain" {
		return false
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	stringData := strings.Split(strings.TrimSpace(string(bodyBytes)), ",")
	if len(stringData) != 7 {
		return false
	}

	data := make([]int, len(stringData))
	for i, s := range stringData {
		var value int
		if _, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &value); err == nil {
			data[i] = value
		} else {
			return false
		}
	}

	processData(data[0], data[1], data[2], data[3], data[4], data[5], data[6])

	return true
}

func incrementError() {
	errorCount++
	if errorCount >= 3 {
		fmt.Println("Unable to fetch server statistic")
	}
}

func processData(loadAvg, ramTotal, ramUsed, diskTotal, diskUsed, bandwidthTotal, bandwidthUsed int) {
	if loadAvg > 30 {
		fmt.Printf("Load Average is too high: %d\n", loadAvg)
	}
	if ramUsed > ramTotal*80/100 {
		fmt.Printf("Memory usage too high: %d%%\n", ramUsed*100/ramTotal)
	}
	if diskUsed > diskTotal*90/100 {
		fmt.Printf("Free disk space is too low: %d Mb left\n", diskTotal-diskUsed)
	}
	if bandwidthUsed > bandwidthTotal*90/100 {
		fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", bandwidthTotal-bandwidthUsed)
	}
}
