package main

import (
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

type FileNode struct {
	Name   string
	Status string
	Active bool
}

var fileNodes []*FileNode

// FileStatusCache represents the cached file status,
// including the time of the last check.
type FileStatusCache struct {
	Status    string
	LastCheck time.Time
}

const CacheDuration = 2 * time.Second

var fileStatusCache = make(map[string]FileStatusCache)

func getFileStatus(filename string) string {
	cache, ok := fileStatusCache[filename]
	now := time.Now()

	// Add a random duration of up to 3 seconds to the cache duration.
	randomDuration := time.Duration(rand.Float32()*2) * time.Second
	if ok && now.Sub(cache.LastCheck) < CacheDuration {
		// If the cache exists and is recent enough, return the cached status.
		return cache.Status
	}

	// If the cache does not exist or is too old, perform the check and update the cache.
	cmd := exec.Command("git", "status", "--short", filename)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	status := strings.TrimSpace(string(output))

	fileStatusCache[filename] = FileStatusCache{Status: status, LastCheck: now.Add(randomDuration)}
	return status
}
