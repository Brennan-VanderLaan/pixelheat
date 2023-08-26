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

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func getChatCompletion(messages []Message) string {
	data := map[string]interface{}{
		"model":       "gpt-4",
		"messages":    messages,
		"temperature": 0.7,
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
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Received %d status from API: %s", resp.StatusCode, bodyBytes)
	}

	var result map[string]interface{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	checkError(err, "Error reading response body")

	err = json.Unmarshal(bodyBytes, &result)
	checkError(err, "Error decoding response")

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		log.Fatal("Unexpected format: 'choices' missing or not an array")
	}

	choice, ok := choices[0].(map[string]interface{})
	checkError(err, "Unexpected format: choice[0] not a map")

	message, ok := choice["message"].(map[string]interface{})
	checkError(err, "Unexpected format: 'message' not a map")

	content, ok := message["content"].(string)
	if !ok {
		log.Fatal("Unexpected format: 'content' not a string")
	}

	return content
}

func main() {
	messages := []Message{
		{
			Role:    "system",
			Content: "You are helping gather some information from the user.",
		},
	}

	fmt.Println(getChatCompletion(messages))
}
