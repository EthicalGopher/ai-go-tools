package gemini

import (
	"context"

	"log"

	list "github.com/EthicalGopher/go-ai-tools/struct"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)



func Maketextembedding (load list.Airesponse,text string) []float32{
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(load.Apikey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	if load.Model == ""{
		load.Model = "gemini-embedding-exp-03-07"
	}
	
	em := client.EmbeddingModel(load.Model)
	res, err := em.EmbedContent(ctx, genai.Text(text))
	
	if err != nil {
		panic(err)
	}
	
	return res.Embedding.Values
}
