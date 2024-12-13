package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/NerdBow/GrindersAPI/api/logs"
	"github.com/NerdBow/GrindersAPI/cmd/db"
)

type Handler struct{}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Hello"))
}

type LogHandler struct{ db db.Database }

func (handler LogHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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

		logId, err := handler.db.PostLog(requestLog)

		if err != nil {
			fmt.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 Internal Server Error"))
			return
		}

		writer.Header().Set("Content-Type", "application/json")

		data := struct {
			Id int `json:"id"`
		}{logId}

		dataBytes, err := json.Marshal(data)

		if err != nil {
			fmt.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 Internal Server Error"))
			return
		}

		writer.Write(dataBytes)

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

	sql, err := db.Start()

	defer sql.Close()

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", &Handler{})

	mux.Handle("/log/", &LogHandler{db: sql})

	http.ListenAndServe(":8080", mux)
}
