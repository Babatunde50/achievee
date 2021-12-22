package main

import (
	// "encoding/json"
	// "log"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"todo-app/data"

	// "todo-app/data"

	"github.com/julienschmidt/httprouter"
)

func createGoal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	goal := data.Goal{
		TotalPercentCompleted: 0,
		TotalPercent:          100,
		Completed:             false,
		Paused:                false,
		UserId:                ctx.Value(userIdKey).(int),
	}

	err := json.NewDecoder(r.Body).Decode(&goal)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusBadRequest)
		return
	}

	err = goal.Create()

	if err != nil {
		log.Fatal(err.Error())
	}

	resp := message(true, "Successful")
	resp["data"] = goal

	respond(w, resp, http.StatusCreated)
}

func userGoals(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	goals, err := data.GetGoalsByUserId(ctx.Value(userIdKey).(int))

	if err != nil {
		log.Fatal(err)
	}

	resp := message(true, "Successful")
	resp["data"] = goals

	respond(w, resp, http.StatusCreated)
}

func deleteGoal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	ctx := r.Context()

	goalIdParam := ps.ByName("id")

	goalId, err := strconv.Atoi(goalIdParam)

	if err != nil {
		// respond(w, message(false, err.Error()), http.StatusUnauthorized)
		return
	}

	goal := data.Goal{
		Id:     goalId,
		UserId: ctx.Value(userIdKey).(int),
	}

	err = goal.Delete()

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	resp := message(true, "Successful")

	respond(w, resp, http.StatusNoContent)

}

func updateGoal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	d := struct {
		Title        string    `json:"title"`
		ColorTag     string    `json:"colorTag"`
		Deadline     time.Time `json:"deadline"`
		TotalPercent int       `json:"totalPercent"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&d)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusBadRequest)
		return
	}

	goalIdParam := ps.ByName("id")

	goalId, _ := strconv.Atoi(goalIdParam)

	updateData := data.UpdateGoalData{
		UserId:       ctx.Value(userIdKey).(int),
		GoalId:       goalId,
		Title:        d.Title,
		ColorTag:     d.ColorTag,
		TotalPercent: d.TotalPercent,
		Deadline:     d.Deadline,
	}

	err = data.UpdateGoal(updateData)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	resp := message(true, "Successful")
	resp["data"] = updateData

	respond(w, resp, http.StatusOK)

}

func updateGoalProgress(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	tcp := struct {
		TotalPercentCompleted int `json:"totalPercentCompleted"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&tcp)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusBadRequest)
		return
	}

	goalIdParam := ps.ByName("id")

	goalId, _ := strconv.Atoi(goalIdParam)

	d := data.UpdateGoalPercentCompletedData{
		UserId:                ctx.Value(userIdKey).(int),
		GoalId:                goalId,
		TotalPercentCompleted: tcp.TotalPercentCompleted,
	}

	err = data.UpdateGoalProgress(d)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	resp := message(true, "Successful")
	resp["data"] = d

	respond(w, resp, http.StatusOK)
}
