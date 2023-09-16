package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	git "github.com/go-git/go-git/v5"
	"github.com/rivo/tview"
)

func getTokens(filename string) int {
	content, err := readFileContents(filename)
	if err != nil {
		log.Printf("Error reading file %s: %v", filename, err)
		return 0
	}
	return len(content) / 4
}

func findChildNode(root *tview.TreeNode, condition func(*tview.TreeNode) bool) *tview.TreeNode {
	for _, child := range root.GetChildren() {
		if condition(child) {
			return child
		}
	}
	return nil
}

// listFiles returns a list of all files in the specified directory
func listFiles(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var filenames []string
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, fmt.Sprintf("%s/%s", dir, file.Name()))
		} else {
			subDir := listFiles(fmt.Sprintf("%s/%s", dir, file.Name()))
			filenames = append(filenames, subDir...)
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

func ReadmeContent(fileNames []string) string {
	readmeNames := []string{"readme.md", "readme.txt", "readme"}

	for _, fileName := range fileNames {
		lowerName := strings.ToLower(fileName)
		for _, readmeName := range readmeNames {
			if lowerName == readmeName {
				content, err := ioutil.ReadFile(fileName)
				if err != nil {
					log.Println("Error reading README file: ", err)
				}
				return string(content)
			}
		}
	}

	return ""
}

func GenerateTitle(fileNames []string, path string, date time.Time) string {
	stack := &MessageStack{}

	// Insert fileNames, path, and date as system messages
	stack.insertSystemMessage("You are creating a tech title for a project based on whatever the user is doing, don't directly reference anything I just said, only respond with the title")
	stack.insertSystemMessage(fmt.Sprintf("File Names: %s", strings.Join(fileNames, ", ")))
	stack.insertSystemMessage(fmt.Sprintf("Path: %s", path))
	stack.insertSystemMessage(fmt.Sprintf("Date: %s", date.Format(time.RFC3339)))

	// Insert README content as a system message, if available
	readmeContent := ReadmeContent(fileNames)
	if readmeContent != "" {
		stack.insertSystemMessage(fmt.Sprintf("README Content: %s", readmeContent))
	}

	// Use AI to generate a title
	title, _ := getChatCompletion(stack.getAllMessages(), GetService("gpt-3.5", "gpt-3.5-turbo"))

	return title
}

func getLatestGitCommit() string {
	r, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}

	ref, err := r.Head()
	if err != nil {
		log.Fatal(err)
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatal(err)
	}

	latestCommit, _ := cIter.Next()

	return latestCommit.ID().String()[:8] + " " + latestCommit.Message
}

// DetermineColorBasedOnStatus returns the color for the file based on its status.
func DetermineColorBasedOnStatus(status string) tcell.Color {
	switch {
	case status == "Modified":
		return tcell.ColorYellow // Modified
	case status == "Unmodified":
		return tcell.ColorGreen // Tracked
	case status == "Untracked":
		return tcell.ColorGray // Untracked
	default:
		return tcell.ColorLightGray // Default color for any other status
	}
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
