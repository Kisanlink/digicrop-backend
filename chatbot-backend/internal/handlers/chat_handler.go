package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"chatbot-backend/internal/models"
	"chatbot-backend/internal/services"
)

func ChatHandler(c *gin.Context) {
	var request models.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Process the query
	response, err := services.HandleUserQuery(request.Query)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Error processing query: %v", err)})
		return
	}

	// Send back response
	c.JSON(200, response)
}
