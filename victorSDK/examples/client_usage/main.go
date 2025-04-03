package main

import (
	"fmt"
	"log"
	"victorgo/victorSDK"
)

func main() {

	client := victorSDK.NewClient(nil)
	// defer client.Close()

	indexParams := victorSDK.CreateIndexCommandInput{
		IndexType: 0,
		Method:    0,
		Dims:      5,
		IndexName: "indice1",
	}

	v1 := victorSDK.InsertVectorCommandInput{
		IndexName: "indice1",
		ID:        1,
		Vector:    []float32{0.1, 0.2, 0.3, 0.4, 0.5},
	}

	v2 := victorSDK.InsertVectorCommandInput{
		IndexName: "indice1",
		ID:        2,
		Vector:    []float32{0.6, 0.7, 0.8, 0.9, 1.0},
	}
	v3 := victorSDK.InsertVectorCommandInput{
		IndexName: "indice1",
		ID:        3,
		Vector:    []float32{1.1, 1.2, 1.3, 1.4, 1.5},
	}

	searchVector := victorSDK.SearchVectorCommandInput{
		IndexName: "indice1",
		TopK:      3,
		Vector:    []float32{1.5, 0.2, 2, -0.4, 0.9},
	}

	indexResult, err := client.CreateIndex(&indexParams)
	if err != nil {
		log.Fatalf("Index creation error: %+v", err)
	}

	fmt.Printf("Index created successfully: %+v\n", indexResult.Results)

	insertResult1, err := client.InsertVector(&v1)
	if err != nil {
		log.Fatalf("Vector insertion error: %+v", err)
	}
	fmt.Printf("Vector insertion successfully: %+v\n", insertResult1.Results)

	insertResult2, err := client.InsertVector(&v2)
	if err != nil {
		log.Fatalf("Insertion vector error: %+v", err)
	}
	fmt.Printf("Vector insertion successfully: %+v\n", insertResult2.Results)

	insertResult3, err := client.InsertVector(&v3)
	if err != nil {
		log.Fatalf("Insertion vector Error: %+v", err)
	}
	fmt.Printf("Vector insertion successfully: %+v\n", insertResult3.Results)

	searchResult, err := client.SearchVector(&searchVector)
	if err != nil {
		log.Fatalf("Search error: %+v", err)
	}
	fmt.Printf("Search results: %+v\n%+v\n", searchResult.Status, searchResult.Results)

	client.Close()
	fmt.Println("Client closed")
}
