package models

// ImageAnalysisRequest represents the incoming request for image analysis
type ImageAnalysisRequest struct {
	ImageBase64 string `json:"image_base64" binding:"required"` // Base64-encoded image
	Prompt      string `json:"prompt" binding:"required"`       // User prompt for analysis
	Model       string `json:"model"`                           // Model to use (e.g., "gpt-4o")
}

// ImageAnalysisResponse represents the response from image analysis
type ImageAnalysisResponse struct {
	Analysis string `json:"analysis"`        // Analysis result from OpenAI
	Error    string `json:"error,omitempty"` // Error message (if any)
}
