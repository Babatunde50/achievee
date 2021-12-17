package main

import (
	"context"
	"fmt"
	"net/http"
	"todo-app/data"

	"github.com/gomodule/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var cache redis.Conn

type key int

const (
	sessionTokenKey key = iota
	userIdKey       key = iota
	userEmailKey    key = iota
)

func pagesAuthMiddleware(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		cookie, err := r.Cookie("session_token")

		if err != nil {
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}

		sessionToken := cookie.Value

		response, _ := cache.Do("GET", sessionToken)

		if response == nil {
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}

		userEmail := fmt.Sprintf("%s", response)

		user, err := data.UserByEmail(userEmail)

		if err != nil {
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}

		ctx := r.Context()

		// Get new context with key-value "settings"
		ctx = context.WithValue(ctx, sessionTokenKey, sessionToken)

		ctx = context.WithValue(ctx, userIdKey, user.Id)
		ctx = context.WithValue(ctx, userEmailKey, user.Email)

		r = r.WithContext(ctx)

		n(w, r, ps)
	}
}

// middleware is used to intercept incoming HTTP calls and apply general functions upon them.
func apiAuthMiddleware(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// do some authentication

		cookie, err := r.Cookie("session_token")

		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				respond(w, message(false, err.Error()), http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			respond(w, message(false, err.Error()), http.StatusBadRequest)
			return
		}

		sessionToken := cookie.Value

		response, err := cache.Do("GET", sessionToken)

		if err != nil {
			// If there is an error fetching from cache, return an internal server error status
			respond(w, message(false, err.Error()), http.StatusInternalServerError)

			return
		}

		if response == nil {
			respond(w, message(false, "Unauthenticated"), http.StatusUnauthorized)
			return
		}

		userEmail := fmt.Sprintf("%s", response)

		user, err := data.UserByEmail(userEmail)

		if err != nil {
			respond(w, message(false, err.Error()), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()

		// Get new context with key-value "settings"
		ctx = context.WithValue(ctx, sessionTokenKey, sessionToken)

		ctx = context.WithValue(ctx, userIdKey, user.Id)
		ctx = context.WithValue(ctx, userEmailKey, user.Email)

		r = r.WithContext(ctx)

		n(w, r, ps)
	}
}

func main() {
	router := httprouter.New()

	// init cache
	initCache()

	// Handles static files
	router.ServeFiles("/static/*filepath", http.Dir("public"))

	// all routes patterns matched here

	// PAGES
	router.GET("/", index)
	router.GET("/planner", pagesAuthMiddleware(planner))

	// APIS

	// AUTH
	router.POST("/api/signup", signup)
	router.POST("/api/login", login)
	router.POST("/api/logout", apiAuthMiddleware(logout))
	router.POST("/api/refresh-token", apiAuthMiddleware(refreshToken))

	// TASKS
	router.POST("/api/tasks", apiAuthMiddleware(createTask))
	router.GET("/api/tasks", apiAuthMiddleware(userTasks))
	router.DELETE("/api/tasks/:id", apiAuthMiddleware(deleteTask))
	router.PATCH("/api/tasks/:id/edits", apiAuthMiddleware(updateTask))
	router.PATCH("/api/tasks/:id/completed", apiAuthMiddleware(updateCompleteTask))

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
