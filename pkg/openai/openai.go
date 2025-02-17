package openai

import (
	"chatbot-backend/config"
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// Declare client variable
var client *openai.Client

// Initialize OpenAI Client
func InitOpenAI() {
	apiKey := config.GetOpenAIKey()
	client = openai.NewClient(apiKey)
}

// Get response from OpenAI
func GetOpenAIResponse(query string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4o",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: query,
				},
			},
			MaxTokens:   1000,
			Temperature: 0.7,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to get response from OpenAI: %v", err)
	}

	// Return first message from OpenAI response
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("empty response from OpenAI")
}
