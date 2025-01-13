package main

import (
	"todoapp/cli"
	clistore "todoapp/store/cli-store"
)

func NewStore() *clistore.TodoStore {
	return &clistore.TodoStore{
		Data: make(map[string]clistore.Todo),
	}
}

func main() {
	cli.Clr()
	store := NewStore()
	store.LogUser()
}
