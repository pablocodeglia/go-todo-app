package clistore

import (
	"fmt"
	"net/http"
	"testing"
	apiv1 "todoapp/api/v1"
	httpV1 "todoapp/http"
)

func initializeApi() {
	mux := http.NewServeMux()
	apiv1.RegisterApiHandlers(mux)
	httpV1.RegisterHttpHandlers(mux)
	go http.ListenAndServe("localhost:8080", mux)
}

func TestAdd(t *testing.T) {
	//initialize server for Api endpoints
	initializeApi()

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

func BenchmarkClearCache(b *testing.B) {
	store := NewStore()

	for i := 0; i < b.N; i++ {
		store.ClearCache()
	}
}

func BenchmarkListTodos(b *testing.B) {
	store := NewStore()

	for i := 0; i < b.N; i++ {
		store.ListTodos()
	}
}
