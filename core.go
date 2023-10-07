package main

import (
	"sync"

	"github.com/go-git/go-git/v5"
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
	files := listFiles(c.projectDir)
	dirs := listDirs(c.projectDir)

	r, _ := git.PlainOpen(c.projectDir)
	w, _ := r.Worktree()
	status, _ := w.Status()

	for _, dirName := range dirs {
		// Check if this directory is already being tracked
		found := false
		for _, fileNode := range c.activeFiles {
			if fileNode.Name == dirName {
				// Update the status of the existing file node
				fileNode.Status = "Directory"
				found = true
				break
			}
		}

		if !found {
			// Add a new file node for this file
			fileNode := &FileNode{Name: dirName, Status: "Directory", Active: false, Directory: true}
			c.activeFiles = append(c.activeFiles, fileNode)
		}
	}

	for _, fileName := range files {
		// Check if this file is already being tracked
		found := false
		for _, fileNode := range c.activeFiles {
			if fileNode.Name == fileName {

				var untracked bool
				if status.IsUntracked(fileName) {
					untracked = true
				}

				// Update the status of the existing file node
				file_status := status.File(fileName)
				status_str := ""
				if !untracked && file_status.Staging == git.Untracked &&
					file_status.Worktree == git.Untracked {
					file_status.Staging = git.Unmodified
					file_status.Worktree = git.Unmodified
					status_str = "Unmodified"
				}
				if file_status.Worktree == git.Modified {
					status_str = "Modified"
				}

				fileNode.Status = status_str
				found = true
				break
			}
		}

		if !found {
			fileNode := &FileNode{Name: fileName, Status: "Unmodified", Active: false}
			c.activeFiles = append(c.activeFiles, fileNode)
		}
	}
}

// Handle Input
func (c *Core) HandleInput(input string) string {
	// Logic for handling input

	stack := c.GetStack()

	//Check active agents at least 0
	if len(c.GetActiveAIAgents()) == 0 {
		return "No active agents."
	}

	//Get the active agent
	agent := c.GetActiveAIAgents()[0]

	response := agent.AIAgent.HandleInput(input, stack, c)

	return response
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

	//Find the index of the agent
	index := -1
	for i, a := range c.activeAIAgents {
		if a == agent {
			index = i
			break
		}
	}

	//Remove the agent
	if index != -1 {
		c.activeAIAgents = append(c.activeAIAgents[:index], c.activeAIAgents[index+1:]...)
	}

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
