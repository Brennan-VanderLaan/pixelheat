package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gdamore/tcell/v2"
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
func getChatCompletion(messages []Message, service Service) (string, float64) {
	data := map[string]interface{}{
		"model":       service.ModelName,
		"messages":    messages,
		"temperature": .7,
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

	// Calculate cost based on the number of tokens used (this is a simplified example)
	// You might need to adjust this based on the exact response format from OpenAI
	usage, ok := result["usage"].(map[string]interface{})
	if !ok {
		log.Fatal("Unexpected format: 'usage' not a map")
	}
	totalTokens, ok := usage["total_tokens"].(float64)
	if !ok {
		log.Fatal("Unexpected format: 'total_tokens' not a number")
	}
	cost := totalTokens / 1000 * (service.InputCost + service.OutputCost) // Assuming cost is per 1K tokens

	return content, cost
}

func main() {
	app := tview.NewApplication()

	stack := &MessageStack{}

	// Title bar
	titleBar := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Your Application Name")
	titleBar.SetBorderPadding(1, 1, 2, 2) // Adjust padding as needed

	// Create panes with borders
	gitCommit := tview.NewTextView().SetText("Git Commit: ...").SetBorder(true).SetTitle("Git Commit")
	trackedFiles := tview.NewList().ShowSecondaryText(true).SetBorder(true).SetTitle("Tracked Files")

	// Backend Services with API Requests displayed side by side
	backendServices := tview.NewTextView().SetDynamicColors(true)
	// Sample data
	services := []string{"Service A", "Service B", "Service C"}
	apiRequestCounts := []int{120, 45, 89}
	var servicesStr string
	for i, service := range services {
		servicesStr += fmt.Sprintf("%s (API Requests: %d)   ", service, apiRequestCounts[i])
	}
	backendServices.SetText(servicesStr).SetBorder(true).SetTitle("Backend Services")

	// Chat tracking pane
	chatTracking := tview.NewTextView()
	chatTracking.SetDynamicColors(true).SetBorder(true).SetTitle("Chat Tracking")

	// Input field for user input
	inputField := tview.NewInputField()
	inputField.SetLabel("Enter Message: ").SetText("").SetDisabled(false)

	// Layout
	// row int, column int, rowSpan int, colSpan int, minGridHeight int, minGridWidth int,
	grid := tview.NewGrid().
		SetRows(2, 3, 0, 2).                        // Rows remain unchanged
		SetColumns(30, 30, 0).                      // Adjusted column widths
		AddItem(titleBar, 0, 0, 1, 3, 0, 0, false). // Span the title bar across the entire width
		AddItem(gitCommit, 1, 0, 1, 1, 0, 0, false).
		AddItem(backendServices, 1, 1, 1, 2, 0, 0, false). // Span backendServices across two columns
		AddItem(trackedFiles, 2, 0, 1, 1, 0, 0, false).
		AddItem(chatTracking, 2, 1, 1, 2, 0, 0, false). // Span chatTracking across two columns
		AddItem(inputField, 3, 1, 1, 2, 0, 0, true)     // Span the input field across two columns

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			userMessage := inputField.GetText()
			stack.insertUserMessage(userMessage)

			// Display user's message in chatTracking
			chatTracking.SetText(chatTracking.GetText(true) + "\n[::b]User::[-] " + userMessage)

			// Send user's message to the API and get the response
			response, _ := getChatCompletion(stack.getAllMessages(), *GetService("gpt-4", "gpt-4"))

			stack.insertAssistantMessage(response)

			// Display API's response in chatTracking
			chatTracking.SetText(chatTracking.GetText(true) + "\n[::b]Assistant::[-] " + response)

			// Clear the inputField
			inputField.SetText("")
		}
	})

	app.SetFocus(inputField)

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
