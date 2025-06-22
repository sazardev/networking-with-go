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

// Product represents a product with ID, Name, and a list of Tags (nested elements).
type Product struct {
	XMLName xml.Name `xml:"product"`
	ID      int      `xml:"id"`
	Name    string   `xml:"name"`
	Tags    []string `xml:"tags>tag"`
}

// ProductList wraps a list of products for XML decoding.
type ProductList struct {
	XMLName  xml.Name  `xml:"products"`
	Products []Product `xml:"product"`
}

func main() {
	// Fetch products (GET)
	resp, err := http.Get("http://localhost:8081/products") // Send GET request
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var plist ProductList
	if err := xml.NewDecoder(resp.Body).Decode(&plist); err != nil { // Decode XML into ProductList
		panic(err)
	}
	fmt.Println("Products from server:")
	for _, p := range plist.Products {
		fmt.Printf("- ID: %d, Name: %s, Tags: %v\n", p.ID, p.Name, p.Tags)
	}

	// Add a new product (POST)
	newProduct := Product{Name: "Thingamajig", Tags: []string{"novelty", "fun"}}
	data, _ := xml.Marshal(newProduct)                                                                  // Encode new product as XML
	resp2, err := http.Post("http://localhost:8081/products", "application/xml", bytes.NewBuffer(data)) // Send POST
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()
	var created Product
	body, _ := io.ReadAll(resp2.Body)
	xml.Unmarshal(body, &created) // Decode response
	fmt.Println("Created product:", created)
}
