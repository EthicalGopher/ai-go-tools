package agenticrag

import (
	"fmt"

	"strings"

	"github.com/EthicalGopher/go-ai-tools/gemini"
	maincom "github.com/EthicalGopher/go-ai-tools/pinecone"
	list "github.com/EthicalGopher/go-ai-tools/struct"
	
)



func Getstored_data(api string,structuredinput []string, storedata maincom.Pincone_info) ([]string) {
	var pincone_data []string
	var embeddingstruct = list.Airesponse{
		Apikey: api,
		Model:  "gemini-embedding-exp-03-07",
		Input:  "",
		About:  "",
	}
	for _,input := range structuredinput{
		vector:=gemini.Maketextembedding(embeddingstruct,input)
		storedata.QueryVector = vector
		data := maincom.SearchData(storedata)
		if len(data) > 0 {
			pincone_data = append(pincone_data, data...)
		}
	}
	var finaluniquedata []string
	for _, data := range pincone_data {
		exists := false
		for _, uniqueData := range finaluniquedata {
			if uniqueData == data {
				exists = true
				break
			}
		}
		if !exists {
			finaluniquedata = append(finaluniquedata, data)
		}
	}

	return finaluniquedata
}

func GenerateFinalResponse(geminiapi,details,input string,storedata maincom.Pincone_info) string {
	firstinput:=input
	fmt.Println("First input: ", firstinput)
	var final_answer string
	for {

	


	// 1st agent



	about := `
	You are an expert query decomposition assistant. 
				Break down complex queries into precise, actionable sub-queries 
				that can be independently researched. Focus on extracting 
				core information needs and potential knowledge domains.
	
	`
	model := ""
	input = fmt.Sprintf(`
	Decompose this query into 2-3 specific sub-queries 
				that will help gather comprehensive information:
				
				Query: %s
				
				Provide sub-queries that are:
				- Specific and focused
				- Can be independently researched
				- Cover different aspects of the original query
	`, input)
	load1 := list.Airesponse{geminiapi, input, about, model}
	res, err := gemini.Generateresponse(load1)
	if err != nil {
		fmt.Println("Error generating response", err)
	}
	fmt.Println(res)
	res1 := strings.Split(res, "\n")
	var context string

	pincone_data := Getstored_data(geminiapi,res1, storedata)
	for i := range pincone_data {
		pincone_data[i] = strings.TrimSpace(pincone_data[i])
		context += fmt.Sprintf("%d %s\n",i, pincone_data[i])
	}

	fmt.Println("Context: ", context)

	// 2nd agent

	about = `
	You are an advanced AI research assistant. Generate 
				comprehensive, well-reasoned responses that:
				- Directly address the query
				- Synthesize information from multiple sources
				- Provide clear, step-by-step explanations
				- Highlight key insights and potential implications
	
	
	`
	model = ""
	input = fmt.Sprintf(`
	Query: %s

				Retrieved Context:
				%s

				Please generate a comprehensive response that:
				1. Answers the original query
				2. Integrates insights from retrieved documents
				3. Explains reasoning and connections
				4. Provides a structured, informative answer
				5. the response should be point-wise and clear no need to add numberings to the point just go to the next line`, 
				firstinput,context)

	load2 := list.Airesponse{geminiapi, input, about, model}
	res2, err := gemini.Generateresponse(load2)
	if err != nil {
		fmt.Println("Error generating response", err)
	}
	fmt.Println(res2)


	// 3rd verifing agent
	about = `
    You are a response verification assistant. Your task is to evaluate 
    whether the generated response is relevant to the original query. 
    Analyze the response and determine if it directly addresses the 
    original input query. Answer with "Yes" if it is relevant, or "No" 
    if it is not.

    Provide a brief explanation for your decision.
    `
    model = ""
    input = fmt.Sprintf(`
    Original Query: %s

    Generated Response:
    %s

    Is the generated response relevant to the original query? 
    Answer "Yes" or "No" and explain briefly.`, 
    firstinput, res2)

    load3 := list.Airesponse{geminiapi, input, about, model}
    res3, err := gemini.Generateresponse(load3)
    if err != nil {
        fmt.Println("Error verifying response", err)
    }
    fmt.Println("Verification Result: ", res3)
res3 = strings.ToLower(res3)
	if strings.Contains(res3, "yes") {
		final_answer = res2
		break
	}

}
if final_answer == "" {
    fmt.Println("Failed to generate a verified response after maximum attempts.")
    return "Unable to generate a response."
}
return final_answer


}