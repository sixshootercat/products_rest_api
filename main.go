package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"product_api/data"
	"product_api/handlers"
	"syscall"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "products-api ", log.LstdFlags)
	v := data.NewValidation()

	// create the handlers
	ph := handlers.NewProducts(l, v)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	getReqs := sm.Methods(http.MethodGet).Subrouter()
	getReqs.HandleFunc("/products", ph.ListAll)
	getReqs.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putReqs := sm.Methods(http.MethodPut).Subrouter()
	putReqs.HandleFunc("/products", ph.Update)
	putReqs.Use(ph.MiddlewareValidateProduct)

	postReqs := sm.Methods(http.MethodPost).Subrouter()
	postReqs.HandleFunc("/products", ph.Create)
	postReqs.Use(ph.MiddlewareValidateProduct)

	deleteReqs := sm.Methods(http.MethodDelete).Subrouter()
	deleteReqs.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getReqs.Handle("/docs", sh)
	getReqs.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// create a new server
	server := http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := server.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	cancel()
	server.Shutdown(ctx)
}