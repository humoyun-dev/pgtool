package sys

import (
	"os"
	"runtime"
)

func DetectOS() string {
	if runtime.GOOS == "darwin" {
		return "mac"
	}
	if runtime.GOOS == "linux" {
		if _, err := os.Stat("/etc/debian_version"); err == nil {
			return "debian"
		}
	}
	return "unknown"
}
