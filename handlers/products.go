package handlers

import (
	"http_server/product-api/data"
	"log"
	"net/http"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle get
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// handle update


	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}