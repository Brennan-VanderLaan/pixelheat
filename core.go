package main

import (
	"sync"
)

type Core struct {
	projectDir       string
	stack            *MessageStack
	activeFiles      []*FileNode
	activeAIAgents   []*AIAgentNode
	backendServices  map[string]*Service
	serviceUsage     map[string]int
	commitMessage    string
	userInput        string
	assistantMessage string
	mu               sync.Mutex
}

func NewCore() *Core {
	return &Core{
		projectDir:      ".",
		stack:           &MessageStack{},
		activeFiles:     []*FileNode{},
		activeAIAgents:  []*AIAgentNode{},
		backendServices: make(map[string]*Service),
		serviceUsage:    make(map[string]int),
	}
}

func (c *Core) Update() {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Logic for updating the core

	//Update files in current project
	c.UpdateFiles()

}

// Update files in current project
func (c *Core) UpdateFiles() {
	files := listFiles(".")

	for _, fileName := range files {
		// Check if this file is already being tracked
		found := false
		for _, fileNode := range c.activeFiles {
			if fileNode.Name == fileName {
				// Update the status of the existing file node
				fileNode.Status = getFileStatus(c.projectDir, fileName)
				found = true
				break
			}
		}

		if !found {
			// Add a new file node for this file
			status := getFileStatus(c.projectDir, fileName)
			fileNode := &FileNode{Name: fileName, Status: status, Active: false}
			c.activeFiles = append(c.activeFiles, fileNode)
		}
	}
}

// Getters
func (c *Core) GetStack() *MessageStack {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stack
}

func (c *Core) GetActiveFiles() []*FileNode {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.activeFiles
}

func (c *Core) GetActiveAIAgents() []*AIAgentNode {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.activeAIAgents
}

func (c *Core) GetBackendServices() map[string]*Service {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.backendServices
}

func (c *Core) GetCommitMessage() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.commitMessage
}

func (c *Core) GetUserInput() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.userInput
}

func (c *Core) GetAssistantMessage() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.assistantMessage
}

// Setters/Modifiers
func (c *Core) AddActiveFile(file *FileNode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.activeFiles = append(c.activeFiles, file)
}

func (c *Core) RemoveActiveFile(file *FileNode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Logic for removing the FileNode from the slice
}

func (c *Core) AddActiveAIAgent(agent *AIAgentNode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.activeAIAgents = append(c.activeAIAgents, agent)
}

func (c *Core) RemoveActiveAIAgent(agent *AIAgentNode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Logic for removing the AIAgentNode from the slice
}

func (c *Core) UpdateBackendService(service *Service) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.backendServices[service.ModelName] = service
}

func (c *Core) UpdateServiceUsage(modelName string, usage int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.serviceUsage[modelName] = usage
}

func (c *Core) SetCommitMessage(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.commitMessage = message
}

func (c *Core) SetUserInput(input string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.userInput = input
}

func (c *Core) SetAssistantMessage(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.assistantMessage = message
}
