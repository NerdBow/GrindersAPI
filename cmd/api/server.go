package main

import (
	"encoding/json"
	"fmt"
	"github.com/NerdBow/GrindersAPI/api/logs"
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
		decoder := json.NewDecoder(request.Body)
		var requestLog logs.Log
		err := decoder.Decode(&requestLog)

		if err != nil {
			fmt.Println(err)
			return
		}

		if !requestLog.Validate() {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("400 Bad Request"))
		}

		fmt.Println("Post")

	case http.MethodGet:
		logId := request.URL.Path[5:]

		fmt.Println(request.URL.Path, logId)

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
