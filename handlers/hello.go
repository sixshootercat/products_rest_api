package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Hello is a simple handler
type Hello struct {
	l *log.Logger
}

// NewHello creates a new hello handler with the given logger
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/http/#Handler
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello world")

	// read the body
	body, err := ioutil.ReadAll(r.Body)

	// handle error if any
	if err != nil {
		http.Error(rw, ":(", http.StatusBadRequest)
		return
	}

	// write the response
	fmt.Fprintf(rw, "Hello %s", body)
}