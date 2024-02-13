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

	// Use Gemini Pro Vision model.
	model := client.GenerativeModel("gemini-pro-vision")

	// Multi-part prompt.
	prompt := []genai.Part{
		// TODO: Provide input image, update description to match the image.
		genai.Text("Screenshot: "),
		genai.ImageData("png", getBytes("astrohorse.png")),
		genai.Text("\nDescription: an astronaut riding a horse\n\nScreenshot: "),

		// TODO: Provide second image for which a description is desired.
		genai.ImageData("png", getBytes("camel.png")),
		genai.Text("\nDescription:"),
	}

	// Invoke API to generate response.
	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		log.Fatal(err)
	}

	// Display the response as formatted JSON output.
	bs, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(bs))
}

// Get file contents as bytes.
func getBytes(path string) []byte {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file %v: %v\n", path, err)
	}
	return bytes
}
