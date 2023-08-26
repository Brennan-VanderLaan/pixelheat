package main

type MessageStack struct {
	messages []Message
}

func (ms *MessageStack) setMessages(messages []Message) {
	ms.messages = messages
}

func (ms *MessageStack) clearMessages() {
	ms.messages = []Message{}
}

func (ms *MessageStack) insertUserMessage(content string) {
	ms.messages = append(ms.messages, Message{Role: "user", Content: content})
}

func (ms *MessageStack) insertSystemMessage(content string) {
	ms.messages = append(ms.messages, Message{Role: "system", Content: content})
}

func (ms *MessageStack) insertAssistantMessage(content string) {
	ms.messages = append(ms.messages, Message{Role: "assistant", Content: content})
}

func (ms *MessageStack) getAllUserMessages() []Message {
	return ms.getMessagesByRole("user")
}

func (ms *MessageStack) getAllAssistantMessages() []Message {
	return ms.getMessagesByRole("assistant")
}

func (ms *MessageStack) getAllSystemMessages() []Message {
	return ms.getMessagesByRole("system")
}

func (ms *MessageStack) getAllMessages() []Message {
	return ms.messages
}

func (ms *MessageStack) getMessagesByRole(role string) []Message {
	var filteredMessages []Message
	for _, msg := range ms.messages {
		if msg.Role == role {
			filteredMessages = append(filteredMessages, msg)
		}
	}
	return filteredMessages
}
