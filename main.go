package main

import (
	"fmt"
	"net/http"
	// "time"
)

func main() {
	mux := http.NewServeMux()

	// TODO: Handle static files

	// all routes patterns matched here

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello! This is the index page!!!!")
	})

	// server configuration and start
	server := &http.Server{
		Addr:    "localhost:8081",
		Handler: mux,
		// ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		// WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
}
