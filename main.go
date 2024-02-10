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

	// Use model that can create text embeddings.
	em := client.EmbeddingModel("embedding-001")
	prompt := genai.Text("The quick brown fox jumps over the lazy dog.")

	// https://cloud.google.com/vertex-ai/docs/generative-ai/embeddings/get-multimodal-embeddings#supported-models
	// Use model that can create multimodal embeddings.
	// em := client.EmbeddingModel("multimodalembedding")
	// prompt := genai.ImageData("jpeg", getJpegImageBytes(jpegUrlBird))

	res, err := em.EmbedContent(ctx, prompt)
	if err != nil {
		log.Fatalf("Unable to create embedding: %v\n", err)
	}

	values := res.Embedding.Values
	fmt.Printf("Embedding dimension: %v\n", len(values))
	fmt.Printf("Embedding values: %v, %v, ...\n", values[0], values[1])
	// fmt.Println(res.Embedding.Values)
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
