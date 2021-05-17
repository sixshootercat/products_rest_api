package handlers

import (
	"net/http"
	"product_api/data"
)

func(p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle POST product")

	product := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&product)
}