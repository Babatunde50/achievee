package main

import (
	// "fmt"
	"log"
	"net/http"
	"todo-app/data"

	"github.com/julienschmidt/httprouter"
)

// GET -> Index Page
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	generateHTML(w, nil, "layout", "index")
}

// PLANNER -> authenticated index page
func planner(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	ctx := r.Context()

	// get tasks
	tasks, err := data.GetTasksByUserId(ctx.Value(userIdKey).(int))

	if err != nil {
		log.Fatal(err)
	}

	generateHTML(w, tasks, "layout", "planner")
}

// TODO:  GET -> Error page
