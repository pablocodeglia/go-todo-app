package main

import (
	"net/http"
	apiv1 "todoapp/api/v1"
	"todoapp/cli"
	httpV1 "todoapp/http"
	clistore "todoapp/store/cli-store"
)

func NewStore() *clistore.TodoStore {
	return &clistore.TodoStore{
		Data: make(map[string]clistore.Todo),
	}
}

func main() {
	mux := http.NewServeMux()
	apiv1.RegisterApiHandlers(mux)
	httpV1.RegisterHttpHandlers(mux)
	// log.Println("Server is running on port 8080...")
	go http.ListenAndServe("localhost:8080", mux)

	cli.Clr()
	store := NewStore()
	store.LogUser()

}
