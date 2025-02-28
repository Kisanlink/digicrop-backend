package handlers

import (
	"chatbot-backend/internal/models"
	"chatbot-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	visionService services.VisionService
}

func NewImageHandler(vs services.VisionService) *ImageHandler {
	return &ImageHandler{visionService: vs}
}

func (h *ImageHandler) HandleImageAnalysis(c *gin.Context) {
	var req models.ImageAnalysisRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ImageAnalysisResponse{Error: "Invalid request payload"})
		return
	}

	// Set default model if not provided
	if req.Model == "" {
		req.Model = "gpt-4o"
	}

	// Call vision service to analyze the image
	analysis, err := h.visionService.AnalyzeImage(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ImageAnalysisResponse{Error: err.Error()})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, models.ImageAnalysisResponse{Analysis: analysis})
}
