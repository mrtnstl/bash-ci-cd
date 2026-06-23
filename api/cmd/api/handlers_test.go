package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"golang.org/x/time/rate"
)

var wg sync.WaitGroup
var shutdownChan chan bool
var app *Application
var server *httptest.Server

func TestMain(m *testing.M) {
	shutdownChan = make(chan bool)
	app = NewApplication(Config{
		Addr: ":8080",
		RlLimit: rate.Limit(1),
		RlBurst: 3,
	},
		&wg,
		&shutdownChan,
	)

	m.Run()

	shutdownChan = nil
	app = nil
	server.Close()
	server = nil
}

func TestGetHealthHandler(t *testing.T) {
	server = httptest.NewServer(http.HandlerFunc(app.getHealthHandler))

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
	server = httptest.NewServer(http.HandlerFunc(app.getStatsPaginatedHandler))

	cases := []struct{
		name string
		params string
		expectedStatusCode int
	}{
		{
			"unspecifiedParams",
			"?limit=&page=",
			200,
		},
		{
			"validParams",
			"?limit=20&page=1",
			200,
		},
		{
			"invalidLimit",
			"?limit=invalid&page=2",
			400,
		},
		{
			"invalidPage",
			"?limit=&page=invalid",
			400,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T){
			resp, err := http.Get(server.URL + v.params)
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != v.expectedStatusCode {
				t.Errorf("expected status code %d but got %d",v.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}

func TestTriggerCICDWorkflowHandler(t *testing.T) {
	// TODO
}

func TestWildcardRouteHandler(t *testing.T) {
	server = httptest.NewServer(http.HandlerFunc(app.wildcardRouteHandler))

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
