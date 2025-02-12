package gemini

import (
	"context"
	
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"


)


func Maketextembedding (API,text,model string) []float32{
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(API))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	
	em := client.EmbeddingModel(model)
	res, err := em.EmbedContent(ctx, genai.Text(text))
	
	if err != nil {
		panic(err)
	}
	
	return res.Embedding.Values
}
