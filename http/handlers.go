package httpV1

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	apiv1 "todoapp/api/v1"
)

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func RegisterHttpHandlers(mux *http.ServeMux) {
	fileServer := http.FileServer(http.Dir("./static/"))
	// mux.HandleFunc("/", HandleRootView)
	mux.HandleFunc("/users", HandleUsersView)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/todos/{userId}", HandleTodosView)

}

func HandleRootView(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func HandleUsersView(w http.ResponseWriter, r *http.Request) {
	template, err := template.New("users.html").ParseFiles("./webapp/users.html")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("./data/users.json")
	if err != nil {
		log.Fatal(err)
	}
	byteValue, _ := io.ReadAll(file)

	var users map[string]User
	json.Unmarshal(byteValue, &users)

	err = template.Execute(w, users)
	if err != nil {
		log.Fatal(err)
	}

}

func HandleTodosView(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	service := apiv1.NewTodoService()

	funcMap := template.FuncMap{}

	template, err := template.New("index.html").Funcs(funcMap).ParseFiles("./webapp/index.html")
	if err != nil {
		log.Fatal(err)
	}

	todosData := service.GetAllByUserId(userId)

	err = template.Execute(w, map[string]interface{}{"data": todosData, "userId": userId})
	if err != nil {
		log.Fatal(err)
	}

}
