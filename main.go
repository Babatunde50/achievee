package main

import (
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var cache redis.Conn

func main() {
	router := httprouter.New()

	// init cache
	initCache()

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
	router.POST("/api/refresh-token", refreshToken)

	http.ListenAndServe(":8081", cors.Default().Handler(router))
}

func initCache() {
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	// Assign the connection to the package level `cache` variable
	cache = conn
}
