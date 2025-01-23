package types

import (
	"time"
)

type TodoStoreData struct {
	Data []map[string]Todo `json:"data"`
}

type Todo struct {
	Task      string    `json:"task" validate:"required"`
	IsDone    bool      `json:"isDone" validate:"boolean"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}

//not really a useful interface but hey
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
