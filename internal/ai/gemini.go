package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// FetchModels gets a list of available Gemini models that support generateContent
func FetchModels(ctx context.Context, apiKey string) ([]string, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	var models []string
	iter := client.ListModels(ctx)
	for {
		m, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		
		isGen := false
		for _, method := range m.SupportedGenerationMethods {
			if method == "generateContent" {
				isGen = true
				break
			}
		}
		
		if isGen {
			name := strings.TrimPrefix(m.Name, "models/")
			
			// Filter to only include Text-out models
			lowerName := strings.ToLower(name)
			
			// Hanya izinkan gemini pro, flash, atau flash-lite
			isGemini := strings.HasPrefix(lowerName, "gemini-")
			isProOrFlash := strings.Contains(lowerName, "-pro") || strings.Contains(lowerName, "-flash")
			
			// Buang varian non-teks dan model lain yang tidak diinginkan
			isTextOut := !strings.Contains(lowerName, "image") &&
				!strings.Contains(lowerName, "tts") &&
				!strings.Contains(lowerName, "audio") &&
				!strings.Contains(lowerName, "lyria") &&
				!strings.Contains(lowerName, "veo") &&
				!strings.Contains(lowerName, "transcribe") &&
				!strings.Contains(lowerName, "live") &&
				!strings.Contains(lowerName, "robotics") &&
				!strings.Contains(lowerName, "computer-use") &&
				!strings.Contains(lowerName, "omni")

			if isGemini && isProOrFlash && isTextOut {
				models = append(models, name)
			}
		}
	}
	
	return models, nil
}

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
