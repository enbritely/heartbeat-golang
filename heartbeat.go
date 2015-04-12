package heartbeat

import (
	"encoding/json"
	"log"
	"net/http"
)

const NotAvailableMessage = "Not available"

var CommitHash string

type HeartbeatMessage struct {
	Status string `json:"status"`
	Build  string `json:"build"`
}

func handler(rw http.ResponseWriter, r *http.Request) {
	hash := CommitHash
	if hash == "" {
		hash = NotAvailableMessage
	}
	err := json.NewEncoder(rw).Encode(HeartbeatMessage{"running", hash})
	if err != nil {
		log.Fatalf("Failed to write heartbeat message. Reason: %s", err.Error())
	}
}

func RunHeartbeatService(address string) {
	http.HandleFunc("/heartbeat", handler)
	http.ListenAndServe(address, nil)
}
