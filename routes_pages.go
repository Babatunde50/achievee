package main

import (
	// "fmt"
	"fmt"
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

	goals, err := data.GetGoalsByUserId(ctx.Value(userIdKey).(int))

	if err != nil {
		log.Fatal(err)
	}

	user, err := data.UserByEmail(ctx.Value(userEmailKey).(string))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user, "user..")

	d := struct {
		Tasks []data.TaskWithSubTasks
		Goals []data.Goal
		User  data.User
	}{
		Tasks: tasks,
		Goals: goals,
		User:  user,
	}

	generateHTML(w, d, "layout", "planner")
}

// TODO:  GET -> Error page
