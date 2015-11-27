package heartbeat

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type testHandler struct {
}

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(w, r)
}

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
			t.Error(errors.New("The server should running"))
		}

		if hm.Build != tv.HashResult {
			t.Error(errors.New("Wrong commit hash"))
		}

		uptime, err := time.ParseDuration(hm.Uptime)
		if err != nil {
			t.Error(err)
		}

		if uptime > time.Since(StartTime) {
			t.Error(errors.New("Wrong uptime"))
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
