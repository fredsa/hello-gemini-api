package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const modelName = "gemini-pro" // OR "gemini-pro-vision".

func main() {
	ctx := context.Background()
	apiKey, found := os.LookupEnv("GEMINI_API_KEY")
	if !found {
		log.Fatal("Missing API key")
	}

	// New client, using API key authorization.
	option := option.WithAPIKey(apiKey)
	client, err := genai.NewClient(ctx, option)
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}
	defer client.Close()

	// Start conversation with Gemini.
	converse(ctx, client.GenerativeModel(modelName))
}

func converse(ctx context.Context, model *genai.GenerativeModel) {
	// Get user input from stdin.
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n>> ")

		// Read user prompt.
		prompt, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v\n", err)
		}

		// Call the Gemini AI API.
		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			log.Fatalf("Error generating content: %v\n", err)
		}

		// Display the response.
		for _, part := range resp.Candidates[0].Content.Parts {
			fmt.Printf("%v\n", part)
		}
	}
}
