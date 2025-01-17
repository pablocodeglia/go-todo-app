package apiv1

import (
	"fmt"
	"net/http"
)

func FileNotFoundErrorHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	http.Error(w, "404 Invalid user or user data not found", 404)
}

func BadRequestErrorHandler(w http.ResponseWriter, errMsg string) {
	w.WriteHeader(http.StatusBadRequest)
	http.Error(w, fmt.Sprintf("400 JSON body not formatted properly\n%s", errMsg), 400)

}

func InternalServerErrorHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}
