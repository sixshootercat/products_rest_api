package handlers

import (
	"net/http"
	"product_api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Deletes a product from the database
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// RemoveProduct handles DELETE requests and removes items from the database
func(p *Products) RemoveProduct(rw http.ResponseWriter, r *http.Request) {
	id := GetProductID(r)

	p.l.Println("handle DELETE product", id)

	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] deleting record id does not exist")
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
}