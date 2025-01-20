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

	"github.com/google/uuid"
)

type TodoStore struct {
	CurrentUserId string
	Data          map[string]Todo
	Mu            sync.Mutex
}

type TodoStoreData struct {
	Data []map[string]Todo `json:"data"`
}

type Todo struct {
	Task      string    `json:"task" validate:"required"`
	IsDone    bool      `json:"isDone" validate:"boolean"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}
type TodoStorer interface {
	Add()
	Delete()
	MarkAsDone()
	ListTodos()
	ClearCache()
	ChangeCurrentUser()
	LoadData()
	SaveChangesToFile()
	DisplayOptions()
}

func (s *TodoStore) Add() {

	cli.Clr()
	taskUuid := uuid.New()

	newTaskName := cli.GetUserInput("New task description: ")

	s.Mu.Lock()
	s.Data[taskUuid.String()] = Todo{
		Task:      newTaskName,
		IsDone:    false,
		CreatedAt: time.Now(),
	}
	s.Mu.Unlock()

	cli.Clr()
	println("\n- Success! Task added.\n")
	s.ListTodos()
	s.DisplayOptions()
}

func (s *TodoStore) Delete() {
	s.Mu.Lock()

	taskId := cli.GetUserInput("Todo's ID to be removed:")
	delete(s.Data, taskId)
	println("\n- Success! Todo removed.\n")

	s.Mu.Unlock()
	s.ListTodos()
	s.DisplayOptions()
}

func (s *TodoStore) MarkAsDone() {
	cli.Clr()
	s.ListTodos()
	taskId := cli.GetUserInput("\n- Todo's ID to mark as done:")

	s.Mu.Lock()

	if entry, ok := s.Data[taskId]; ok {
		entry.IsDone = true
		s.Data[taskId] = entry
	}

	s.Mu.Unlock()
	cli.Clr()
	fmt.Printf("\n- Success! Todo marked as done.\n")
	s.ListTodos()
	s.DisplayOptions()
}

func (s *TodoStore) ListTodos() {
	fmt.Printf("\n######### TODOS ##########\n")
	i := 1
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
}

func (s *TodoStore) ClearCache() {
	s.Mu.Lock()

	s.Data = map[string]Todo{}

	s.Mu.Unlock()
}

func (s *TodoStore) LogUser() {
	s.CurrentUserId = cli.GetUserInput("Username: ")
	s.ClearCache()
	s.LoadData()
	s.ListTodos()
	s.DisplayOptions()
}

func (s *TodoStore) LoadData() {

	url := fmt.Sprintf("http://127.0.0.1:8080/api/v1/todo/%s", s.CurrentUserId)
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

	cli.Clr()

	s.Mu.Lock()
	var userData TodoStoreData
	json.Unmarshal(body, &userData)

	for i := 0; i < len(userData.Data); i++ {
		for k, v := range userData.Data[i] {
			s.Data[k] = v
		}
	}
	s.Mu.Unlock()
}

func (s *TodoStore) SaveChangesToFile() {
	s.Mu.Lock()

	dataTosave := TodoStoreData{Data: []map[string]Todo{}}
	for k, v := range s.Data {
		dataTosave.Data = append(dataTosave.Data, map[string]Todo{k: v})
	}

	jsonBytes, err := json.MarshalIndent(dataTosave, " "," ")
	if err != nil {
		log.Fatal(err)
	}
	//
	//
	//
	url := fmt.Sprintf("http://127.0.0.1:8080/api/v1/todo/new/%s", s.CurrentUserId)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Print(err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer res.Body.Close()

	// body, readErr := io.ReadAll(res.Body)
	// if readErr != nil {
	// 	fmt.Print(err.Error())
	// }
	//
	//
	//
	// os.WriteFile(fmt.Sprintf("data/%s.json", s.CurrentUserId), jsonBytes, os.ModePerm)
	s.Mu.Unlock()

	cli.Clr()
	fmt.Printf("\n- Success! Changes saved!\n")
	s.ListTodos()
	s.DisplayOptions()
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
	switch choiceInt {
	case 1:
		s.Add()
	case 2:
		s.MarkAsDone()
	case 3:
		s.Delete()
	case 4:
		s.SaveChangesToFile()
	case 5:
		s.LogUser()
	}
}
