package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Get API key from environment variable.
func getApiKey() string {
	envKey := "API_KEY"
	apiKey, found := os.LookupEnv(envKey)
	if !found {
		log.Fatalf("Environment variable '%v' not set\n", envKey)
	}
	return apiKey
}

func main() {
	ctx := context.Background()

	// New client, using API key authorization.
	client, err := genai.NewClient(ctx, option.WithAPIKey(getApiKey()))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// For text-only input, use the Gemini Pro model.
	model := client.GenerativeModel("gemini-pro")

	// Establish chat history.
	session := model.StartChat()
	session.History = []*genai.Content{{
		Parts: []genai.Part{genai.Text("Hello, I have 2 dogs in my house.")},
		Role:  "user",
	}, {
		Parts: []genai.Part{genai.Text("Great to meet you. What would you like to know?")},
		Role:  "model",
	}}

	// Send next message in chat and generate response.
	resp, err := session.SendMessage(ctx, genai.Text("How many paws are in my house?"))
	if err != nil {
		log.Fatal(err)
	}

	// Display the response as formatted JSON output.
	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
}
