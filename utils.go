package main

import (
	"crypto/rand"
	"encoding/json"
	"log"

	// "errors"
	"fmt"
	"html/template"
	"net/http"
	// "todo-app/data"
)

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, "layout", data)
}

func message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

func respond(w http.ResponseWriter, data map[string]interface{}, statusCode int) {
	w.Header().Add("Content-Type", "application/json")

	// cookie := http.Cookie{
	// 	Name:     "_test_cookie",
	// 	Value:    "123123",
	// 	HttpOnly: true,
	// }

	w.WriteHeader(statusCode)
	// http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(data)
}

func session(writer http.ResponseWriter, request *http.Request) (sessionToken string, userEmail string, err error) {
	cookie, err := request.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			respond(writer, message(false, err.Error()), http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		respond(writer, message(false, err.Error()), http.StatusBadRequest)
		return
	}

	sessionToken = cookie.Value

	response, err := cache.Do("GET", sessionToken)

	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		respond(writer, message(false, err.Error()), http.StatusInternalServerError)

		return
	}

	userEmail = fmt.Sprintf("%s", response)

	if response == nil {
		// If the session token is not present in cache, return an unauthorized error
		respond(writer, message(false, err.Error()), http.StatusUnauthorized)
		return
	}

	return
}
