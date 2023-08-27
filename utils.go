package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
)

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
