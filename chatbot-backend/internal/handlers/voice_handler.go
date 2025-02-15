package handlers

import (
	"fmt"
	"net/http"
	"os"

	"chatbot-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func HandleVoiceInput(c *gin.Context) {
	// Retrieve the file from the form data
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Save the file to a temporary location
	filePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error saving file: %v", err)})
		return
	}

	// Process the voice input
	response, err := services.HandleVoiceInput(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error processing voice input: %v", err)})
		return
	}

	// Send back response
	c.JSON(http.StatusOK, response)

	// Clean up the temporary file
	os.Remove(filePath)
}
