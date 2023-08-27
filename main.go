package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const openaiURL = "https://api.openai.com/v1/chat/completions"

var openaiAPIKey = os.Getenv("OPENAI_KEY")

// ServiceUsage tracks the number of API requests made for each service.
var serviceUsage = make(map[string]int)

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

func updateBackendServicesDisplay() string {
	var servicesStr string
	for _, model := range models {
		for _, service := range model.Services {
			count := serviceUsage[service.ModelName]
			cost := float64(count) * (service.InputCost + service.OutputCost) // Assuming cost is per request
			if cost > 0 || count > 0 {
				servicesStr += fmt.Sprintf("%s (API Requests: %d, Cost: $%.2f)   ", service.ModelName, count, cost)
			}
		}
	}
	return servicesStr
}

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

	choice, _ := choices[0].(map[string]interface{})
	checkError(err, "Unexpected format: choice[0] not a map")

	message, _ := choice["message"].(map[string]interface{})
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

	// Increment the service usage count
	serviceUsage[service.ModelName]++

	return content, cost
}

func main() {
	app := tview.NewApplication()

	stack := &MessageStack{}

	// Title bar
	titleBar := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("[titiw] Tape It Till It Works")
	titleBar.SetBorderPadding(1, 1, 2, 2) // Adjust padding as needed

	// Create panes with borders
	gitCommit := tview.NewTextView()
	gitCommit.SetText("Git Commit: ...").SetBorder(true).SetTitle("Git Commit")
	trackedFiles := tview.NewTreeView()

	// Create a root node
	root := tview.NewTreeNode("Project Files").SetColor(tcell.ColorWhite)
	trackedFiles.SetRoot(root).SetCurrentNode(root)

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
	chatTracking.SetScrollable(true)

	// Input field for user input
	inputField := tview.NewInputField()
	inputField.SetLabel("Enter Message: ").SetText("").SetDisabled(false)

	// Layout

	// row int, column int, rowSpan int, colSpan int, minGridHeight int, minGridWidth int,
	grid := tview.NewGrid().
		SetRows(3, 4, 0, 2).                        // Rows remain unchanged
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

			inputField.SetText("...")
			inputField.SetDisabled(true)

			// Display user's message in chatTracking
			chatTracking.SetText(chatTracking.GetText(true) + "\n[::b]User::[-] " + userMessage)

			// Use a goroutine to make the API call asynchronously
			go func() {
				// Send user's message to the API and get the response
				response, _ := getChatCompletion(stack.getAllMessages(), *GetService("gpt-4", "gpt-4"))

				stack.insertAssistantMessage(response)

				// Update the UI in the main goroutine
				app.QueueUpdateDraw(func() {
					// Display API's response in chatTracking
					chatTracking.SetText(chatTracking.GetText(true) + "\n[::b]Assistant::[-] " + response)

					// Clear the inputField and enable it
					inputField.SetText("")
					inputField.SetDisabled(false)

					// Scroll to the end of chatTracking after adding a new message
					chatTracking.ScrollToEnd()

					// Update the backend services display after making a request
					backendServices.SetText(updateBackendServicesDisplay())
				})
			}()
		}
	})

	// Populate the TreeView with files and their status
	for _, file := range listFiles() {
		status := getFileStatus(file)
		color := DetermineColorBasedOnStatus(status)
		if status != "" {
			node := tview.NewTreeNode(file + " (" + status + ")").SetColor(color)
			root.AddChild(node)
		} else {
			node := tview.NewTreeNode(file).SetColor(color)
			root.AddChild(node)
		}
	}

	// Set the latest git commit message
	gitCommit.SetText(getLatestGitCommit())

	// Update the backend services display
	backendServices.SetText(updateBackendServicesDisplay())

	// Create a slice of focusable primitives.
	focusablePrimitives := []tview.Primitive{inputField, trackedFiles}

	// Create a variable to keep track of the currently focused primitive.
	var currentFocus int

	// Set the initial focus.
	app.SetFocus(focusablePrimitives[currentFocus])

	// Capture user input to switch focus.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Capture the Tab key to switch focus.
		if event.Key() == tcell.KeyTab {
			// Increment the current focus index, wrapping around if necessary.
			currentFocus = (currentFocus + 1) % len(focusablePrimitives)
			// Set the new focus.
			app.SetFocus(focusablePrimitives[currentFocus])
			return nil // Don't propagate the handled event.
		}
		// Propagate all other events.
		return event
	})

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
