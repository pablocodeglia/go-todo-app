package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TodoStore struct {
	CurrentUserId string
	Data          map[string]Todo
	Mu            sync.Mutex
}

type TodoStoreJson struct {
	Data []map[string]Todo `json:"data"`
}

type UserStore struct {
	Data map[string]interface{}
	Mu   sync.Mutex
}

type User struct {
	FirstName string
	LastName  string
}

type Todo struct {
	Task      string    `json:"task"`
	IsDone    bool      `json:"isDone"`
	CreatedAt time.Time `json:"createdAt"`
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
}

func (s *TodoStore) ListTodos() {
	fmt.Println("##########################")
	fmt.Println("######### TODOS ##########")
	fmt.Println("##########################")
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
}

func (s *TodoStore) LoadData() {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	file, err := os.Open(fmt.Sprintf("data/%s.json", s.CurrentUserId))
	if err != nil {
		fmt.Println("No previous saved data found.")
	}

	byteValue, _ := io.ReadAll(file)

	var userData TodoStoreJson
	json.Unmarshal([]byte(byteValue), &userData)

	for i := 0; i < len(userData.Data); i++ {
		for k, v := range userData.Data[i] {
			s.Data[k] = v
		}
	}
}

func (s *TodoStore) SaveChangesToFile() {
	//TODO add logic
	dataTosave := TodoStoreJson{Data: []map[string]Todo{}}
	for k, v := range s.Data {
		dataTosave.Data = append(dataTosave.Data, map[string]Todo{k: v})
	}

	jsonBytes, err := json.Marshal(dataTosave)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(fmt.Sprintf("data/%s.json", s.CurrentUserId), jsonBytes, os.ModePerm)
}

func (s *TodoStore) Add() {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	taskUuid := uuid.New()

	newTaskName := getUserInput("New task description: ")

	s.Data[taskUuid.String()] = Todo{
		Task:      newTaskName,
		IsDone:    false,
		CreatedAt: time.Now(),
	}

	s.ListTodos()
	displayOptions(s)
}
func (s *TodoStore) Delete() {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	taskId := getUserInput("Todo's ID to be removed:")
	delete(s.Data, taskId)
	println("Success! Todo removed.")
}
func (s *TodoStore) MarkAsDone() {
	taskId := getUserInput("Todo's ID to mark as done:")

	s.Mu.Lock()
	defer s.Mu.Unlock()

	if entry, ok := s.Data[taskId]; ok {
		entry.IsDone = true
		s.Data[taskId] = entry
	}
	fmt.Printf("Success! Todo marked as done.\n")
	s.ListTodos()
}
func (s *TodoStore) ClearCache() {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Data = map[string]Todo{}
}
func (s *TodoStore) ChangeCurrentUser() {
	s.CurrentUserId = getUserInput("Username: ")
	s.ClearCache()
	s.LoadData()
	s.ListTodos()
	displayOptions(s)
}

func NewStore() *TodoStore {
	return &TodoStore{
		Data: make(map[string]Todo),
	}
}

func main() {
	store := NewStore()
	store.ChangeCurrentUser()

	displayOptions(store)

}

func getUserInput(q string) string {
	fmt.Println(q)
	reader := bufio.NewReader(os.Stdin)
	userId, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(userId)
}

func displayOptions(store *TodoStore) {
	fmt.Printf("What do you want to do next? \n\n")
	fmt.Printf("1. Add new TODO\n")
	fmt.Printf("2. Mark a TODO as \"done\"\n")
	fmt.Printf("3. Delete TODO\n")
	fmt.Printf("4. Save changes to file\n")
	fmt.Printf("5. Change current user\n")
	choice := getUserInput("\nSelect one option number:")
	choiceInt, err := strconv.Atoi(choice)
	for (err != nil) || !(choiceInt < 6 && choiceInt > 0) {
		fmt.Printf("Invalid option!!!\n")
		displayOptions(store)
		break
	}
	switch choiceInt {
	case 1:
		store.Add()
	case 2:
		store.MarkAsDone()
	case 3:
		store.Delete()
	case 4:
		store.SaveChangesToFile()
	case 5:
		store.ChangeCurrentUser()
	}
}
