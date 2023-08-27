package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

const openaiURL = "https://api.openai.com/v1/chat/completions"

var openaiAPIKey = os.Getenv("OPENAI_KEY")

// ServiceUsage tracks the number of API requests made for each service.
var serviceUsage = make(map[string]int)

type FileNode struct {
	Name   string
	Status string
	Active bool
}

var fileNodes []*FileNode

func getTokens(filename string) int {
	content, err := readFileContents(filename)
	if err != nil {
		log.Printf("Error reading file %s: %v", filename, err)
		return 0
	}
	return len(content) / 4
}

func listFiles() []string {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	var filenames []string
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}
	return filenames
}

func readFileContents(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getLatestGitCommit() string {
	cmd := exec.Command("git", "log", "-1", "--pretty=%h %B")
	output, err := cmd.Output()
	if err != nil {
		return "Error fetching commit: " + err.Error()
	}
	return string(output)
}

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

// DetermineColorBasedOnStatus returns the color for the file based on its status.
func DetermineColorBasedOnStatus(status string) tcell.Color {
	switch {
	case status == "":
		return tcell.ColorGreen // Tracked
	case strings.HasPrefix(status, "M"):
		return tcell.ColorYellow // Modified
	case strings.HasPrefix(status, "??"):
		return tcell.ColorRed // Untracked
	default:
		return tcell.ColorWhite // Default color for any other status
	}
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {

	stack := &MessageStack{}
	ui := NewUI(stack)
	ui.Draw()

}
