package main

import (
	"encoding/json"
	"fmt"
	"github.com/NerdBow/GrindersAPI/api/logs"
	"net/http"
	"strconv"
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
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("400 Bad Request"))
			fmt.Println(err)
			return
		}

		if !requestLog.Validate() {
			fmt.Println("Log was not valid")
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("400 Bad Request"))
			return
		}

		fmt.Println("Post")

	case http.MethodGet:
		// This can be factored out into a new function
		requestedId := request.URL.Path[5:]

		logId, err := strconv.Atoi(requestedId)

		if err != nil {
			fmt.Println("Get method attempted with:", requestedId)
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("400 Bad Request"))
			return
		}

		fmt.Printf("Fetch Log %d", logId)

		fmt.Println("Get")

	case http.MethodPut:
		requestedId := request.URL.Path[5:]

		logId, err := strconv.Atoi(requestedId)

		if err != nil {
			fmt.Println("Put method attempted with:", requestedId)
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("400 Bad Request"))
			return
		}
		fmt.Printf("Update Log %d", logId)
		fmt.Println("Put")

	case http.MethodDelete:
		requestedId := request.URL.Path[5:]

		logId, err := strconv.Atoi(requestedId)

		if err != nil {
			fmt.Println("Delete method attempted with:", requestedId)
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("400 Bad Request"))
			return
		}

		fmt.Printf("Delete Log %d", logId)

		fmt.Println("Delete")
	}
}
func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &Handler{})

	mux.Handle("/log/", &LogHandler{})

	http.ListenAndServe(":8080", mux)
}
