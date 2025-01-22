package clistore

import (
	"fmt"
	"net/http"
	"testing"
	apiv1 "todoapp/api/v1"
	httpV1 "todoapp/http"
)

func TestAdd(t *testing.T) {
	//initialize server for Api endpoints
	mux := http.NewServeMux()
	apiv1.RegisterApiHandlers(mux)
	httpV1.RegisterHttpHandlers(mux)
	// log.Println("Server is running on port 8080...")
	go http.ListenAndServe("localhost:8080", mux)

	//
	store := NewStore()
	store.CurrentUserId = "testuser"
	// log user
	store.ClearCache()
	go store.LoadData()

	fmt.Println(store.CurrentUserId)

	store.Add("TEST NEW TASK")
	store.Add("TEST ANOTHER NEW TASK")

	got := len(store.Data)
	want := 2

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}

}
