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

type FileNode struct {
	Name   string
	Status string
	Active bool
}

var fileNodes []*FileNode

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

func updateBackendServicesDisplay() string {
	var servicesStr string
	for _, model := range models {
		for _, service := range model.Services {
			count := serviceUsage[service.ModelName]
			cost := (float64(service.InputTokens)/1000)*(service.InputCost) + (float64(service.OutputTokens)/1000)*(service.OutputCost) // Assuming cost is per 1K tokens
			if cost > 0 || count > 0 {
				servicesStr += fmt.Sprintf("%s [API Requests: %d, Cost: $%.2f (%d|%d)]   ", service.ModelName, count, cost, service.InputTokens, service.OutputTokens)
			}
		}
	}
	return servicesStr
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func getChatCompletion(messages []Message, service *Service) (string, float64) {
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

	// Calculate cost based on the number of characters used
	totalCharacters := 0
	for _, msg := range messages {
		totalCharacters += len(msg.Content)
	}
	totalTokens := float64(totalCharacters) / 4.0    // 1 token is approximately 4 characters
	cost := totalTokens / 1000 * (service.InputCost) // Assuming cost is per 1K tokens

	// Increment the service usage count and total tokens processed
	serviceUsage[service.ModelName]++
	service.InputTokens += int(totalTokens)
	service.OutputTokens += len(content) / 4

	return content, cost
}

func main() {
	app := tview.NewApplication()

	stack := &MessageStack{}

	showFormattedText := true

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

			//Clear the old messages
			stack.clearMessagesByRole("system")
			stack.insertSystemMessage("You are a meta application for helping building other applications. You are helping the user with whatever content they have selected. Follow best practices for the content you are helping with. Ask questions when neccessary.")

			// Insert contents of active files as new messages
			for _, fileNode := range fileNodes {
				if fileNode.Active {
					content, err := readFileContents(fileNode.Name)
					if err != nil {
						log.Printf("Error reading file %s: %v", fileNode.Name, err)
						continue
					}
					stack.insertSystemMessage(fmt.Sprintf("File: %s\nContent:\n%s", fileNode.Name, content))
				}
			}

			userMessage := inputField.GetText()
			stack.insertUserMessage(userMessage)

			inputField.SetText("...")
			inputField.SetDisabled(true)

			// Display user's message in chatTracking
			chatTracking.SetText(chatTracking.GetText(true) + "\n[::b]User::[-] " + userMessage)

			// Use a goroutine to make the API call asynchronously
			go func() {
				// Send user's message to the API and get the response
				response, _ := getChatCompletion(stack.getAllMessages(), GetService("gpt-4", "gpt-4"))

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

					// Set the latest git commit message
					gitCommit.SetText(getLatestGitCommit())
				})
			}()
		}
	})

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentRow, _ := chatTracking.GetScrollOffset()
		switch event.Key() {
		case tcell.KeyUp:
			if currentRow > 0 {
				chatTracking.ScrollTo(currentRow-1, 0)
			}
			return nil // Don't propagate the handled event.
		case tcell.KeyDown:
			chatTracking.ScrollTo(currentRow+1, 0)
			return nil // Don't propagate the handled event.
		}
		return event // Propagate all other events.
	})

	for _, file := range listFiles() {
		status := getFileStatus(file)
		node := tview.NewTreeNode(file)
		if status != "" {
			node.SetText(file)
		}
		node.SetColor(DetermineColorBasedOnStatus(status))
		root.AddChild(node)
		fileNode := &FileNode{Name: file, Status: status, Active: false}
		fileNodes = append(fileNodes, fileNode)
		node.SetReference(fileNode) // Store the FileNode as a reference in the TreeNode
	}

	trackedFiles.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}
		fileNode := ref.(*FileNode)
		fileNode.Active = !fileNode.Active // Toggle the active state

		// Update the display based on the active state
		if fileNode.Active {
			node.SetColor(tcell.ColorBlue)
			node.SetIndent(4)
		} else {
			color := DetermineColorBasedOnStatus(fileNode.Status)
			node.SetColor(color)
			node.SetIndent(2)
		}
	})

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
		if (event.Key() == tcell.KeyTab) && showFormattedText {
			// Increment the current focus index, wrapping around if necessary.
			currentFocus = (currentFocus + 1) % len(focusablePrimitives)
			// Set the new focus.
			app.SetFocus(focusablePrimitives[currentFocus])
			return nil // Don't propagate the handled event.
		}

		if event.Key() == tcell.KeyF1 && event.Modifiers() == tcell.ModShift {
			showFormattedText = !showFormattedText
			if showFormattedText {
				app.SetRoot(grid, true)
			} else {
				plainTextView := tview.NewTextView().SetText(stack.getPlainText())
				app.SetRoot(plainTextView, true)
			}
			return nil
		}

		// Propagate all other events.
		return event
	})

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
