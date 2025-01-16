package apiv1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	clistore "todoapp/store/cli-store"
)

func saveUserJsonFile(w http.ResponseWriter, userId string, data clistore.TodoStoreData) {
	bytesToSave, err := json.MarshalIndent(data, " ", " ")
	if err != nil {
		InternalServerErrorHandler(w)
	}

	os.WriteFile(fmt.Sprintf("data/%s.json", userId), bytesToSave, os.ModePerm)
}
