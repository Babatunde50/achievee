package main

import (
	// "encoding/json"
	// "log"
	"encoding/json"
	"net/http"
	"todo-app/data"

	// "todo-app/data"

	"github.com/julienschmidt/httprouter"
)

func createGoal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	goal := data.Goal{
		TotalPercentCompleted: 0,
		Completed:             false,
		Paused:                false,
		UserId:                ctx.Value(userIdKey).(int),
	}

	err := json.NewDecoder(r.Body).Decode(&goal)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusBadRequest)
		return
	}

	// TODO: do the needful

	resp := message(true, "Successful")
	resp["data"] = goal
}

func userGoals(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: do the needful
}

func deleteGoal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: do the needful
}

func updateGoal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: do the needful
}

func updateGoalProgress(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: do the needful
}
