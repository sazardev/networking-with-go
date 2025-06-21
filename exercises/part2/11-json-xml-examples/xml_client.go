// xml_client.go
// HTTP client that fetches products and posts a new product to the XML server.
package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Fetch products (GET)
	resp, err := http.Get("http://localhost:8081/products")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var plist ProductList
	if err := xml.NewDecoder(resp.Body).Decode(&plist); err != nil {
		panic(err)
	}
	fmt.Println("Products from server:")
	for _, p := range plist.Products {
		fmt.Printf("- ID: %d, Name: %s, Tags: %v\n", p.ID, p.Name, p.Tags)
	}

	// Add a new product (POST)
	newProduct := Product{Name: "Thingamajig", Tags: []string{"novelty", "fun"}}
	data, _ := xml.Marshal(newProduct)
	resp2, err := http.Post("http://localhost:8081/products", "application/xml", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()
	var created Product
	body, _ := io.ReadAll(resp2.Body)
	xml.Unmarshal(body, &created)
	fmt.Println("Created product:", created)
}
