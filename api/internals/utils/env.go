package utils

import (
	"os"
)

const (
	PORT            = "PORT"
	ALLOWED_DOMAINS = "ALLOWED_DOMAINS"
	GO_ENV          = "GO_ENV"
)


func GetEnvString(key string) string {
	return os.Getenv(key)
}
