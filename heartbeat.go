package heartbeat

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const NotAvailableMessage = "Not available"

var CommitHash string
var StartTime time.Time

type HeartbeatMessage struct {
	Status string `json:"status"`
	Build  string `json:"build"`
	Uptime string `json:"uptime"`
}

func init() {
	StartTime = time.Now()
}

func handler(rw http.ResponseWriter, r *http.Request) {
	hash := CommitHash
	if hash == "" {
		hash = NotAvailableMessage
	}
	uptime := time.Now().Sub(StartTime).String()
	err := json.NewEncoder(rw).Encode(HeartbeatMessage{"running", hash, uptime})
	if err != nil {
		log.Fatalf("Failed to write heartbeat message. Reason: %s", err.Error())
	}
}

func RunHeartbeatService(address string) {
	http.HandleFunc("/heartbeat", handler)
	http.ListenAndServe(address, nil)
}
