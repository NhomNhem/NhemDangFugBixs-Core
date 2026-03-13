package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	var baseURL string
	var levelID string
	var concurrency int
	var requests int

	flag.StringVar(&baseURL, "url", "http://localhost:8080", "Base URL of the API")
	flag.StringVar(&levelID, "level", "level_1", "Level ID to query")
	flag.IntVar(&concurrency, "c", 50, "Number of concurrent workers")
	flag.IntVar(&requests, "n", 1000, "Total number of requests to send")
	flag.Parse()

	fmt.Printf("Starting load test on %s/api/v1/leaderboards/%s\n", baseURL, levelID)
	fmt.Printf("Concurrency: %d, Total Requests: %d\n", concurrency, requests)

	var wg sync.WaitGroup
	reqChan := make(chan int, requests)
	resChan := make(chan time.Duration, requests)
	errChan := make(chan error, requests)

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{Timeout: 5 * time.Second}
			for range reqChan {
				start := time.Now()
				resp, err := client.Get(fmt.Sprintf("%s/api/v1/leaderboards/%s?page=1&perPage=20", baseURL, levelID))
				if err != nil {
					errChan <- err
					continue
				}
				resp.Body.Close()
				if resp.StatusCode >= 400 {
					errChan <- fmt.Errorf("bad status code: %d", resp.StatusCode)
					continue
				}
				resChan <- time.Since(start)
			}
		}()
	}

	// Send requests
	startTime := time.Now()
	for i := 0; i < requests; i++ {
		reqChan <- i
	}
	close(reqChan)

	// Wait for workers to finish
	wg.Wait()
	close(resChan)
	close(errChan)
	totalTime := time.Since(startTime)

	// Collect results
	var totalDuration time.Duration
	var minDuration = time.Hour
	var maxDuration time.Duration
	successCount := 0
	errorCount := 0

	for d := range resChan {
		successCount++
		totalDuration += d
		if d < minDuration {
			minDuration = d
		}
		if d > maxDuration {
			maxDuration = d
		}
	}

	for range errChan {
		errorCount++
	}

	fmt.Println("\n--- Load Test Results ---")
	fmt.Printf("Total Time: %v\n", totalTime)
	fmt.Printf("Successful Requests: %d\n", successCount)
	fmt.Printf("Failed Requests: %d\n", errorCount)
	if successCount > 0 {
		fmt.Printf("Requests/sec: %.2f\n", float64(successCount)/totalTime.Seconds())
		fmt.Printf("Average Latency: %v\n", totalDuration/time.Duration(successCount))
		fmt.Printf("Min Latency: %v\n", minDuration)
		fmt.Printf("Max Latency: %v\n", maxDuration)
	}
}
