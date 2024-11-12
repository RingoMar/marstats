package rqueue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Queue represents a simple in-memory queue with a mutex for thread safety
type USERNOTICE struct {
	ID         string
	Channel    string
	Name       string
	SubMethod  string
	SubAmount  string
	GiftAmount string
	SubPlan    string
	Created    string
}

type Queue struct {
	items []USERNOTICE
	mu    sync.Mutex
}

// Add an item to the queue
func (q *Queue) AddQueue(ID string, Channel string, Name string, SubMethod string, SubAmount string, GiftAmount string, SubPlan string, Created string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	// Build Usernotice struct
	var subUsernotice = new(USERNOTICE)
	subUsernotice.ID = ID
	subUsernotice.Channel = Channel
	subUsernotice.Name = Name
	subUsernotice.SubMethod = SubMethod
	subUsernotice.SubAmount = SubAmount
	subUsernotice.GiftAmount = GiftAmount
	subUsernotice.SubPlan = SubPlan
	subUsernotice.Created = Created

	q.items = append(q.items, *subUsernotice)
}

// Dump the contents of the queue
func (q *Queue) Dump() {
	q.mu.Lock()
	defer q.mu.Unlock()
	fmt.Println("Dumping queue contents:", q.items)
	dumpapi := uploadDump(q.items)
	if dumpapi {
		q.items = nil
	}
}

// Find the length of the queue
func (q *Queue) Length() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Monitor for inactivity and dump queue if no activity occurs within timeout
func MonitorQueue(q *Queue, timeout time.Duration) {
	timer := time.NewTimer(timeout)

	for {
		// `<-timer.C` blocks the execution until the timer signals that the specified duration has elapsed.
		<-timer.C
		if q.Length() > 0 {
			q.Dump() // Dump queue contents if there are items
		}
		timer.Reset(timeout) // Reset timer for next inactivity period
	}
}

func uploadDump(itmes []USERNOTICE) bool {
	upStreamURL := "localhost"
	if os.Getenv("GO_ENV") == "production" {
		upStreamURL = "mongoserver"
	}
	url := fmt.Sprintf("http://%s:5284/api/dump", upStreamURL)
	queueData, err := json.Marshal(itmes)
	if err != nil {
		fmt.Printf("Failed to marshal queue items: %v", err)
		return false
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(queueData))
	if err != nil {
		fmt.Printf("Failed to send request: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("Queue sent successfully.")
		return true
	} else {
		fmt.Printf("Failed to send queue: %s\n", resp.Status)
		return false
	}
}
