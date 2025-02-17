package services

import (
	"chatbot-backend/pkg/openai"
	"chatbot-backend/internal/models"
)

func HandleUserQuery(query string) (*models.ChatResponse, error) {
	// Call the OpenAI function to get the response
	response, err := openai.GetOpenAIResponse(query)
	if err != nil {
		return nil, err
	}

	return &models.ChatResponse{Response: response}, nil
}
