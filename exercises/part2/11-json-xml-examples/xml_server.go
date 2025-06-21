// xml_server.go
// HTTP server that serves a list of products as XML and receives new products via POST.
package main

import (
	"encoding/xml"
	"io"
	"net/http"
)

type Product struct {
	XMLName xml.Name `xml:"product"`
	ID      int      `xml:"id"`
	Name    string   `xml:"name"`
	Tags    []string `xml:"tags>tag"`
}

type ProductList struct {
	XMLName  xml.Name  `xml:"products"`
	Products []Product `xml:"product"`
}

var (
	products = []Product{
		{ID: 1, Name: "Widget", Tags: []string{"gadget", "tool"}},
		{ID: 2, Name: "Gizmo", Tags: []string{"device"}},
	}
)

// getProducts handles GET /products and returns the list as XML.
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	mu.Lock()
	defer mu.Unlock()
	xml.NewEncoder(w).Encode(ProductList{Products: products})
}

// addProduct handles POST /products and adds a new product from XML body.
func addProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if err := xml.Unmarshal(body, &p); err != nil {
		http.Error(w, "Invalid XML", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	p.ID = len(products) + 1
	products = append(products, p)
	w.WriteHeader(http.StatusCreated)
	xml.NewEncoder(w).Encode(p)
}

func main() {
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getProducts(w, r)
		} else if r.Method == http.MethodPost {
			addProduct(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":8081", nil)
}
