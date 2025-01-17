package apiv1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	clistore "todoapp/store/cli-store"
)

func saveUserJsonFile(w http.ResponseWriter, userId string, data clistore.TodoStoreData) {
	bytesToSave, err := json.MarshalIndent(data, " ", " ")
	if err != nil {
		InternalServerErrorHandler(w)
	}

	os.WriteFile(fmt.Sprintf("data/%s.json", userId), bytesToSave, os.ModePerm)
}

func findIndexByTodoIdFunc(userData clistore.TodoStoreData, todoId string) int {
	i := slices.IndexFunc(userData.Data, func(data map[string]clistore.Todo) bool {
		for k := range data {
			if k == todoId {
				return true
			}
		}
		return false
	})
	return i
}
