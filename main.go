package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

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

	// Use Gemini Pro model.
	model := client.GenerativeModel("gemini-pro")

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
		StopSequences: []string{
			`a`,
			`a\b`,
			`'a'`,
			`"a"`,
		},
	}

	// Configure safety settings.
	model.SafetySettings = []*genai.SafetySetting{
		{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockMediumAndAbove},
	}

	// Call the Gemini AI API.
	part := genai.Text(`Write a creative\exciting children's bedtime "story"`)
	resp, err := model.GenerateContent(ctx, part)
	if err != nil {
		log.Fatalf("Error sending message: %v\n", err)
	}

	// Display the response.
	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Printf("%v\n", part)
	}
}
