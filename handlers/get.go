package handlers

import (
	"net/http"
	"product_api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 200: productsResponse

// GetProducts returns the products from the data store
func(p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle GET products")

	// fetch the products from the data store
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}