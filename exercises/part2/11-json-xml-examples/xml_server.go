// xml_server.go
// HTTP server that serves a list of products as XML and receives new products via POST.
package xmlserver

import (
	"encoding/xml"
	"io"
	"net/http"
	"sync"
)

// Product represents a product with ID, Name, and a list of Tags (nested elements).
type Product struct {
	XMLName xml.Name `xml:"product"` // Root element for each product
	ID      int      `xml:"id"`      // Product ID
	Name    string   `xml:"name"`    // Product name
	Tags    []string `xml:"tags>tag"`// Nested <tags><tag>...</tag></tags>
}

// ProductList wraps a list of products for XML encoding.
type ProductList struct {
	XMLName  xml.Name  `xml:"products"` // Root element for the list
	Products []Product `xml:"product"`  // Each product as a child
}

var (
	// products holds the in-memory list of products.
	products = []Product{
		{ID: 1, Name: "Widget", Tags: []string{"gadget", "tool"}},
		{ID: 2, Name: "Gizmo", Tags: []string{"device"}},
	}
	// mu is a mutex to protect concurrent access to products slice.
	mu sync.Mutex
)

// getProducts handles GET /products and returns the list as XML.
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml") // Set response type
	mu.Lock()   // Lock to safely read products
	defer mu.Unlock()
	xml.NewEncoder(w).Encode(ProductList{Products: products}) // Encode as XML
}

// addProduct handles POST /products and adds a new product from XML body.
func addProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	body, err := io.ReadAll(r.Body) // Read request body
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if err := xml.Unmarshal(body, &p); err != nil { // Parse XML into Product
		http.Error(w, "Invalid XML", http.StatusBadRequest)
		return
	}
	mu.Lock() // Lock to safely modify products
	defer mu.Unlock()
	p.ID = len(products) + 1 // Assign a new ID
	products = append(products, p) // Add to list
	w.WriteHeader(http.StatusCreated)
	xml.NewEncoder(w).Encode(p) // Respond with the created product
}

func main() {
	// Route /products to GET and POST handlers
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getProducts(w, r)
		} else if r.Method == http.MethodPost {
			addProduct(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// Start the HTTP server on port 8081
	http.ListenAndServe(":8081", nil)
}
