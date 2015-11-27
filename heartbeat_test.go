package heartbeat

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testHandler struct {
}

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(w, r)
}

func init() {
	StartTime = time.Date(2015, 11, 27, 11, 47, 00, 00, time.UTC)
}

func TestHandler(t *testing.T) {

	var TestValues = []struct {
		Hash       string
		HashResult string
	}{
		{"", NotAvailableMessage},
		{"testHash", "testHash"},
	}

	for _, tv := range TestValues {
		CommitHash = tv.Hash

		ts := httptest.NewServer(testHandler{})
		defer ts.Close()

		res, err := http.Get(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		respJson, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		var hm HeartbeatMessage
		err = json.Unmarshal(respJson, &hm)
		if err != nil {
			t.Fatal(err)
		}

		if hm.Status != "running" {
			t.Fatal(errors.New("The server should running"))
		}

		if hm.Build != tv.HashResult {
			t.Fatal(errors.New("Wrong commit hash"))
		}

		uptime, err := time.ParseDuration(hm.Uptime)
		if err != nil {
			t.Fatal(err)
		}

		if uptime > time.Since(StartTime) {
			t.Fatal(errors.New("Wrong uptime"))
		}
	}
}

func TestHGet(t *testing.T) {

}
