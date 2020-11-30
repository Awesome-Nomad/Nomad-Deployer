package config

import (
	"fmt"
	"os"
	"strings"
)

func HTTPURLEnhancer(address string) string {
	if strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://") {
		return address
	}
	return fmt.Sprintf("http://%s", address)
}

func FileExists(fileName string) bool {
	if _, err := os.OpenFile(fileName, os.O_RDONLY, 0644); !os.IsNotExist(err) {
		return true
	}
	return false
}
