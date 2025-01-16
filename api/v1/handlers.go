package apiv1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	clistore "todoapp/store/cli-store"

	"github.com/google/uuid"
)

func RegisterApiHandlers(mux *http.ServeMux) {
	mux.Handle("/api/v1/{userId}", &TodosHandler{})
	mux.Handle("/api/v1/{userId}/{todoId}", &TodosHandler{})
}

type TodosHandler struct{}

func (h *TodosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		h.GetTodosByUser(w, r)
	case r.Method == http.MethodDelete:
		h.DeleteTodo(w, r)
	case r.Method == http.MethodPost:
		h.CreateTodo(w, r)
	}
}

func (h *TodosHandler) GetTodosByUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")

	file, err := os.Open(fmt.Sprintf("./data/%s.json", userId))
	if err != nil {
		FileNotFoundErrorHandler(w)
	} else {

		byteValue, _ := io.ReadAll(file)

		w.WriteHeader(http.StatusOK)
		w.Write(byteValue)
	}
}

func (h *TodosHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")

	// new todo
	// newTodo := clistore.Todo{}
	var newTodo clistore.Todo

	// validate json
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&newTodo)
	if err != nil {
		BadRequestErrorHandler(w)
	}

	// open file
	file, err := os.Open(fmt.Sprintf("./data/%s.json", userId))
	if err != nil {
		FileNotFoundErrorHandler(w)
	}

	byteValue, _ := io.ReadAll(file)

	var userData clistore.TodoStoreData
	json.Unmarshal([]byte(byteValue), &userData)

	// add new todo to file
	newTodoUUID := uuid.New().String()
	newTodoObject := map[string]clistore.Todo{newTodoUUID: {Task: newTodo.Task, IsDone: newTodo.IsDone, CreatedAt: newTodo.CreatedAt}}
	userData.Data = append(userData.Data, newTodoObject)
	
	// save new file
	saveUserJsonFile(w, userId, userData)
	w.WriteHeader(http.StatusOK)
}

func (h *TodosHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	todoId := r.PathValue("todoId")

	file, err := os.Open(fmt.Sprintf("./data/%s.json", userId))
	if err != nil {
		FileNotFoundErrorHandler(w)
	}

	byteValue, _ := io.ReadAll(file)

	var userData clistore.TodoStoreData
	json.Unmarshal([]byte(byteValue), &userData)

	// get toDeleteIndex of item to be deleted
	toDeleteIndex := slices.IndexFunc(userData.Data, func(data map[string]clistore.Todo) bool {
		for k := range data {
			if k == todoId {
				return true
			}
		}
		return false
	})
	// delete item or return not found error
	if toDeleteIndex == -1 {
		FileNotFoundErrorHandler(w)
	} else {
		userData.Data = append(userData.Data[:toDeleteIndex], userData.Data[toDeleteIndex+1:]...)
		saveUserJsonFile(w, userId, userData)
		w.WriteHeader(http.StatusOK)
	}

	file.Close()
}
