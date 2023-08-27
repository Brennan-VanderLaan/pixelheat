package main

// MessageStack represents a stack of messages.
type MessageStack struct {
	messages []Message
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// setMessages sets the messages in the stack.
func (ms *MessageStack) setMessages(messages []Message) {
	ms.messages = messages
}

// clearMessages clears all messages from the stack.
func (ms *MessageStack) clearMessages() {
	ms.messages = []Message{}
}

// insertUserMessage inserts a user message into the stack.
func (ms *MessageStack) insertUserMessage(content string) {
	ms.messages = append(ms.messages, Message{Role: "user", Content: content})
}

// insertSystemMessage inserts a system message into the stack.
func (ms *MessageStack) insertSystemMessage(content string) {
	ms.messages = append(ms.messages, Message{Role: "system", Content: content})
}

// insertAssistantMessage inserts an assistant message into the stack.
func (ms *MessageStack) insertAssistantMessage(content string) {
	ms.messages = append(ms.messages, Message{Role: "assistant", Content: content})
}

// getAllUserMessages returns all user messages in the stack.
func (ms *MessageStack) getAllUserMessages() []Message {
	return ms.getMessagesByRole("user")
}

// getAllAssistantMessages returns all assistant messages in the stack.
func (ms *MessageStack) getAllAssistantMessages() []Message {
	return ms.getMessagesByRole("assistant")
}

// getAllSystemMessages returns all system messages in the stack.
func (ms *MessageStack) getAllSystemMessages() []Message {
	return ms.getMessagesByRole("system")
}

// getAllMessages returns all messages in the stack.
func (ms *MessageStack) getAllMessages() []Message {
	return ms.messages
}

// getFormattedText returns the formatted text representation of the messages.
func (ms *MessageStack) getFormattedText() string {
	var formattedText string
	for _, msg := range ms.messages {
		switch msg.Role {
		case "user":
			formattedText += "[::b]User::[-] " + msg.Content + "\n"
		case "assistant":
			formattedText += "[::b]Assistant::[-] " + msg.Content + "\n"
		case "system":
			formattedText += "[::b]System::[-] " + msg.Content + "\n"
		}
	}
	return formattedText
}

// getPlainText returns the plain text representation of the messages.
func (ms *MessageStack) getPlainText() string {
	var plainText string
	for _, msg := range ms.messages {
		switch msg.Role {
		case "user":
			plainText += "User: " + msg.Content + "\n"
		case "assistant":
			plainText += "Assistant: " + msg.Content + "\n"
		case "system":
			//plainText += "System: " + msg.Content + "\n"
		}
	}
	return plainText
}

// clearMessagesByRole returns all messages in the stack with the specified role.
func (ms *MessageStack) clearMessagesByRole(role string) []Message {
	var filteredMessages []Message
	for _, msg := range ms.messages {
		if msg.Role != role {
			filteredMessages = append(filteredMessages, msg)
		}
	}

	ms.messages = filteredMessages
	return filteredMessages
}

// getMessagesByRole returns all messages in the stack with the specified role.
func (ms *MessageStack) getMessagesByRole(role string) []Message {
	var filteredMessages []Message
	for _, msg := range ms.messages {
		if msg.Role == role {
			filteredMessages = append(filteredMessages, msg)
		}
	}
	return filteredMessages
}
