package util

import (
	"os/exec"
)

// cache commandExists
var commandExistsCache = make(map[string]bool)

func CommandExists(cmd string) bool {
	// check the cache first
	if exists, ok := commandExistsCache[cmd]; ok {
		return exists
	}

	_, err := exec.LookPath(cmd)
	exists := err == nil

	// save to cache
	commandExistsCache[cmd] = exists

	return exists
}
