package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const openaiURL = "https://api.openai.com/v1/chat/completions"

var openaiAPIKey = os.Getenv("OPENAI_KEY")

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
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Fatalf("Received %d status from API: %s", resp.StatusCode, bodyBytes)
	}

	var result map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
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
