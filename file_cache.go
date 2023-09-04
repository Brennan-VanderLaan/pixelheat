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

func getFileStatus(path string, filename string) string {
	filepath := path + "/" + filename
	cache, ok := fileStatusCache[filepath]
	now := time.Now()

	// Add a random duration of up to 3 seconds to the cache duration.
	randomDuration := time.Duration(rand.Float32()*2) * time.Second
	if ok && now.Sub(cache.LastCheck) < CacheDuration {
		// If the cache exists and is recent enough, return the cached status.
		return cache.Status
	}

	// If the cache does not exist or is too old, perform the check and update the cache.
	cmd := exec.Command("git", "status", "--short", filepath)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	status := strings.TrimSpace(string(output))

	fileStatusCache[filepath] = FileStatusCache{Status: status, LastCheck: now.Add(randomDuration)}
	return status
}
