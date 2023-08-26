package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/rivo/tview"
)

const openaiURL = "https://api.openai.com/v1/chat/completions"

var openaiAPIKey = os.Getenv("OPENAI_KEY")

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func getChatCompletion(messages []Message) string {
	data := map[string]interface{}{
		"model":       "gpt-4",
		"messages":    messages,
		"temperature": 0.7,
	}

	reqBody, err := json.Marshal(data)
	checkError(err, "Error encoding data")

	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(reqBody))
	checkError(err, "Error creating request")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiAPIKey)

	resp, err := http.DefaultClient.Do(req)
	checkError(err, "Error making request")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Received %d status from API: %s", resp.StatusCode, bodyBytes)
	}

	var result map[string]interface{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	checkError(err, "Error reading response body")

	err = json.Unmarshal(bodyBytes, &result)
	checkError(err, "Error decoding response")

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		log.Fatal("Unexpected format: 'choices' missing or not an array")
	}

	choice, ok := choices[0].(map[string]interface{})
	checkError(err, "Unexpected format: choice[0] not a map")

	message, ok := choice["message"].(map[string]interface{})
	checkError(err, "Unexpected format: 'message' not a map")

	content, ok := message["content"].(string)
	if !ok {
		log.Fatal("Unexpected format: 'content' not a string")
	}

	return content
}

func main() {
	app := tview.NewApplication()

	// Title bar
	titleBar := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Your Application Name")
	titleBar.SetBorderPadding(1, 1, 2, 2) // Adjust padding as needed

	// Create panes with borders
	gitCommit := tview.NewTextView().SetText("Git Commit: ...").SetBorder(true).SetTitle("Git Commit")
	trackedFiles := tview.NewList().ShowSecondaryText(true).SetBorder(true).SetTitle("Tracked Files")

	// Backend Services with API Requests as secondary text
	backendServices := tview.NewList().ShowSecondaryText(true)

	// Sample data: You can replace this with actual data
	services := []string{"Service A", "Service B", "Service C"}
	apiRequestCounts := []int{120, 45, 89} // Sample API request counts for each service
	for i, service := range services {
		backendServices.AddItem(service, fmt.Sprintf("API Requests: %d", apiRequestCounts[i]), 0, nil)
	}
	backendServices.SetBorder(true).SetTitle("Backend Services")

	// Layout
	grid := tview.NewGrid().
		SetRows(2, 3, 0).
		SetColumns(30, 0, 30).
		AddItem(titleBar, 0, 0, 1, 3, 0, 0, false). // Span the title bar across the entire width
		AddItem(gitCommit, 1, 0, 1, 1, 0, 0, false).
		AddItem(trackedFiles, 1, 1, 1, 2, 0, 0, false).
		AddItem(backendServices, 2, 0, 1, 3, 0, 0, false) // Span the backend services across the entire width

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}

// func main() {
// 	stack := &MessageStack{}
// 	stack.insertSystemMessage("You are a helpful assistant.")
// 	stack.insertUserMessage("Hello, I need help with my computer.")

// 	fmt.Println(getChatCompletion(stack.getAllMessages()))
// }
