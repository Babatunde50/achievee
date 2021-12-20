package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"todo-app/data"

	"github.com/julienschmidt/httprouter"
)

func userTasks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	ctx := r.Context()

	tasks, err := data.GetTasksByUserId(ctx.Value(userIdKey).(int))

	if err != nil {
		log.Fatal(err)
	}

	resp := message(true, "Successful")
	resp["data"] = tasks

	respond(w, resp, http.StatusCreated)

}

func updateCompleteSubTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	complete := struct {
		Completed bool `json:"completed"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&complete)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusBadRequest)
		return
	}

	taskId := ps.ByName("taskId")
	id := ps.ByName("id")

	err = data.UpdateSubTaskCompletion(complete.Completed, id, taskId)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	resp := message(true, "Successful")
	resp["data"] = complete

	respond(w, resp, http.StatusOK)

}

func updateCompleteTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	complete := struct {
		Completed bool `json:"completed"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&complete)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusBadRequest)
		return
	}

	taskIdParam := ps.ByName("id")

	taskId, _ := strconv.Atoi(taskIdParam)

	completeTaskData := data.Completed{
		IsCompleted: complete.Completed,
		TaskId:      taskId,
		UserId:      ctx.Value(userIdKey).(int),
	}

	err = data.UpdateTaskCompletion(completeTaskData)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	resp := message(true, "Successful")
	resp["data"] = complete

	respond(w, resp, http.StatusOK)
}

func updateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	ctx := r.Context()

	d := struct {
		Title    string    `json:"title"`
		ColorTag string    `json:"colorTag"`
		Deadline time.Time `json:"deadline"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&d)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusBadRequest)
		return
	}

	taskIdParam := ps.ByName("id")

	taskId, _ := strconv.Atoi(taskIdParam)

	updateData := data.UpdateTaskData{
		Title:    d.Title,
		ColorTag: d.ColorTag,
		UserId:   ctx.Value(userIdKey).(int),
		TaskId:   taskId,
		Deadline: d.Deadline,
	}

	err = data.UpdateTaskEdits(updateData)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	resp := message(true, "Successful")
	resp["data"] = updateData

	respond(w, resp, http.StatusOK)

}

func deleteTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	ctx := r.Context()

	taskIdParam := ps.ByName("id")

	taskId, err := strconv.Atoi(taskIdParam)

	if err != nil {
		// respond(w, message(false, err.Error()), http.StatusUnauthorized)
		return
	}

	task := data.Task{
		Id:     taskId,
		UserId: ctx.Value(userIdKey).(int),
	}

	err = task.Delete()

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	resp := message(true, "Successful")

	respond(w, resp, http.StatusNoContent)
}

func deleteSubTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	taskId := ps.ByName("taskId")
	id := ps.ByName("id")

	err := data.DeleteSubTask(id, taskId)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	resp := message(true, "Successful")

	respond(w, resp, http.StatusOK)
}

func createSubTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	subTask := data.SubTask{}

	err := json.NewDecoder(r.Body).Decode(&subTask)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusBadRequest)
		return
	}

	taskIdParam := ps.ByName("taskId")

	taskId, err := strconv.Atoi(taskIdParam)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusUnauthorized)
		return
	}

	subTask.TaskId = taskId

	// create a subTask
	err = subTask.Create()

	if err != nil {
		log.Fatal(err.Error())
	}

	resp := message(true, "Successful")
	resp["data"] = subTask

	respond(w, resp, http.StatusCreated)
}

func createTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	ctx := r.Context()

	task := data.Task{}

	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusBadRequest)
		return
	}

	task.UserId = ctx.Value(userIdKey).(int)

	// create a task
	err = task.Create()

	if err != nil {
		log.Fatal(err.Error())
	}

	resp := message(true, "Successful")
	resp["data"] = task

	respond(w, resp, http.StatusCreated)
}
