// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"product_api/data"
	"strconv"

	"github.com/gorilla/mux"
)

// A list of products returns in the response
// swagger:response productsResponse
type productsResponseWrapper struct {
	// All products in the datastore
	// in: body
	Body []data.Product
}

// swagger:parameters deleteProduct
type productIDParameterWrapper struct {
	// The id of the product to delete from the database
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:response noContent
type productsNoContent struct {

}

// Products is an http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

type KeyProduct struct{}

func(p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		product := data.Product{}

		// deserialize the product
		err := product.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product")
			http.Error(rw, "Error readig product", http.StatusBadRequest)
			return
		}

		// validate the product
		err = product.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product")
			http.Error(
				rw,
				fmt.Sprintf("Incorrect product field: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		// calls the next handler which can be another middleware in the chain of the final handler
		next.ServeHTTP(rw, r)
	})
}

func GetProductID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}