package apiv1

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	clistore "todoapp/store/cli-store"
)

type TodoService struct{}

type TodoServiceInterface interface {
	CreateTodo()
	UpdateTodo()
	DeleteTodo()
	GetAllByUserId()
}

func NewTodoService() *TodoService {
	return &TodoService{}
}

func (service *TodoService) CreateTodo(userId string, newTodo clistore.Todo) clistore.TodoStoreData {
	todos := clistore.TodoStoreData{}
	return todos
	//TODO
}

func (service *TodoService) UpdateTodo() {

	//TODO
}

func (service *TodoService) DeleteTodo(userId, todoId string) {
	fmt.Println("userId :", userId)
	fmt.Println("todoId :", todoId)
	// currentTodosResponse := service.GetAllByUserId(userId)

	// currentTodosData := Json.Unmarshal(currentTodosResponse, )

	// url := fmt.Sprintf("http://127.0.0.1:8080/api/v1/todo/new/%s", s.CurrentUserId)

	// req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBytes))
	// if err != nil {
	// 	fmt.Print(err.Error())
	// }
	//TODO
}

func (service *TodoService) GetAllByUserId(userId string) []map[string]clistore.Todo {
	url := fmt.Sprintf("http://127.0.0.1:8080/api/v1/todo/%s", userId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Print(err.Error())
	}

	var userData clistore.TodoStoreData
	json.Unmarshal(body, &userData)

	//TODO
	response := []map[string]clistore.Todo{}
	response = userData.Data
	return response
}
