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

	// // For text-only input, use the Gemini Pro model.
	// model := client.GenerativeModel("gemini-pro")

	// text := "Parrots can be green and live a long time."
	// res, err := model.CountTokens(ctx, genai.Text(text))
	// if err != nil {
	// 	log.Fatalf("Unable to count tokens: %v\n", err)
	// }

	// fmt.Printf("Input text (%v bytes): %v\n", len(text), text)
	// fmt.Printf("Total tokens: %v\n", res.TotalTokens)

	// For text-and-image input (multimodal), use the Gemini Pro Vision model.
	model := client.GenerativeModel("gemini-pro-vision")

	text := "Parrots can be green and live a long time."
	imageBytes := getJpegImageBytes(string(jpegUrlBird))

	parts := []genai.Part{
		genai.Text(text),
		genai.ImageData("jpeg", imageBytes),
	}

	res, err := model.CountTokens(ctx, parts...)
	if err != nil {
		log.Fatalf("Unable to count tokens: %v\n", err)
	}

	fmt.Printf("Input text (%v bytes): %v\n", len(text), text)
	fmt.Printf("Input image (%v bytes): %v\n", len(imageBytes), jpegUrlBird)
	fmt.Printf("Total: %v bytes, %v tokens\n", len(text)+len(imageBytes), res.TotalTokens)
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
