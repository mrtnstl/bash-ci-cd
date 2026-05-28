package utils

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
)

func GetIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarder-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}

		return ip, nil
	}

	return "", errors.New("IP not found!")
}

func GetEnvString(key string) string {
	return os.Getenv(key)
}
