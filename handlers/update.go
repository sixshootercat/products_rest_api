package handlers

import (
	"net/http"
	"product_api/data"
)

func (p Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	id := GetProductID(r)

	p.l.Println("handle PUT product", id)
	product := r.Context().Value(KeyProduct{}).(data.Product) // cast to data.Product

	err := data.UpdateProduct(id, &product)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
}