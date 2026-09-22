package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// SolveMCQ sends the image to Gemini and asks for the multiple choice answer
func SolveMCQ(ctx context.Context, apiKey string, imageBytes []byte) (string, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	modelName := os.Getenv("GEMINI_MODEL")
	if modelName == "" {
		modelName = "gemini-3.8-flash"
	}
	model := client.GenerativeModel(modelName)
	prompt := genai.Text("Baca soal pilihan ganda dari gambar ini beserta pilihan jawabannya. Pilih satu jawaban yang paling tepat dan kembalikan HANYA huruf pilihannya saja (contoh: A, B, C, D, atau E).")
	imgData := genai.ImageData("image/png", imageBytes)

	resp, err := model.GenerateContent(ctx, prompt, imgData)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from model")
	}

	part := resp.Candidates[0].Content.Parts[0]
	if text, ok := part.(genai.Text); ok {
		return strings.TrimSpace(string(text)), nil
	}

	return "", fmt.Errorf("unexpected response type")
}
