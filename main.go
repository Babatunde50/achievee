package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	router := httprouter.New()

	// Handles static files
	router.ServeFiles("/static/*filepath", http.Dir("public"))

	// all routes patterns matched here

	// PAGES
	router.GET("/", index)
	router.GET("/planner", planner)

	// APIS

	// AUTH
	router.POST("/api/signup", signup)
	router.POST("/api/login", login)
	router.POST("/api/logout", logout)

	http.ListenAndServe(":8081", cors.Default().Handler(router))
}
