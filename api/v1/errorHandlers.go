package apiv1

import "net/http"

func FileNotFoundErrorHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	http.Error(w, "404 Invalid user or user data not found", 404)
}

func BadRequestErrorHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	http.Error(w, "400 JSON body not formatted properly", 400)
}

func InternalServerErrorHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}
