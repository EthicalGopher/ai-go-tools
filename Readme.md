# Go RAG Agent System

This project implements a Retrieval-Augmented Generation (RAG) agent system in Go, leveraging the Gemini API for natural language processing and the Pinecone vector database for efficient information retrieval. The system is designed to answer complex queries by breaking them down into sub-queries, retrieving relevant information, and synthesizing a comprehensive response.

## Overview

The core idea behind this project is to create an intelligent agent that can:

1.  **Understand Complex Queries:** Decompose complex user queries into smaller, more manageable sub-queries.
2.  **Retrieve Relevant Information:** Search a vector database (Pinecone) for information relevant to the sub-queries.
3.  **Generate Comprehensive Responses:** Synthesize the retrieved information into a well-structured and informative answer to the original query.
4. **Verify the response:** Check if the generated response is relevant to the original query.

## Project Structure

The project is organized into several packages, each responsible for a specific aspect of the system:

*   **`agenticrag`:** Contains the main logic for the agent, including query decomposition, information retrieval, response generation, and verification.
*   **`gemini`:** Handles interactions with the Google Gemini API, including text embedding and response generation.
*   **`pinecone`:** Manages interactions with the Pinecone vector database, including data storage and retrieval.
*   **`struct`:** Defines custom data structures used throughout the project.

## Key Components

### 1. `agenticrag` Package (`agent.go`)

This package contains the core logic for the RAG agent.

*   **`Getstored_data(api string, structuredinput []string, storedata maincom.Pincone_info) []string`**
    *   **Purpose:** Retrieves relevant data from the Pinecone vector database based on a list of input strings (sub-queries).
    *   **Functionality:**
        1.  Takes a list of strings (`structuredinput`) representing sub-queries.
        2.  Uses the Gemini API to generate embeddings for each sub-query.
        3.  Queries the Pinecone database using the generated embeddings.
        4.  Collects the retrieved data from Pinecone.
        5.  Removes duplicate entries from the retrieved data.
        6.  Returns a list of unique data strings.
    *   **Parameters:**
        *   `api`: Gemini API key.
        *   `structuredinput`: A slice of strings representing sub-queries.
        *   `storedata`: A `Pincone_info` struct containing Pinecone connection details (index name, namespace, API key, field, etc.).
    * **Return:**
        * `[]string`: A slice of unique string retrieved from pinecone.

*   **`GenerateFinalResponse(geminiapi, details, input string, storedata maincom.Pincone_info) string`**
    *   **Purpose:** Orchestrates the entire RAG process to generate a final response to a user's query.
    *   **Functionality:**
        1.  **Query Decomposition (Agent 1):**
            *   Uses the Gemini API to decompose the original query into 2-3 sub-queries.
            *   Employs a specific prompt to guide the Gemini model in this task.
        2.  **Information Retrieval:**
            *   Calls `Getstored_data` to retrieve relevant information from Pinecone based on the sub-queries.
        3.  **Response Generation (Agent 2):**
            *   Uses the Gemini API again to generate a comprehensive response.
            *   Provides the original query and the retrieved context to the Gemini model.
            *   Employs a prompt that instructs the model to synthesize information, explain reasoning, and provide a structured answer.
        4. **Verification (Agent 3):**
            * Uses the Gemini API to verify if the generated response is relevant to the original query.
            * If the response is not relevant it will try again.
        5.  **Iteration:**
            *   The process iterates until a verified response is generated.
        6. **Return:**
            * `string`: The final answer or a message if it fails to generate a response.
    *   **Parameters:**
        *   `geminiapi`: Gemini API key.
        *   `details`: Not used in the current implementation.
        *   `input`: The original user query.
        *   `storedata`: A `Pincone_info` struct containing Pinecone connection details.

### 2. `gemini` Package (`response.go`, `textembeding.go`)

This package handles interactions with the Google Gemini API.

*   **`Generateresponse(load list.Airesponse) (string, error)`**
    *   **Purpose:** Sends a prompt to the Gemini API and returns the generated response.
    *   **Functionality:**
        1.  Initializes a Gemini client using the provided API key.
        2.  Sets the system instruction (context/role) for the Gemini model.
        3.  Sends the user's input to the model.
        4.  Returns the generated text response.
    *   **Parameters:**
        *   `load`: An `Airesponse` struct containing the API key, input text, system instruction, and model name.
    *   **Return:**
        *   `string`: The generated response.
        *   `error`: Any error encountered during the process.

*   **`Maketextembedding(load list.Airesponse, text string) []float32`**
    *   **Purpose:** Generates a text embedding (vector representation) for a given text using the Gemini API.
    *   **Functionality:**
        1.  Initializes a Gemini client.
        2.  Uses the specified embedding model (default: `gemini-embedding-exp-03-07`).
        3.  Sends the text to the model for embedding.
        4.  Returns the resulting vector.
    *   **Parameters:**
        *   `load`: An `Airesponse` struct containing the API key and model name.
        *   `text`: The text to embed.
    *   **Return:**
        *   `[]float32`: The generated vector embedding.

### 3. `pinecone` Package (`searchdata.go`, `adddata.go`)

This package manages interactions with the Pinecone vector database.

*   **`SearchData(load Pincone_info) []string`**
    *   **Purpose:** Searches the Pinecone database for data similar to a given query vector.
    *   **Functionality:**
        1.  Connects to the Pinecone index using the provided credentials.
        2.  Queries the index using the provided query vector.
        3.  Retrieves the top-K (default: 3) most similar data entries.
        4.  Extracts the text data from the metadata of the retrieved entries.
        5.  Returns a list of the retrieved text data.
    *   **Parameters:**
        *   `load`: A `Pincone_info` struct containing Pinecone connection details, the query vector, and the desired number of results (TopK).
    *   **Return:**
        *   `[]string`: A list of retrieved text data.

*   **`Addonestring(indexName, namespace, API, text string, vector []float32)`**
    *   **Purpose:** Adds a single text entry and its corresponding vector to the Pinecone database.
    *   **Functionality:**
        1.  Connects to the Pinecone index.
        2.  Creates a new vector with a unique ID, the provided vector values, and metadata containing the text.
        3.  Upserts (inserts or updates) the vector into the index.
    *   **Parameters:**
        *   `indexName`: The name of the Pinecone index.
        *   `namespace`: The namespace within the index.
        *   `API`: The Pinecone API key.
        *   `text`: The text to store.
        *   `vector`: The vector representation of the text.

*   **`AddTextsToMetadata(indexName, namespace, API string, texts []string, vector []float32) error`**
    *   **Purpose:** Adds multiple texts as metadata to a single vector in the Pinecone database.
    *   **Functionality:**
        1.  Connects to the Pinecone index.
        2.  Creates a metadata map where each text is stored under a unique key (e.g., "text1", "text2").
        3.  Creates a new vector with a unique ID, the provided vector values, and the metadata map.
        4.  Upserts the vector into the index.
    *   **Parameters:**
        *   `indexName`: The name of the Pinecone index.
        *   `namespace`: The namespace within the index.
        *   `API`: The Pinecone API key.
        *   `texts`: A slice of texts to store as metadata.
        *   `vector`: The vector representation.
    *   **Return:**
        *   `error`: Any error encountered during the process.

### 4. `struct` Package

*   **`Airesponse`**
    *   **Purpose:** A struct to hold the data for the gemini api.
    *   **Fields:**
        * `Apikey`: The api key for the gemini api.
        * `Input`: The input for the gemini api.
        * `About`: The context for the gemini api.
        * `Model`: The model for the gemini api.

*   **`Pincone_info`**
    *   **Purpose:** A struct to hold the data for the pinecone api.
    *   **Fields:**
        * `IndexName`: The name of the pinecone index.
        * `Namespace`: The namespace of the pinecone index.
        * `API`: The api key for the pinecone api.
        * `Field`: The field to search in the pinecone index.
        * `QueryVector`: The query vector for the pinecone api.
        * `TopK`: The number of results to return from the pinecone api.

## How to Use

1.  **Prerequisites:**
    *   Go installed on your system.
    *   A Google Gemini API key.
    *   A Pinecone API key and a configured index.
    *   The required go packages.
        ```bash
        go get github.com/google/generative-ai-go/genai
        go get github.com/pinecone-io/go-pinecone/v3/pinecone
        go get github.com/google/uuid
        ```

2.  **Configuration:**
    *   Set your Gemini API key and Pinecone API key in your code.
    *   Configure the `Pincone_info` struct with your Pinecone index name, namespace, and other relevant details.

3.  **Running the Agent:**
    *   Call the `GenerateFinalResponse` function from your main program, providing the Gemini API key, the user's query, and the `Pincone_info` struct.

    ```go
    package main

    import (
        "fmt"
        "github.com/EthicalGopher/go-ai-tools/agenticrag"
        maincom "github.com/EthicalGopher/go-ai-tools/pinecone"
    )

    func main() {
        geminiAPIKey := "YOUR_GEMINI_API_KEY"
        pineconeAPIKey := "YOUR_PINECONE_API_KEY"
        indexName := "YOUR_PINECONE_INDEX_NAME"
        namespace := "YOUR_PINECONE_NAMESPACE"
        field := "text" // The field in your Pinecone metadata that contains the text

        storedata := maincom.Pincone_info{
            IndexName: indexName,
            Namespace: namespace,
            API:       pineconeAPIKey,
            Field:     field,
        }

        userQuery := "What are the benefits of using a RAG system?"
        response := agenticrag.GenerateFinalResponse(geminiAPIKey, "", userQuery, storedata)
        fmt.Println("Final Response:\n", response)
    }
    ```

## Future Improvements

*   **Enhanced Iteration:** Implement a more sophisticated iteration strategy to refine responses if the verification agent rejects them.
*   **Error Handling:** Add more robust error handling to gracefully handle issues like API failures or empty search results.
*   **Dynamic Sub-query Generation:** Allow the agent to dynamically adjust the number of sub-queries based on the complexity of the input.
*   **Memory:** Implement a memory mechanism to allow the agent to learn from past interactions.
*   **Concurrency:** Use goroutines to perform tasks concurrently and improve performance.
* **More testing:** Add more test to the code.

## Contributing

Contributions to this project are welcome! Please feel free to submit pull requests or open issues to discuss potential improvements.

## License

This project is licensed under the MIT License.
