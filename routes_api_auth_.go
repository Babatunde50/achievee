package main

import (
	"encoding/json"
	"net/http"
	"time"
	"todo-app/data"

	"github.com/julienschmidt/httprouter"
)

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

// POST -> Login
func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	loginData := LoginData{}

	err := json.NewDecoder(r.Body).Decode(&loginData)

	if err != nil {
		respond(w, message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	user, err := data.UserByEmail(loginData.Email)

	if err != nil {
		respond(w, message(false, "Invalid email or password"), http.StatusBadRequest)
		return
	}

	err = data.ComparePasswords(user.Password, loginData.Password)

	if err != nil {
		respond(w, message(false, "Invalid email or password"), http.StatusBadRequest)
		return
	}

	// create a new random session token
	sessionToken := createUUID()

	_, err = cache.Do("SETEX", sessionToken, "300", user.Email)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HttpOnly: true,
		Expires:  time.Now().Add(7 * time.Hour),
		Path:     "/",
	}

	http.SetCookie(w, &cookie)

	resp := message(true, "Successful")
	resp["data"] = user

	respond(w, resp, http.StatusOK)
}

func refreshToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionToken, userEmail, err := session(w, r)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusUnauthorized)
		return
	}

	newSessionToken := createUUID()

	_, err = cache.Do("SETEX", newSessionToken, "300", userEmail)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	// Delete the older session token
	_, err = cache.Do("DEL", sessionToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(300 * time.Second),
	})
}

func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionToken, _, err := session(w, r)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusUnauthorized)
		return
	}

	// Delete the session from cache
	_, err = cache.Do("DEL", sessionToken)

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, &cookie)

	respond(w, message(true, "logout successs"), http.StatusOK)
}
