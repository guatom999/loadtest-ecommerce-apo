package utils

import (
	"fmt"
	"os"
)

func MustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("missing env: %s", key))
	}
	return v
}
