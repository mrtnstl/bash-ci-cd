package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestGetHealthHandler(t *testing.T) {
	var wg sync.WaitGroup
	shutdownChan := make(chan bool)
	app := NewApplication(Config{
		Addr: ":8080",
	},
		&wg,
		&shutdownChan,
	)
	server := httptest.NewServer(http.HandlerFunc(app.getHealthHandler))

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code is 200 but got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	var parsedRespBody map[string]any
	err = json.Unmarshal(b, &parsedRespBody)
	if err != nil {
		t.Error(err)
	}

	// expected to contain "is_alive" "uptime" "last_wf_start" "last_wf_finish"
	if _, ok := parsedRespBody["is_alive"]; !ok {
		t.Errorf("expected the response body to contain %s but it didn't", "is_alive")
	}
	if _, ok := parsedRespBody["uptime"]; !ok {
		t.Errorf("expected the response body to contain %s but it didn't", "uptime")
	}
	if _, ok := parsedRespBody["last_wf_start"]; !ok {
		t.Errorf("expected the response body to contain %s but it didn't", "last_wf_start")
	}
	if _, ok := parsedRespBody["last_wf_finish"]; !ok {
		t.Errorf("expected the response body to contain %s but it didn't", "last_wf_finish")
	}
}

func TestGetStatsPaginatedHandler(t *testing.T) {
	var wg sync.WaitGroup
	shutdownChan := make(chan bool)
	app := NewApplication(Config{
		Addr: ":8080",
	},
		&wg,
		&shutdownChan,
	)
	server := httptest.NewServer(http.HandlerFunc(app.getStatsPaginatedHandler))

	// case: no query params
	resp, err := http.Get(server.URL + "?limit=&page=")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200 but got %d", resp.StatusCode)
	}

	// case: with valid query params

	resp, err = http.Get(server.URL + "?limit=20&page=1")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200 but got %d", resp.StatusCode)
	}

	// TODO: validate body

	// case: with invalid query params
	resp, err = http.Get(server.URL + "?limit=invalid&page=2")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("expected status code 400 but got %d", resp.StatusCode)
	}

	resp, err = http.Get(server.URL + "?limit=&page=invalid")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("expected status code 400 but got %d", resp.StatusCode)
	}

}

func TestTriggerCICDWorkflowHandler(t *testing.T) {
	// TODO
}

func TestWildcardRouteHandler(t *testing.T) {
	var wg sync.WaitGroup
	shutdownChan := make(chan bool)
	app := NewApplication(Config{
		Addr: ":8080",
	},
		&wg,
		&shutdownChan,
	)
	server := httptest.NewServer(http.HandlerFunc(app.wildcardRouteHandler))

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("expected status code 404 but got %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if string(b) != "Not Found!\n" {
		t.Errorf("expected response body to be 'Not Found!' but got %s", string(b))
	}
}
