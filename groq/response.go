package groq

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestPayload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type StreamChoice struct {
	Delta struct {
		Content *string `json:"content"`
	} `json:"delta"`
}

type StreamResponse struct {
	Choices []StreamChoice `json:"choices"`
}



func Ragfromgroq(apiKey, input string) string {
	url := "https://api.groq.com/openai/v1/chat/completions"

	payload := RequestPayload{
		Model: "llama-3.3-70b-versatile",
		Messages: []Message{
			{
				Role:    "user",
				Content: input,
			},
		},
		Stream: true,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling payload: %v\n", err)
		return "Error marshaling payload"
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return "Error creating request"
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return "Error making request"
	}
	defer resp.Body.Close()

	// Use bufio.Scanner to read SSE stream
	scanner := bufio.NewScanner(resp.Body)
	var fullResponse strings.Builder // Accumulate the full response
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines or non-data lines
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		// Remove "data: " prefix
		jsonStr := strings.TrimPrefix(line, "data: ")

		// Check for stream end marker
		if jsonStr == "[DONE]" {
			break
		}

		// Parse the JSON chunk
		var chunk StreamResponse
		err := json.Unmarshal([]byte(jsonStr), &chunk)
		if err != nil {
			fmt.Printf("Error decoding JSON: %v\nRaw data: %s\n", err, jsonStr)
			continue
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != nil {
			output := *chunk.Choices[0].Delta.Content
			fullResponse.WriteString(output) // Append to the full response
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading stream: %v\n", err)
		return "Error reading stream"
	}
	return fullResponse.String() // Return the accumulated response
}
