package voice

import (
	"bytes"
	"chatbot-backend/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sashabaranov/go-openai"
)

// Declare client variable
var client *openai.Client

// Initialize OpenAI Client for voice processing
func InitOpenAI() {
	apiKey := config.GetOpenAIKey()
	client = openai.NewClient(apiKey)
}

// Convert speech to text using OpenAI
func ConvertSpeechToText(audioFilePath string) (string, error) {
	audioFile, err := os.Open(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %v", err)
	}
	defer audioFile.Close()

	// Prepare the request
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", filepath.Base(audioFilePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}
	if _, err = io.Copy(fw, audioFile); err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}
	w.WriteField("model", "whisper-1")
	w.WriteField("language", "en")
	w.Close()

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &b)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+config.GetOpenAIKey())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error, status code: %d, status: %s", resp.StatusCode, resp.Status)
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return result.Text, nil
}

// Get response from OpenAI for the converted text
func GetOpenAIResponse(query string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini-2024-07-18",
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
