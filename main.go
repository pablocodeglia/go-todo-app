package main

import (
	"net/http"
	cli "todoapp/cli"
	apiv1 "todoapp/api/v1"
	httpV1 "todoapp/http"
	clistore "todoapp/store/cli-store"
)

func main() {
	go initServer()
	initCliStore()
}

func initServer() {
	mux := http.NewServeMux()
	apiv1.RegisterApiHandlers(mux)
	httpV1.RegisterHttpHandlers(mux)
	// log.Println("Server is running on port 8080...")
	http.ListenAndServe("localhost:8080", mux)
}

func initCliStore() {
	store := clistore.NewStore()

	userId := cli.GetUserInput("Username: ")
	store.LogUser(userId)
	store.ListTodos()
	store.DisplayOptions()
}
