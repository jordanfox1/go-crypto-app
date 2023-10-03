package test_helper

import (
	"os"
)

func SetSkipAuth() {
	if os.Getenv("SKIP_AUTH") != "" {
		return
	}
	os.Setenv("SKIP_AUTH", "1")
}

func UnsetSkipAuth() {
	if os.Getenv("SKIP_AUTH") == "" {
		return
	}
	os.Setenv("SKIP_AUTH", "")
}
