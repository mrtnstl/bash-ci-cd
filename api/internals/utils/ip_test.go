package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// only happy paths are tested here
func TestGetIP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		reqIP, err := GetIP(r)
		if err != nil {
			w.WriteHeader(400)
			return	
		}
		w.Write([]byte(reqIP))
	}))

	client := &http.Client{}

	cases := []struct{
		name string
		xForwardedForHeader string
		expectedIP string
	}{
		{
			"ipOnXForwHeader",
			"8.8.8.8,127.0.0.1",
			"127.0.0.1",
		},
		{
			"ipOnRemoteAddr",
			"",
			"127.0.0.1",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T){
			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			req.Header.Set("X-Forwarder-For", v.xForwardedForHeader)
		
			res, err := client.Do(req)
			if err != nil {
				t.Error(err)
			}
			defer res.Body.Close()

			b, err := io.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
			}

			if string(b) != v.expectedIP {
				t.Errorf("expected ip to be %v, got %v instead", v.expectedIP, string(b))
			}
		})
	}
}