package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func MockMiddlewareChainEnd(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(200)
		w.Write([]byte{})
	}
}

func MockMiddlewareReqIP(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		
		writerContextWithIP := context.WithValue(r.Context(), REQ_IP_KEY, "192.168.0.1")
		r = r.WithContext(writerContextWithIP)
		next.ServeHTTP(w, r)
	}
}

func TestRequireHeaderSecretMiddleware(t *testing.T) {
	server = httptest.NewServer(app.RequireHeaderSecretMiddleware(MockMiddlewareChainEnd(nil)))
	client := &http.Client{}
	
	cases := []struct{
		name string
		authorizationHeaderValue string
		expectedStatusCode int
	}{
		{
			"validHeader",
			"Bearer token",
			200,
		},
		{
			"invalidHeader",
			"some-invalid-value",
			401,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T){
			req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
			
			req.Header.Set("Authorization", v.authorizationHeaderValue)

			resp, err := client.Do(req)
			
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != v.expectedStatusCode {
				t.Errorf("expected status code %d but got %d", v.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}


func TestRateLimiterMiddleware(t *testing.T) {
	server = httptest.NewServer(MockMiddlewareReqIP(app.RateLimiterMiddleware(MockMiddlewareChainEnd(nil))))
	client := &http.Client{}
	cases := []struct{
		name string
		burstCount int
		expectedStatusCode int
	}{
		{
			"allowed",
			2,
			200,
		},
		{
			"tooManyRequests",
			8,
			429,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T){
			var response *http.Response
			var respErr error

			i := 0
			for i < v.burstCount {
				i++
				req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
				
				nthResponse, err := client.Do(req)

				if i == v.burstCount {
					response = nthResponse
					respErr = err
				}
			}
			
			if respErr != nil {
				t.Error(respErr)
			}

			if response == nil {
				t.Error("response is a nil pointer!")
			}

			if response.StatusCode != v.expectedStatusCode {
				t.Errorf("expected status code %d but got %d",v.expectedStatusCode, response.StatusCode)
			}
			
		})
	}
}

func TestCheckAllowedDomainsMiddleware(t *testing.T){
	server = httptest.NewServer(MockMiddlewareReqIP(app.CheckAllowedDomainsMiddleware(MockMiddlewareChainEnd(nil))))
	
	client := &http.Client{}
	cases := []struct{
		name string
		allowedIPs string
		expectedStatusCode int
	}{
		{
			"allowed",
			"192.168.0.1;8.8.8.8",
			200,
		},
		{
			"notInAllowedEmpty",
			"",
			500,
		},
		{
			"notInAllowed",
			"8.8.8.8",
			401,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T){
			os.Unsetenv("ALLOWED_DOMAINS")
			os.Setenv("ALLOWED_DOMAINS", v.allowedIPs)
			req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

			resp, err := client.Do(req)
			
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != v.expectedStatusCode {
				t.Errorf("expected status code %d but got %d", v.expectedStatusCode, resp.StatusCode)
			}
		})
	}

	os.Unsetenv("ALLOWED_DOMAINS")
}