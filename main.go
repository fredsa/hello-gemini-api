package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	// "cloud.google.com/go/vertexai/genai"    // Vertex AI.
	"github.com/google/generative-ai-go/genai" // Google AI.
	"google.golang.org/api/option"             // Google AI.
)

func main() {
	ctx := context.Background()

	// New Google AI client, using API key authorization.
	option := option.WithAPIKey(os.Getenv("API_KEY")) // Google AI.
	client, err := genai.NewClient(ctx, option)       // Google AI.

	// New Vertex AI client, using application default credentials.
	// client, err := genai.NewClient(ctx, "PROJECT_ID", "LOCATION") // Vertex AI.

	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}
	defer client.Close()

	// Start conversation with Gemini.
	converse(ctx, client.GenerativeModel("gemini-pro"))
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

		// Call the Gemini API.
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
