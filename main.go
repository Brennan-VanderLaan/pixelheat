package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

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

func getFileStatus(filename string) string {
	cmd := exec.Command("git", "status", "--short", filename)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
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
