package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const modelName = "gemini-pro"

func getApiKey() string {
	envKey := "GEMINI_API_KEY"
	apiKey, found := os.LookupEnv(envKey)
	if !found {
		log.Fatalf("Environment variable '%v' not set\n", envKey)
	}
	return apiKey
}

func main() {
	ctx := context.Background()

	// New client, using API key authorization.
	option := option.WithAPIKey(getApiKey())
	client, err := genai.NewClient(ctx, option)
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}
	defer client.Close()

	// Start conversation with Gemini.
	converse(ctx, client.GenerativeModel(modelName))
}

func converse(ctx context.Context, model *genai.GenerativeModel) {
	// Configure generation settings.
	temperature := float32(0.9)
	topK := int32(1)
	topP := float32(1.0)
	maxOutputTokens := int32(2048)

	model.GenerationConfig = genai.GenerationConfig{
		Temperature:     &temperature,
		TopK:            &topK,
		TopP:            &topP,
		MaxOutputTokens: &maxOutputTokens,
		StopSequences:   []string{`🎉`},
	}

	// Configure safety settings.
	model.SafetySettings = []*genai.SafetySetting{
		{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockMediumAndAbove},
	}

	// Multi-part request.
	parts := []genai.Part{
		genai.Text("Describe the character"),
		genai.Text("char: 🥞"),
		genai.Text("description: pancakes emoji"),
		genai.Text("char: 木"),
		genai.Text("description: Mandarin character mù"),
		genai.Text("char: 💩"),
		genai.Text("description: "),
	}
	// Call the Gemini AI API.
	resp, err := model.GenerateContent(ctx, parts...)
	if err != nil {
		log.Fatalf("Error sending message: %v\n", err)
	}

	// Display the response.
	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Printf("%v\n", part)
	}
}
