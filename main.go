package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"product_api/handlers"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// create the handlers
	ph := handlers.NewProducts(l)

	// create new serve mux and register the handlers
	sm := mux.NewRouter()

	// get router
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	// put router
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)


	// post router
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)


	// create a new server
	server := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		ErrorLog:     l,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// start the server
	go func() {
		l.Println("listening on port :9090...")
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// trap sigterm and interrupt or gracefully shutdown server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	cancel()
	server.Shutdown(tc)
}
