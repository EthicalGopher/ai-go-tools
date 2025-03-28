package maincom

import (
	"context"
	"fmt"
	"log"

	"github.com/pinecone-io/go-pinecone/v3/pinecone"
)
type Pincone_info struct{
	IndexName string
	Namespace string
	API string
	Field string
	QueryVector []float32
	TopK int
}





func SearchData(load Pincone_info) []string {
	indexName := load.IndexName
	namespace := load.Namespace
	API := load.API
	field := load.Field



	queryVector := load.QueryVector
	ctx := context.Background()

	clientParams := pinecone.NewClientParams{
		ApiKey: API,
	}
	pc, err := pinecone.NewClient(clientParams)
	if err!=nil{
		log.Println(err)
	}
	idxModel, err := pc.DescribeIndex(ctx, indexName)
	if err != nil {
		log.Fatalf("Failed to describe index \"%v\": %v", indexName, err)
	}

	idxConnection, err := pc.Index(pinecone.NewIndexConnParams{Host: idxModel.Host, Namespace: namespace})
	if err != nil {
		log.Fatalf("Failed to create IndexConnection1 for Host %v: %v", idxModel.Host, err)
	}

	var K int
	if load.TopK==0{
		K=3
	}else{
		K=load.TopK
	}

	



	res, err := idxConnection.QueryByVectorValues(ctx, &pinecone.QueryByVectorValuesRequest{
		Vector:        queryVector,
		TopK:          uint32(K),
		IncludeValues: true,
		IncludeMetadata: true,
	})
	var values []string
	if err != nil {
		log.Fatalf("Error encountered when querying by vector: %v", err)
	} else {
		for _, match := range res.Matches {
			
			value:=fmt.Sprintln(match.Vector.Metadata.Fields[field])
			values=append(values, value)
		}
	}
	return values
}