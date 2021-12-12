package main

import (
	// "fmt"
	// "fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GET -> Index Page
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	generateHTML(w, nil, "layout", "index")
}

// PLANNER -> authenticated index page
func planner(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// _, err := session(w, r)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	// }

	generateHTML(w, nil, "layout", "planner")
}

// TODO:  GET -> Error page
