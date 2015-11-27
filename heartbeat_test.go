package heartbeat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type testGetHandler struct {
	JsonReturn string
	HttpError  string
	HttpCode   int
}

func (t testGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if t.HttpError != "" {
		http.Error(w, t.HttpError, t.HttpCode)
	} else {
		w.Write([]byte(t.JsonReturn))
	}
}

func TestHandler(t *testing.T) {

	var TestValues = []struct {
		RW           *httptest.ResponseRecorder
		Hash         string
		ExpectedHash string
	}{
		{
			httptest.NewRecorder(),
			"testHash",
			"testHash",
		},
		{
			httptest.NewRecorder(),
			"",
			NotAvailableMessage,
		},
	}

	for _, tv := range TestValues {
		CommitHash = tv.Hash

		handler(tv.RW, nil)

		hm := HeartbeatMessage{}
		err := json.NewDecoder(tv.RW.Body).Decode(&hm)
		if err != nil {
			t.Fatal(err)
		}

		if hm.Build != tv.ExpectedHash {
			t.Error("Wrong hash! Expected:", tv.ExpectedHash, "Got:", hm.Build)
		}
		if hm.Status != "running" {
			t.Error("Wrong status! Expected: running", "Got:", hm.Status)
		}
	}
}

func TestGet(t *testing.T) {

	var TestValues = []struct {
		Json      string
		HttpError string
		HttpCode  int
		RetMsg    HeartbeatMessage
		RetErr    string
	}{
		{
			`{"status":"running","build":"testHash","uptime":"5m31.5s"}`,
			"",
			http.StatusOK,
			HeartbeatMessage{"running", "testHash", "5m31.5s"},
			"",
		},
		{
			"",
			"Error",
			http.StatusBadRequest,
			HeartbeatMessage{},
			"Wrong status code: 400",
		},
		{
			"",
			"",
			http.StatusOK,
			HeartbeatMessage{},
			"Error occured unmarshalling the response",
		},
	}

	for _, tv := range TestValues {

		ts := httptest.NewServer(testGetHandler{tv.Json, tv.HttpError, tv.HttpCode})
		defer ts.Close()

		hm, err := Get(ts.URL)
		if err != nil && err.Error() != tv.RetErr {
			t.Fatal("Wrong error result! Expected:", tv.RetErr, "Got:", err)
		}

		if !reflect.DeepEqual(hm, tv.RetMsg) {
			t.Fatal("Wrong result object! Expected:", tv.RetMsg, "Got:", hm)
		}
	}
}
