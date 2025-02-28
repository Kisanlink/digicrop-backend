package services

import (
	"chatbot-backend/internal/models"
	"chatbot-backend/pkg/vision"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
)

type VisionService interface {
	AnalyzeImage(ctx context.Context, req models.ImageAnalysisRequest) (string, error)
}

type visionService struct {
	visionClient vision.VisionClient
}

func NewVisionService(vc vision.VisionClient) VisionService {
	return &visionService{visionClient: vc}
}

func (s *visionService) AnalyzeImage(ctx context.Context, req models.ImageAnalysisRequest) (string, error) {
	// Validate base64 image
	if _, err := base64.StdEncoding.DecodeString(req.ImageBase64); err != nil {
		return "", errors.New("invalid base64 image format")
	}

	// Prepare vision request
	visionReq := vision.ImageRequest{
		Image:     req.ImageBase64,
		Prompt:    req.Prompt,
		Model:     req.Model,
		MaxTokens: 1000,
	}

	// Send request to OpenAI Vision API
	response, err := s.visionClient.SendImageRequest(ctx, visionReq)
	if err != nil {
		return "", fmt.Errorf("vision API error: %v", err)
	}

	return response.Analysis, nil
}
