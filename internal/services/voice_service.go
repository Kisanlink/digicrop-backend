package services

import (
	"chatbot-backend/internal/models"
	"chatbot-backend/pkg/voice"
)

func HandleVoiceInput(audioFilePath string) (*models.VoiceResponse, error) {
	// Convert audio to text
	text, err := voice.ConvertSpeechToText(audioFilePath)
	if err != nil {
		return nil, err
	}

	// Call the OpenAI function to get the response
	response, err := voice.GetOpenAIResponse(text)
	if err != nil {
		return nil, err
	}

	return &models.VoiceResponse{Response: response}, nil
}
