package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"todo-app/data"
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

func session(writer http.ResponseWriter, request *http.Request) (sess data.Session, err error) {
	cookie, err := request.Cookie("_cookie")

	if err == nil {
		sess = data.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("invalid session")
		}
	}

	return
}
