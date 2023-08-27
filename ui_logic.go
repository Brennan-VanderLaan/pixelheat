package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	App               *tview.Application
	Grid              *tview.Grid
	TitleBar          *tview.TextView
	GitCommit         *tview.TextView
	TrackedFiles      *tview.TreeView
	AIView            *tview.TreeView
	BackendServices   *tview.TextView
	ChatTracking      *tview.TextView
	InputField        *tview.TextArea
	FileRoot          *tview.TreeNode
	AIViewRoot        *tview.TreeNode
	CurrentFocus      int
	Primitives        []tview.Primitive
	ShowFormattedText bool
}

// NewUI creates a new UI instance
func NewUI(stack *MessageStack) *UI {
	ui := &UI{
		App:               tview.NewApplication(),
		Grid:              tview.NewGrid(),
		TitleBar:          tview.NewTextView(),
		GitCommit:         tview.NewTextView(),
		TrackedFiles:      tview.NewTreeView(),
		AIView:            tview.NewTreeView(),
		BackendServices:   tview.NewTextView(),
		ChatTracking:      tview.NewTextView(),
		InputField:        tview.NewTextArea(),
		FileRoot:          tview.NewTreeNode("Project Files"),
		AIViewRoot:        tview.NewTreeNode("AI Agents"),
		CurrentFocus:      0,
		Primitives:        []tview.Primitive{},
		ShowFormattedText: true,
	}

	ui.TitleBar.SetTextAlign(tview.AlignCenter).SetText("[pixelheat] Assisted Hallucinator")
	ui.BackendServices.SetDynamicColors(true).SetBorder(true).SetTitle(" Backend Services ")
	ui.TitleBar.SetBorderPadding(1, 1, 2, 2) // Adjust padding as needed

	// Create panes with borders
	ui.GitCommit.SetText("...").SetBorder(true).SetTitle(" Git Commit ")

	// Create a root node for each of the trees
	ui.FileRoot.SetColor(tcell.ColorWhite)
	ui.AIViewRoot.SetColor(tcell.ColorWhite)
	ui.TrackedFiles.SetRoot(ui.FileRoot).SetCurrentNode(ui.FileRoot)
	ui.AIView.SetRoot(ui.AIViewRoot).SetCurrentNode(ui.AIViewRoot)

	// Chat tracking pane
	ui.ChatTracking.SetDynamicColors(true).SetBorder(true).SetTitle(" Chat Tracking ")
	ui.ChatTracking.SetScrollable(true)

	// Input field for user input
	ui.InputField.SetBorder(true)
	ui.InputField.SetTitle(" Human Input (Shift-F2 to send) ")
	ui.InputField.SetPlaceholder("Enter your message here...\nPress Shift-F2 to send.")

	// Create a slice of focusable primitives.
	ui.AddPrimitive(ui.InputField)
	ui.AddPrimitive(ui.TrackedFiles)
	ui.AddPrimitive(ui.AIView)
	ui.SetupKeybinds(stack)

	// Layout
	ui.Grid.
		//     t git chat ai input
		SetRows(3, 4, 0, 0, 4). // Rows remain unchanged
		SetColumns(30, 0, 0).   // Adjusted column widths
		//                r  c  rs cs mh mw
		AddItem(ui.TitleBar, 0, 0, 1, 3, 0, 0, false).
		AddItem(ui.GitCommit, 1, 0, 1, 1, 0, 0, false).
		AddItem(ui.BackendServices, 1, 1, 1, 2, 0, 0, false).
		AddItem(ui.TrackedFiles, 2, 0, 1, 1, 0, 0, false).
		//                   r  c  rs cs mh mw
		AddItem(ui.AIView, 3, 0, 1, 1, 0, 0, false).
		AddItem(ui.ChatTracking, 2, 1, 3, 2, 0, 0, false).
		AddItem(ui.InputField, 5, 1, 1, 2, 0, 0, true)

	ui.Draw()

	ticker := time.NewTicker(time.Millisecond * 1000)
	go func() {
		for range ticker.C {
			ui.App.QueueUpdateDraw(func() {
				ui.Draw()
			})
		}
	}()

	if err := ui.App.SetRoot(ui.Grid, true).Run(); err != nil {
		panic(err)
	}

	return ui
}

func (ui *UI) SetupKeybinds(stack *MessageStack) {
	// Capture user input to switch focus.
	ui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Capture the Tab key to switch focus.
		if (event.Key() == tcell.KeyTab) && ui.ShowFormattedText {
			// Increment the current focus index, wrapping around if necessary.
			ui.CurrentFocus = (ui.CurrentFocus + 1) % len(ui.Primitives)
			// Set the new focus.
			ui.App.SetFocus(ui.Primitives[ui.CurrentFocus])
			return nil // Don't propagate the handled event.
		}

		if event.Key() == tcell.KeyF1 && event.Modifiers() == tcell.ModShift {
			ui.ShowFormattedText = !ui.ShowFormattedText
			if ui.ShowFormattedText {
				ui.App.SetRoot(ui.Grid, true)
			} else {
				plainTextView := tview.NewTextView().SetText(stack.getPlainText())
				ui.App.SetRoot(plainTextView, true)
			}
			return nil
		}

		if event.Key() == tcell.KeyF2 && event.Modifiers() == tcell.ModShift {
			ui.HandleInput(stack)
		}

		// Propagate all other events.
		return event
	})

	ui.TrackedFiles.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}
		fileNode := ref.(*FileNode)
		fileNode.Active = !fileNode.Active // Toggle the active sta	te

		// If the fileNode is active, add a child node with "*ACTIVE*"
		// If the fileNode is not active, remove all child nodes (assuming it only has the "*ACTIVE*" node)
		if fileNode.Active {
			// Add "*ACTIVE*" as a child node
			activeNode := tview.NewTreeNode(fmt.Sprintf("*ACTIVE* (%d)", getTokens(fileNode.Name)))
			activeNode.SetColor(tcell.ColorBlue)
			node.AddChild(activeNode)
		} else {
			color := DetermineColorBasedOnStatus(fileNode.Status)
			node.SetColor(color)
			// Remove all child nodes
			for _, child := range node.GetChildren() {
				node.RemoveChild(child)
			}
		}
	})
}

// AddPrimitive adds a primitive to the UI
func (ui *UI) AddPrimitive(p tview.Primitive) {
	ui.Primitives = append(ui.Primitives, p)
	ui.App.SetFocus(ui.Primitives[ui.CurrentFocus])
}

// SwitchFocus handles focus switching
func (ui *UI) SwitchFocus() {
	ui.CurrentFocus = (ui.CurrentFocus + 1) % len(ui.Primitives)
	ui.App.SetFocus(ui.Primitives[ui.CurrentFocus])
}

// HandleInput handles user input
func (ui *UI) HandleInput(stack *MessageStack) {

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

	userMessage := ui.InputField.GetText()
	stack.insertUserMessage(userMessage)

	ui.InputField.SetText("", false)
	ui.InputField.SetDisabled(true)

	// Display user's message in chatTracking
	ui.ChatTracking.SetText(ui.ChatTracking.GetText(true) + "\n[::b]User::[-] " + userMessage)

	// Use a goroutine to make the API call asynchronously
	go func() {
		// Send user's message to the API and get the response
		response, _ := getChatCompletion(stack.getAllMessages(), GetService("gpt-4", "gpt-4"))

		stack.insertAssistantMessage(response)

		// Update the UI in the main goroutine
		ui.App.QueueUpdateDraw(func() {
			// Display API's response in chatTracking
			ui.ChatTracking.SetText(ui.ChatTracking.GetText(true) + "\n[::b]Assistant::[-] " + response)

			// Clear the inputField and enable it
			ui.InputField.SetText("", false)
			ui.InputField.SetDisabled(false)

			// Scroll to the end of chatTracking after adding a new message
			ui.ChatTracking.ScrollToEnd()

			// Update the backend services display after making a request
			ui.UpdateBackendServices()

			// Set the latest git commit message
			ui.GitCommit.SetText(getLatestGitCommit())
		})
	}()

}

// Draw draws the UI to the screen
func (ui *UI) Draw() {
	ui.UpdateAIView()
	ui.UpdateTrackedFiles()
	ui.UpdateGitCommit()
	ui.UpdateBackendServices()
}

// UpdateGitCommit updates the git commit text view
func (ui *UI) UpdateGitCommit() {

	// Set the latest git commit message
	ui.GitCommit.SetText(getLatestGitCommit())
}

// UpdateBackendServices updates the backend services text view
func (ui *UI) UpdateBackendServices() {
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

	ui.BackendServices.SetText(servicesStr)
}

// UpdateChatTracking updates the chat tracking text view
func (ui *UI) UpdateChatTracking() {
	// ...
}

// UpdateInputField updates the input field
func (ui *UI) UpdateInputField() {
	// ...
}

// UpdateTrackedFiles updates the tracked files tree view
func (ui *UI) UpdateTrackedFiles() {
	for _, file := range listFiles() {

		foundMatch := false
		matchingNode := &tview.TreeNode{}
		for _, existingNode := range ui.FileRoot.GetChildren() {
			if existingNode.GetText() == file {
				foundMatch = true
				matchingNode = existingNode
				break
			} else {
				continue
			}
		}

		if !foundMatch {
			// We haven't seen this file yet, make a new node
			status := getFileStatus(file)
			node := tview.NewTreeNode(file)
			if status != "" {
				node.SetText(file)
			}
			node.SetColor(DetermineColorBasedOnStatus(status))
			ui.FileRoot.AddChild(node)
			fileNode := &FileNode{Name: file, Status: status, Active: false}
			fileNodes = append(fileNodes, fileNode)
			node.SetReference(fileNode) // Store the FileNode as a reference in the TreeNode
		} else {
			status := getFileStatus(file)
			matchingNode.SetText(file)
			matchingNode.SetColor(DetermineColorBasedOnStatus(status))

		}
	}

	for _, existingNode := range ui.FileRoot.GetChildren() {
		foundMatch := false
		for _, file := range listFiles() {
			if existingNode.GetText() == file {
				foundMatch = true
				break
			} else {
				continue
			}
		}

		if !foundMatch {
			// We have a node that doesn't match any files, remove it
			ui.FileRoot.RemoveChild(existingNode)
		}
	}

}

// UpdateAIView updates the AI view tree view
func (ui *UI) UpdateAIView() {

	for _, agent := range aiAgents {

		foundMatch := false
		matchingAgent := &tview.TreeNode{}
		for _, existingNode := range ui.AIViewRoot.GetChildren() {
			if existingNode.GetText() == agent.Name {
				foundMatch = true
				matchingAgent = existingNode
				break
			} else {
				continue
			}
		}

		if !foundMatch {
			// We haven't seen this agent yet, make a new node
			agentNode := tview.NewTreeNode(agent.Name).SetColor(tcell.ColorWhite)
			ui.AIViewRoot.AddChild(agentNode)
			matchingAgent = agentNode

		}

		matchingAgent.ClearChildren()
		for _, service := range agent.Services {
			serviceNode := tview.NewTreeNode(service.ModelName).SetColor(tcell.ColorCadetBlue)
			matchingAgent.AddChild(serviceNode)
		}
	}

}
