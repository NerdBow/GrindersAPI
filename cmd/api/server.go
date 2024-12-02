package main

import (
	"fmt"
	"net/http"
)

type Handler struct{}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Hello"))
}

type LogHandler struct{}

func (handler *LogHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		fmt.Println("Post")

	case http.MethodGet:
		fmt.Println("Get")

	case http.MethodPut:
		fmt.Println("Put")

	case http.MethodDelete:
		fmt.Println("Delete")
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &Handler{})

	mux.Handle("/log/", &LogHandler{})

	http.ListenAndServe(":8080", mux)
}
