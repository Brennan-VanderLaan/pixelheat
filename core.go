package main

import (
	"sync"
)

type Core struct {
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
		stack:           &MessageStack{},
		activeFiles:     []*FileNode{},
		activeAIAgents:  []*AIAgentNode{},
		backendServices: make(map[string]*Service),
		serviceUsage:    make(map[string]int),
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
