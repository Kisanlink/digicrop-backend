package main

import (
	"chatbot-backend/config"
	"chatbot-backend/internal/handlers"
	"chatbot-backend/pkg/openai"
	"chatbot-backend/pkg/voice"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Ensure uploads directory exists
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		if err := os.Mkdir("uploads", 0755); err != nil {
			fmt.Printf("Error creating uploads directory: %v\n", err)
			return
		}
	}
	// Initialize OpenAI and voice processing
	openai.InitOpenAI()
	voice.InitOpenAI() // Initialize the voice module

	// Set up Gin router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Define routes
	r.POST("/api/v1/chat", handlers.ChatHandler) // Handles text-based chat queries
	r.POST("/api/v1/voice", handlers.HandleVoiceInput)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
