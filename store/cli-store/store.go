package clistore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	cli "todoapp/cli"
	types "todoapp/types"

	"github.com/google/uuid"
)

type TodoStore struct {
	CurrentUserId string
	Data          map[string]types.Todo
	Mu            sync.Mutex
}

func NewStore() *TodoStore {
	return &TodoStore{
		Data: make(map[string]types.Todo),
	}
}

func (s *TodoStore) Add(newTaskName string) {

	taskUuid := uuid.New()

	s.Mu.Lock()
	s.Data[taskUuid.String()] = types.Todo{
		Task:      newTaskName,
		IsDone:    false,
		CreatedAt: time.Now(),
	}
	s.Mu.Unlock()

	cli.Clr()
	println("\n- Success! Task added.\n")
}

func (s *TodoStore) Delete(todoId string) {
	s.Mu.Lock()
	delete(s.Data, todoId)
	s.Mu.Unlock()

}

func (s *TodoStore) MarkAsDone(todoId string) {
	s.Mu.Lock()

	if entry, ok := s.Data[todoId]; ok {
		entry.IsDone = true
		s.Data[todoId] = entry
	}

	s.Mu.Unlock()

	cli.Clr()
	fmt.Printf("\n- Success! Todo marked as done.\n")

}

func (s *TodoStore) ListTodos() {
	fmt.Printf("\n######### TODOS ##########\n")
	i := 1

	s.Mu.Lock()
	for k, v := range s.Data {
		fmt.Printf("- Task %d (Id: %s)\n", i, k)
		fmt.Printf("\t- %s\n", v.Task)
		var isDone string
		if v.IsDone {
			isDone = "[🗸]"
		} else {
			isDone = "[ ]"
		}
		fmt.Printf("\t- Done? %s\n", isDone)
		i++
	}
	fmt.Println("")
	s.Mu.Unlock()

}

func (s *TodoStore) ClearCache() {
	s.Mu.Lock()

	s.Data = map[string]types.Todo{}

	s.Mu.Unlock()
}

func (s *TodoStore) LogUser(userId string) {
	s.CurrentUserId = userId
	s.ClearCache()
	s.LoadData()
}

func (s *TodoStore) LoadData() {

	client := &http.Client{}

	url := fmt.Sprintf("http://localhost:8080/api/v1/todo/%s", s.CurrentUserId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("client error")
		fmt.Print(err.Error())
	}

	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Print(err.Error())
	}

	cli.Clr()

	s.Mu.Lock()
	var userData types.TodoStoreData
	json.Unmarshal(body, &userData)

	for i := 0; i < len(userData.Data); i++ {
		for k, v := range userData.Data[i] {
			s.Data[k] = v
		}
	}
	s.Mu.Unlock()
}

func (s *TodoStore) SaveChangesToFile() {

	client := &http.Client{}

	s.Mu.Lock()

	dataTosave := types.TodoStoreData{Data: []map[string]types.Todo{}}
	for k, v := range s.Data {
		dataTosave.Data = append(dataTosave.Data, map[string]types.Todo{k: v})
	}

	s.Mu.Unlock()

	jsonBytes, err := json.MarshalIndent(dataTosave, " ", " ")
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://127.0.0.1:8080/api/v1/todo/new/%s", s.CurrentUserId)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Print(err.Error())
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer res.Body.Close()
}

func (s *TodoStore) DisplayOptions() {
	fmt.Printf("What do you want to do next? \n\n")
	fmt.Printf("1. Add new TODO\n")
	fmt.Printf("2. Mark a TODO as \"done\"\n")
	fmt.Printf("3. Delete TODO\n")
	fmt.Printf("4. Save changes to file\n")
	fmt.Printf("5. Change current user\n")

	choice := cli.GetUserInput("\nSelect one option by number:")
	choiceInt, err := strconv.Atoi(choice)
	for (err != nil) || !(choiceInt < 6 && choiceInt > 0) {
		fmt.Printf("Invalid option!!!\n")
		s.DisplayOptions()
		break
	}

	cli.Clr()

	switch choiceInt {
	case 1:
		newTaskName := cli.GetUserInput("New task description: ")
		s.Add(newTaskName)

	case 2:
		s.ListTodos()
		taskId := cli.GetUserInput("\n- Todo's ID to mark as done:")

		s.MarkAsDone(taskId)

	case 3:
		s.ListTodos()
		todoId := cli.GetUserInput("Todo's ID to be removed:")
		s.Delete(todoId)
		println("\n- Success! Todo removed.\n")

	case 4:
		s.SaveChangesToFile()
		fmt.Printf("\n- Success! Changes saved!\n")

	case 5:
		userId := cli.GetUserInput("Username: ")
		s.LogUser(userId)
	}
	s.ListTodos()
	s.DisplayOptions()
}
