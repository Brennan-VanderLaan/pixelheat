package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const openaiURL = "https://api.openai.com/v1/chat/completions"

var openaiAPIKey = os.Getenv("OPENAI_KEY")

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func getChatCompletion(messages []Message) string {
	reqBody, err := json.Marshal(map[string]interface{}{
		"model":       "gpt-4",
		"messages":    messages,
		"temperature": 0.7,
	})
	if err != nil {
		log.Fatalf("Error encoding data: %v", err)
	}

	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openaiAPIKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Received %d status from API: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		log.Fatalf("Unexpected format: 'choices' missing or not an array")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		log.Fatalf("Unexpected format: choice[0] not a map")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		log.Fatalf("Unexpected format: 'message' not a map")
	}

	content, ok := message["content"].(string)
	if !ok {
		log.Fatalf("Unexpected format: 'content' not a string")
	}

	return content
}

func main() {
	// Example usage of getChatCompletion
	messages := []Message{
		{
			Role:    "system",
			Content: "You are helping gather some information from the user.",
		},
	}

	response := getChatCompletion(messages)
	fmt.Println(response)
}
