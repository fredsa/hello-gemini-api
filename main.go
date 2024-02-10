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

	// Use PaLM 2 for Text (text-bison) model.
	model := client.GenerativeModel("text-bison-001")

	// Configure generation settings.
	temperature := float32(0.7)
	candidateCount := int32(4)
	topK := int32(40)
	topP := float32(0.95)
	maxOutputTokens := int32(1024)

	model.GenerationConfig = genai.GenerationConfig{
		Temperature:     &temperature,
		CandidateCount:  &candidateCount,
		TopK:            &topK,
		TopP:            &topP,
		MaxOutputTokens: &maxOutputTokens,
		StopSequences: []string{
			"a",
			"a\b",
			"'a'",
			`"a"`,
		},
	}

	// Configure safety settings.
	model.SafetySettings = []*genai.SafetySetting{
		{Category: genai.HarmCategoryDerogatory, Threshold: genai.HarmBlockLowAndAbove},
		{Category: genai.HarmCategoryToxicity, Threshold: genai.HarmBlockLowAndAbove},
		{Category: genai.HarmCategoryViolence, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategorySexual, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategoryMedical, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategoryDangerous, Threshold: genai.HarmBlockMediumAndAbove},
	}

	mood := "exciting"
	prompt := fmt.Sprintf("Write a creative%s children&#39;s bedtime &quot;story&quot;", mood)

	// Call the Gemini AI API.
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatalf("Error sending message: %v\n", err)
	}

	// Display the response.
	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Printf("%v\n", part)
	}
}
