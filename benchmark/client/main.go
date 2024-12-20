package main

import (
	"os"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsURL           = "ws://localhost:8080/ws"
	pollURL         = "http://localhost:8081/poll"
	messageURL      = "http://localhost:8081/message"
	maxClients      = 1000
	messagesPerClient = 10
	clientStep      = 50
	pollTimeout       = 5 * time.Second
)

type Message struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

func benchmarkWebSocket(numClients int) time.Duration {
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			for j := 0; j < messagesPerClient; j++ {
				conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
				if err != nil {
					log.Printf("Error connecting to WebSocket: %v", err)
					time.Sleep(time.Second) // Wait before retrying
					continue
				}

				message := fmt.Sprintf("Message %d from client %d", j, clientID)
				err = conn.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					log.Printf("Error sending message: %v", err)
					conn.Close()
					time.Sleep(time.Second) // Wait before retrying
					continue
				}

				_, _, err = conn.ReadMessage()
				if err != nil {
					log.Printf("Error reading message: %v", err)
				}

				conn.Close()
			}
		}(i)
	}

	wg.Wait()
	return time.Since(start)
}

func benchmarkPolling(numClients int) time.Duration {
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			lastID := 0

			for j := 0; j < messagesPerClient; j++ {
				// Post a message
				message := fmt.Sprintf("Message %d from client %d", j, clientID)
				postMessage(message)

				// Poll for new messages with a timeout
				ctx, cancel := context.WithTimeout(context.Background(), pollTimeout)
				newMessages := pollMessages(ctx, lastID)
				cancel()

				if len(newMessages) > 0 {
					lastID = newMessages[len(newMessages)-1].ID
				}
			}
		}(i)
	}

	wg.Wait()
	return time.Since(start)
}

func postMessage(text string) {
	message := map[string]string{"text": text}
	jsonData, _ := json.Marshal(message)
	
	ctx, cancel := context.WithTimeout(context.Background(), pollTimeout)
	defer cancel()
	
	req, _ := http.NewRequestWithContext(ctx, "POST", messageURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error posting message: %v", err)
		return
	}
	defer resp.Body.Close()
}

func pollMessages(ctx context.Context, lastID int) []Message {
	url := fmt.Sprintf("%s?lastId=%d", pollURL, lastID)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Printf("Polling request timed out")
		} else {
			log.Printf("Error polling messages: %v", err)
		}
		return nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var messages []Message
	json.Unmarshal(body, &messages)
	return messages
}

func main() {
    fmt.Println("Benchmarking WebSocket and Polling servers")
    fmt.Println("==========================================")

    // Create and open a CSV file
    file, err := os.Create("benchmark_results.csv")
    if err != nil {
        log.Fatal("Cannot create file", err)
    }
    defer file.Close()

    // Write CSV header
    file.WriteString("Clients,WebSocket Total (ms),WebSocket Latency (ms),Polling Total (ms),Polling Latency (ms),Percent Difference (%)\n")

    var lastWSLatency, lastPollLatency time.Duration

    for numClients := clientStep; numClients <= maxClients; numClients += clientStep {
        wsDuration := benchmarkWebSocket(numClients)
        wsLatency := wsDuration / time.Duration(numClients*messagesPerClient)

        pollDuration := benchmarkPolling(numClients)
        pollLatency := pollDuration / time.Duration(numClients*messagesPerClient)

        fmt.Printf("Clients: %d\n", numClients)
        fmt.Printf("WebSocket - Total: %v, Latency: %v\n", wsDuration, wsLatency)
        fmt.Printf("Polling   - Total: %v, Latency: %v\n", pollDuration, pollLatency)

        // Calculate percent difference
        percentDiff := (float64(pollLatency) - float64(wsLatency)) / float64(wsLatency) * 100
        fmt.Printf("Percent difference: %.2f%% (positive means polling is slower)\n", percentDiff)

        // Write results to CSV
        csvLine := fmt.Sprintf("%d,%d,%d,%d,%d,%.2f\n",
            numClients,
            wsDuration.Milliseconds(),
            wsLatency.Milliseconds(),
            pollDuration.Milliseconds(),
            pollLatency.Milliseconds(),
            percentDiff)
        file.WriteString(csvLine)

        if numClients > clientStep {
            wsLatencyChange := (wsLatency - lastWSLatency) / lastWSLatency * 100
            pollLatencyChange := (pollLatency - lastPollLatency) / lastPollLatency * 100

            if wsLatencyChange > 10 && pollLatencyChange > 10 {
                fmt.Println("Both latencies increased by more than 10%. Stopping benchmark.")
                break
            }
        }

        lastWSLatency = wsLatency
        lastPollLatency = pollLatency

        fmt.Println("------------------------------------------")

        // Optional: add a small delay between tests to allow servers to stabilize
        time.Sleep(1 * time.Second)
    }
}