package maincom

import (
	"context"

	"log"
"fmt"
	"github.com/google/uuid"
	"github.com/pinecone-io/go-pinecone/v3/pinecone"
	"google.golang.org/protobuf/types/known/structpb"
)

func Addonestring(indexName,namespace,API, text string, vector []float32) {
	ctx := context.Background()

	clientParams := pinecone.NewClientParams{
		ApiKey: API,
	}
	pc, err := pinecone.NewClient(clientParams)
	if err != nil {
		log.Fatalf("Failed to create Client: %v", err)
	}

	vectorData := []struct {
		Values []float32
		Text   string
	}{
		{
			Values: vector,
			Text:   text,
		},
	}

	vectors := make([]*pinecone.Vector, len(vectorData))

	metadataMap := make(map[string]interface{})
	metadataMap["text"] = vectorData[0].Text

	metadata, err := structpb.NewStruct(metadataMap)
	if err != nil {
		panic(err)
	}

	vectors[0] = &pinecone.Vector{
		Id:       IDgen(),
		Values:   &vectorData[0].Values, // No & here
		Metadata: metadata,
	}

	idxModel, err := pc.DescribeIndex(ctx, indexName)
	if err != nil {
		log.Fatalf("Failed to describe index \"%v\": %v", indexName, err)
	}

	idxConnection, err := pc.Index(pinecone.NewIndexConnParams{Host: idxModel.Host, Namespace: namespace})
	if err != nil {
		log.Fatalf("Failed to create IndexConnection1 for Host %v: %v", idxModel.Host, err)
	}

	_, err = idxConnection.UpsertVectors(ctx, vectors)
	if err != nil {
		log.Println(err)
	}

}

func IDgen() string {
		u,err:=uuid.NewUUID()	
		if err!=nil{
			log.Fatalln(err)
		}

return u.String()
}



func AddTextsToMetadata(indexName,namespace, API string, texts []string, vector []float32) error {
	ctx := context.Background()

	clientParams := pinecone.NewClientParams{
			ApiKey: API,
	}
	pc, err := pinecone.NewClient(clientParams)
	if err != nil {
			return err
	}

	idxModel, err := pc.DescribeIndex(ctx, indexName)
	if err != nil {
			return err
	}

	idxConnection, err := pc.Index(pinecone.NewIndexConnParams{Host: idxModel.Host, Namespace: namespace})
	if err != nil {
			return err
	}

	metadataMap := make(map[string]interface{})
	for i, text := range texts {
			key := fmt.Sprintf("text%d", i+1) // Create keys like "text1", "text2", etc.
			metadataMap[key] = text
	}

	metadata, err := structpb.NewStruct(metadataMap)
	if err != nil {
			return err
	}
	vectorCopy := make([]float32, len(vector))
	copy(vectorCopy, vector)
	vectorPtr := &vectorCopy

	v := &pinecone.Vector{
			Id:       IDgen(),
			Values:   vectorPtr,
			Metadata: metadata,
	}

	_, err = idxConnection.UpsertVectors(ctx, []*pinecone.Vector{v})
	return err
}