package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GET -> Index Page
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	generateHTML(w, nil, "layout", "index")
}

// PLANNER -> authenticated index page
func planner(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := r.Cookie("session_token")

	if err != nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	generateHTML(w, nil, "layout", "planner")
}

// TODO:  GET -> Error page
