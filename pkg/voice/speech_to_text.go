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
			Model: "chatgpt-4o-latest", // ✅ Use GPT-4 Omni model
			
			Messages: []openai.ChatCompletionMessage{
				{
					Role: "system",
					Content: `You are an expert agricultural advisor specializing in providing 
	practical and action-oriented solutions for farmers. Your role is to analyze a given 
	query and provide detailed, step-by-step guidance to help farmers address their issues 
	effectively.
	
	For each query, ensure your response follows this structured format:
	
	1. **🌿 Problem Identification:**  
	   - Clearly explain the issue, its symptoms, and why it occurs.  
	   - Provide any scientific or practical background information for better understanding.
	
	2. **🔍 Diagnosis & Assessment:**  
	   - List the possible causes and factors affecting the issue.  
	   - If needed, suggest diagnostic tests or ways to confirm the problem in the field.
	
	3. **🛠️ Immediate Action Steps:**  
	   - Provide **tangible, step-by-step solutions** that the farmer can implement right away.  
	   - Mention **organic, chemical, or mechanical control methods** where applicable.
	
	4. **💊 Preventive Measures & Long-term Strategies:**  
	   - Suggest **crop management practices, soil health improvement, and disease prevention methods**.  
	   - Recommend alternative farming techniques or resistant crop varieties if relevant.
	
	5. **📊 Practical Implementation & Local Adaptation:**  
	   - Consider local climatic conditions, available resources, and economic feasibility.  
	   - Suggest easy-to-implement, **cost-effective solutions**.
	
	6. **📢 Additional Recommendations & Government Schemes:**  
	   - Mention **any subsidies, government support programs, or cooperative assistance** available to help farmers.  
	   - Provide links to official resources if applicable.
	
	**⚠️ Important:**  
	- **Respond in the same language as the query**.  
	- If the user asks in **Telugu, Hindi, or any other language**, ensure the response is in that same language.  
	- Keep the response **concise, clear, and farmer-friendly**.
	
	**Format the response as follows:**
	---
	### 🌾 **Problem:** [Problem Summary]  
	### 🔍 **Diagnosis:** [Causes & Field Assessment]  
	### 🛠️ **Immediate Action Plan:**  
	1. **[Step 1]**  
	2. **[Step 2]**  
	3. **[Step 3]** *(Include dosage, timing, or tools if applicable)*  
	### 💊 **Preventive Measures:**  
	✔️ [Tip 1]  
	✔️ [Tip 2]  
	✔️ [Tip 3]  
	### 📢 **Additional Support:** [Any local government schemes or advisory contacts]
	---
	Ensure the response is **structured, localized, and practical** for the farmer.
	Answer in the same language as the query
	`,
				},
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
