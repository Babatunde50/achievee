package main

import (
	"encoding/json"
	"net/http"
	"todo-app/data"

	"github.com/julienschmidt/httprouter"
)

// POST -> signup
func signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := data.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respond(w, message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	// check if email doesn't already exists
	emailExists, err := data.UserEmailExists(user.Email)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusBadRequest)
		return
	}

	if emailExists {
		respond(w, message(false, "A user with this email already exists"), http.StatusBadRequest)
		return
	}

	// check if password is atleast 6 characters.. then
	if len(user.Password) < 6 {
		respond(w, message(false, "Password must be atleast 6 characters"), http.StatusBadRequest)
		return
	}

	// hashed password and replace user password with the hashed password
	hashedPassword, err := data.GeneratePassword(user.Password)

	if err != nil {
		// log error
		panic(err)
	}

	user.Password = hashedPassword

	err = user.Create()

	if err != nil {
		//log error
		panic(err)
	}

	user.Password = ""

	resp := message(true, "Successful")
	resp["data"] = user

	respond(w, resp, http.StatusCreated)
}

// TODO: Login
