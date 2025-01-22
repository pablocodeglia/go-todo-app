package apiv1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	types "todoapp/types"
	"github.com/google/uuid"
)

func RegisterApiHandlers(mux *http.ServeMux) {
	mux.Handle("/api/v1/todo/new/{userId}", &TodoCreateHandler{})
	mux.Handle("/api/v1/todo/{userId}", &TodosHandler{})
	mux.Handle("/api/v1/todo/{userId}/{todoId}", &TodosHandler{})
}

type TodosHandler struct{}
type TodoCreateHandler struct{}

func (h *TodoCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		h.SaveUserChanges(w, r)
	}
}

func (h *TodosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		h.GetTodosByUser(w, r)
	case r.Method == http.MethodDelete:
		h.DeleteTodo(w, r)
	case r.Method == http.MethodPost:
		h.CreateTodo(w, r)
	case r.Method == http.MethodPut:
		h.UpdateTodo(w, r)
	}
}

func (h *TodosHandler) GetTodosByUser(w http.ResponseWriter, r *http.Request) types.TodoStoreData {
	userId := r.PathValue("userId")
	var todos types.TodoStoreData

	file, err := os.Open(fmt.Sprintf("./data/%s.json", userId))
	if err != nil {
		FileNotFoundErrorHandler(w)
	} else {
		byteValue, _ := io.ReadAll(file)

		json.Unmarshal(byteValue, &todos)

		w.WriteHeader(http.StatusOK)
		w.Write(byteValue)
	}
	return todos
}

func (h *TodosHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")

	var newTodo types.Todo

	dec := json.NewDecoder(r.Body)

	// validate json
	err := dec.Decode(&newTodo)

	if err != nil {
		fmt.Printf("%s", err.Error())
		BadRequestErrorHandler(w, fmt.Sprintf("%s", err))
	}

	// open file
	file, err := os.Open(fmt.Sprintf("./data/%s.json", userId))
	if err != nil {
		FileNotFoundErrorHandler(w)
	}

	byteValue, _ := io.ReadAll(file)

	var userData types.TodoStoreData
	json.Unmarshal([]byte(byteValue), &userData)

	// add new todo to file
	newTodoUUID := uuid.New().String()
	newTodoObject := map[string]types.Todo{newTodoUUID: {Task: newTodo.Task, IsDone: newTodo.IsDone, CreatedAt: newTodo.CreatedAt}}
	userData.Data = append(userData.Data, newTodoObject)

	// save new file
	saveUserJsonFile(w, userId, userData)
	w.WriteHeader(http.StatusOK)
}

func (h *TodoCreateHandler) SaveUserChanges(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	byteValue, _ := io.ReadAll(r.Body)

	os.WriteFile(fmt.Sprintf("data/%s.json", userId), byteValue, os.ModePerm)
	w.WriteHeader(http.StatusOK)
}

func (h *TodosHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	todoId := r.PathValue("todoId")

	// open file and get existing data
	file, err := os.Open(fmt.Sprintf("./data/%s.json", userId))
	if err != nil {
		FileNotFoundErrorHandler(w)
	}

	byteValue, _ := io.ReadAll(file)
	file.Close()

	var userData types.TodoStoreData
	json.Unmarshal([]byte(byteValue), &userData)
	toUpdateIndex := findIndexByTodoIdFunc(userData.Data, todoId)

	// decode request body
	// var updatedTodo clistore.Todo

	updatedTodo := userData.Data[toUpdateIndex][todoId]
	json.NewDecoder(r.Body).Decode(&updatedTodo)

	// update todo
	userData.Data[toUpdateIndex][todoId] = updatedTodo

	// save file
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
	file.Close()

	var userData types.TodoStoreData
	json.Unmarshal([]byte(byteValue), &userData)

	// get toDeleteIndex of item to be deleted
	toDeleteIndex := findIndexByTodoIdFunc(userData.Data, todoId)
	// delete item or return not found error
	if toDeleteIndex == -1 {
		FileNotFoundErrorHandler(w)
	} else {
		userData.Data = append(userData.Data[:toDeleteIndex], userData.Data[toDeleteIndex+1:]...)
		saveUserJsonFile(w, userId, userData)
		w.WriteHeader(http.StatusOK)
	}

}
