package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	// Handles static files
	router.ServeFiles("/static/*filepath", http.Dir("public"))

	// all routes patterns matched here
	router.GET("/", index)

	http.ListenAndServe(":8081", router)
}
