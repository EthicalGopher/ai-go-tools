package gemini

import (
	"context"
	"strings"

	list "github.com/EthicalGopher/go-ai-tools/struct"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func Generateresponse(load list.Airesponse) (string, error) {
	

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(load.Apikey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel(load.Model)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(load.About)},
	}
	resp, err := model.GenerateContent(ctx, genai.Text(load.Input))
	if err != nil {
		return "", err
	}

	
	var responseText strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			responseText.WriteString(string(text))
		}
	}

	return responseText.String(), nil
}
