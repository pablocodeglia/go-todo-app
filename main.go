package main

import (
	"log"
	"net/http"

	// "todoapp/cli"
	apiv1 "todoapp/api/v1"
	clistore "todoapp/store/cli-store"
)

func NewStore() *clistore.TodoStore {
	return &clistore.TodoStore{
		Data: make(map[string]clistore.Todo),
	}
}

func main() {
	// cli.Clr()
	// store := NewStore()
	// store.LogUser()
	mux := http.NewServeMux()
	apiv1.RegisterApiHandlers(mux)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
