package main

import (
	"fmt"
	"log"
)

type AIAgent struct {
	Name             string
	Directive        string
	PrefferedService *Service
	Services         []*Service
}

// AIAgent handleInput method takes a string and a stack of messages and returns a string
func (a *AIAgent) HandleInput(input string, stack *MessageStack, core *Core) string {

	//Clear the old messages
	stack.clearMessagesByRole("system")
	stack.insertSystemMessage(a.Directive)

	// Insert contents of active files as new messages
	for _, fileNode := range core.GetActiveFiles() {
		if fileNode.Active {
			content, err := readFileContents(fileNode.Name)
			if err != nil {
				log.Printf("Error reading file %s: %v", fileNode.Name, err)
				continue
			}
			stack.insertSystemMessage(fmt.Sprintf("File: %s\nContent:\n%s", fileNode.Name, content))
		}
	}

	stack.insertUserMessage(input)

	// Send user's message to the API and get the response
	response, _ := getChatCompletion(stack.getAllMessages(), GetService("gpt-4", "gpt-4"))
	stack.insertAssistantMessage(response)

	return response
}

var aiAgents = []*AIAgent{
	{
		Name:      "PixelHeat",
		Directive: "You are a meta application for helping building other applications. You are helping the user with whatever content they have selected. Follow best practices for the content you are helping with. Ask questions when neccessary.",
		Services: []*Service{
			GetService("gpt-4", "gpt-4"),
		},
	},
	{
		Name:      "PixelHeat (Pirate)",
		Directive: "You are a meta application for helping building other applications. You are helping the user with whatever content they have selected. Follow best practices for the content you are helping with. Ask questions when neccessary. You only respond as a helpful pirate, do everything you can to stay in character. Never break character.",
		Services: []*Service{
			GetService("gpt-4", "gpt-4"),
		},
	},
	{
		Name:      "Chat Assistant (smart)",
		Directive: "You are a helpful chat assistant. You can help with nearly anything, if you are unsure of the validity of an answer state as much.",
		Services: []*Service{
			GetService("gpt-4", "gpt-4"),
			GetService("gpt-3.5", "gpt-3.5-turbo"),
		},
		PrefferedService: GetService("gpt-4", "gpt-4"),
	},
	{
		Name:      "Chat Assistant (eh)",
		Directive: "You are a helpful chat assistant. You can help with nearly anything, if you are unsure of the validity of an answer state as much.",
		Services: []*Service{
			GetService("gpt-3.5", "gpt-3.5-turbo"),
		},
		PrefferedService: GetService("gpt-3.5", "gpt-3.5-turbo"),
	},
	{
		Name:      "Code Reviewer (friendly)",
		Directive: "You are a friendly code reviewer. You are not too strict, but you are not too lenient either. You are a good balance of both.",
		Services: []*Service{
			GetService("gpt-4", "gpt-4"),
		},
	},
	// ... add more agents as needed
}
