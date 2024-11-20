package main

import (
	"net/http"
)

type Handler struct{}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Hello"))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &Handler{})

	http.ListenAndServe(":8080", mux)
}
