package vision

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type VisionClient interface {
	SendImageRequest(ctx context.Context, req ImageRequest) (*ImageResponse, error)
}

type visionClient struct {
	apiKey string
}

func NewVisionClient(apiKey string) VisionClient {
	return &visionClient{apiKey: apiKey}
}

type ImageRequest struct {
	Image     string `json:"image"`
	Prompt    string `json:"prompt"`
	Model     string `json:"model"`
	MaxTokens int    `json:"max_tokens"`
}

type ImageResponse struct {
	Analysis string `json:"analysis"`
}

func (vc *visionClient) SendImageRequest(ctx context.Context, req ImageRequest) (*ImageResponse, error) {
	// Prepare the payload for OpenAI Vision API
	messages := []map[string]interface{}{
		{
			"role": "user",
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": req.Prompt,
				},
				map[string]interface{}{
					"type": "image_url",
					"image_url": map[string]string{
						"url": fmt.Sprintf("data:image/jpeg;base64,%s", req.Image),
					},
				},
			},
		},
	}

	payload := map[string]interface{}{
		"model":      req.Model,
		"messages":   messages,
		"max_tokens": req.MaxTokens,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Send HTTP request to OpenAI API
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+vc.apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	// Parse the response
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no response from model")
	}

	return &ImageResponse{
		Analysis: result.Choices[0].Message.Content,
	}, nil
}
