package main

import (
	// "log"
	"log"
	"net/http"
	apiv1 "todoapp/api/v1"
	"todoapp/cli"
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
	// log.Println("Server is running on port 8080...")
	go http.ListenAndServe("localhost:8080", mux)

	cli.Clr()
	store := NewStore()
	store.LogUser()

}
