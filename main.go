package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const jpegUrlBird = "https://t2.gstatic.com/licensed-image?q=tbn:ANd9GcRYL5x5KonnTYgeE-C-s09bBKuupEp0F0VsPYKNwpW6Xp-TFfoZ_iTporuLNOBWwm6HhOoek5cF"
const jpegUrlCat = "https://t0.gstatic.com/licensed-image?q=tbn:ANd9GcRbNuiexLb-Bsa2FR_HX5wRzIfI79zoNJk1F0kjbPvQ0O_MK3T-xobhQjN4fbYwDBa-RNdCx36t"

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

	// Use Gemini Pro Vision model.
	model := client.GenerativeModel("gemini-pro-vision")

	// Configure generation settings.
	model.SetTemperature(0.4)
	model.SetTopK(32)
	model.SetTopP(1.0)
	model.SetMaxOutputTokens(4096)

	// Configure safety settings.
	model.SafetySettings = []*genai.SafetySetting{
		{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockMediumAndAbove},
	}

	// Image only request.
	part := genai.ImageData("jpeg", getJpegImageBytes(jpegUrlBird))

	// Call the Gemini AI API.
	resp, err := model.GenerateContent(ctx, part)
	if err != nil {
		log.Fatalf("Error sending message: %v\n", err)
	}

	// Display the response.
	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Printf("%v\n", part)
	}
}

func getJpegImageBytes(url string) []byte {
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
