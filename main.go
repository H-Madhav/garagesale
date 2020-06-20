package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//==========================simple server implementation==============================

// func main() {
// 	// Convert the Echo function to a type that implements http.Handler
// 	h := http.HandlerFunc(Echo)
// 	fmt.Println("server is running on localhost: 8000")
// 	if err := http.ListenAndServe("localhost:8000", h); err != nil {
// 		log.Fatalf("error: listening and serving: %s", err)
// 	}
// }

// // Echo is a basic HTTP Handler.
// func Echo(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "You asked to %s %s\n", r.Method, r.URL.Path)
// }

//======================================================================================

// ========================server implementaion with gracefull shutdown=================

func main() {

	log.Printf("main: started")
	defer log.Println("main : Completed")

	api := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(ListProducts),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	serverError := make(chan error, 1)

	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverError <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Fatalf("error: listening and serving: %s", err)
	case <-shutdown:
		log.Println("main : Start shutdown")

		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}

//Product schema
type Product struct {
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Quantity int    `json:"quantity"`
}

// ListProducts is a basic HTTP Handler.
func ListProducts(w http.ResponseWriter, r *http.Request) {
	list := []Product{}

	if true {
		list = append(list, Product{Name: "Iphone", Cost: 29000, Quantity: 1})
		list = append(list, Product{Name: "McDonalds Toys", Cost: 75, Quantity: 120})
	}

	data, err := json.Marshal(list)

	if err != nil {
		log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error writing result", err)
	}
}
