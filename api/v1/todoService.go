package apiv1

import (
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

func (service *TodoService) DeleteTodo() {
	//TODO
}

func (service *TodoService) GetAllByUserId() {
	//TODO
}
