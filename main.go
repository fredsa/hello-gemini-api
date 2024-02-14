package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const jpegUrlBird = "https://t2.gstatic.com/licensed-image?q=tbn:ANd9GcRYL5x5KonnTYgeE-C-s09bBKuupEp0F0VsPYKNwpW6Xp-TFfoZ_iTporuLNOBWwm6HhOoek5cF"
const jpegUrlCat = "https://t0.gstatic.com/licensed-image?q=tbn:ANd9GcRbNuiexLb-Bsa2FR_HX5wRzIfI79zoNJk1F0kjbPvQ0O_MK3T-xobhQjN4fbYwDBa-RNdCx36t"

func getAPIKey() string {
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
	option := option.WithAPIKey(getAPIKey())
	client, err := genai.NewClient(ctx, option)
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}
	defer client.Close()

	// Use Gemini Pro Vision model.
	model := client.GenerativeModel("gemini-pro-vision")

	// Configure generation settings.
	temperature := float32(0.4)
	topK := int32(32)
	topP := float32(1.0)
	maxOutputTokens := int32(4096)

	model.GenerationConfig = genai.GenerationConfig{
		Temperature:     &temperature,
		TopK:            &topK,
		TopP:            &topP,
		MaxOutputTokens: &maxOutputTokens,
	}

	// Configure safety settings.
	model.SafetySettings = []*genai.SafetySetting{
		// {Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockMediumAndAbove},
		// {Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockMediumAndAbove},
		// {Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockMediumAndAbove},
		// {Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockMediumAndAbove},
	}

	// Multi-part request.
	parts := []genai.Part{
		genai.ImageData("jpeg", getBytes(jpegUrlCat)),
		genai.Text("Tell me a story about this animal"),
	}

	// Call the Gemini AI API.
	it := model.GenerateContentStream(ctx, parts...)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error iterating through response: %v\n", err)
		}

		// Display the the next part of streamed response.
		for _, part := range resp.Candidates[0].Content.Parts {
			fmt.Print(part)
			fmt.Println()
		}
	}
}

func getBytes(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching image %v: %v\n", url, err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading image bytes %v: %v\n", url, err)
	}
	return bytes
}
