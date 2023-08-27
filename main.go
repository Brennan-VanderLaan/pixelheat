package main

import (
	"os"
)

const openaiURL = "https://api.openai.com/v1/chat/completions"

var openaiAPIKey = os.Getenv("OPENAI_KEY")

// ServiceUsage tracks the number of API requests made for each service.
var serviceUsage = make(map[string]int)

func main() {

	stack := &MessageStack{}
	ui := NewUI(stack)
	ui.Draw()

}
