package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	session, err := user.CreateSession()

	if err != nil {
		respond(w, message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    session.Uuid,
		MaxAge:   7000,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	}

	http.SetCookie(w, &cookie)

	resp := message(true, "Successful")
	resp["session"] = session

	respond(w, resp, http.StatusOK)
}

func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie, err := r.Cookie("_cookie")

	fmt.Println(cookie, "cookie", err.Error(), "error")

	if err != http.ErrNoCookie {
		session := data.Session{Uuid: cookie.Value}
		err = session.DeleteByUUID()

		if err != nil {
			respond(w, message(false, err.Error()), http.StatusInternalServerError)
			return
		}

		respond(w, message(true, "Logout successfully"), http.StatusInternalServerError)
	}

	respond(w, message(false, err.Error()), http.StatusBadRequest)
}
