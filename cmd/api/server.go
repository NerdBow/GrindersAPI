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

func (handler *LogHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(request.Body)
		decoder.DisallowUnknownFields()
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

		writer.Header().Set("Content-Type", "application/json")

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

		queriedLog, err := handler.db.GetLog(logId)

		if err != nil {
			fmt.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 Internal Server Error"))
			return
		}

		dataBytes, err := json.Marshal(queriedLog)

		if err != nil {
			fmt.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 Internal Server Error"))
			return
		}

		writer.Header().Set("Content-Type", "application/json")

		writer.Write(dataBytes)

		fmt.Printf("Fetch Log %d\n", logId)

		fmt.Println("Get")

	case http.MethodPut:
		decoder := json.NewDecoder(request.Body)
		decoder.DisallowUnknownFields()
		var requestLog logs.Log
		err := decoder.Decode(&requestLog)

		if err != nil || requestLog.Id == 0 {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("400 Bad Request"))
			fmt.Println(err)
			return
		}

		result, err := handler.db.UpdateLog(&requestLog)

		if err != nil || !result {
			fmt.Println(result)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 Internal Server Error"))
			fmt.Println(err)
			return
		}

		data := struct {
			Id int `json:"id"`
		}{requestLog.Id}

		dataBytes, err := json.Marshal(data)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 Internal Server Error"))
			fmt.Println(err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")

		writer.Write(dataBytes)

		fmt.Printf("Update Log %d", requestLog.Id)

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

		result, err := handler.db.DeleteLog(logId)

		if err != nil || !result {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 Internal Server Error"))
			fmt.Println(err)
			return
		}

		data := struct {
			Id int `json:"id"`
		}{logId}

		dataBytes, err := json.Marshal(data)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 Internal Server Error"))
			fmt.Println(err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")

		writer.Write(dataBytes)

		fmt.Printf("Delete Log %d", logId)

		fmt.Println("Delete")
	}
}

type LogsHandler struct{ db db.Database }

func (handler *LogsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	const categoryParameter string = "category"
	const pageParameter string = "page"

	switch request.Method {
	case http.MethodGet:
		queryValues := request.URL.Query()
		category := queryValues.Get(categoryParameter)
		page, err := strconv.Atoi(queryValues.Get(pageParameter))

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("400 Bad Request"))
			fmt.Println(err)
			return
		}

		logs, err := handler.db.GetLogs(page, category)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 Internal Server Error"))
			fmt.Println(err)
			return
		}

		dataBytes, err := json.Marshal(*logs)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 Internal Server Error"))
			fmt.Println(err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")

		writer.Write(dataBytes)

	}
}

type SignInHandler struct{ db db.Database }

func (handler *SignInHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(request.Body)
		decoder.DisallowUnknownFields()
		var requestedUser struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := decoder.Decode(&requestedUser)

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("400 Bad Request"))
			fmt.Println(err)
			return
		}

		result, err := handler.db.SignIn(requestedUser.Username, requestedUser.Password)

		// Add different error when the user is not in the db
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("400 Bad Request"))
			fmt.Println(err)
			return
		}

		token := struct {
			Token string `json:"token"`
		}{result}

		dataBytes, err := json.Marshal(token)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 Internal Server Error"))
			fmt.Println(err)
			return
		}

		writer.WriteHeader(http.StatusOK)

		writer.Header().Set("Content-Type", "application/json")

		writer.Write(dataBytes)

		return

	}
	return
}

type SignUpHandler struct{ db db.Database }

func (handler *SignUpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(request.Body)
		decoder.DisallowUnknownFields()
		var requestedUser struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := decoder.Decode(&requestedUser)

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("400 Bad Request"))
			fmt.Println(err)
			return
		}

		result, err := handler.db.SignUp(requestedUser.Username, requestedUser.Password)

		if result && err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("400 Bad Request. Username is taken"))
			fmt.Println(err)
			return
		} else if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("500 Internal Server Error"))
			fmt.Println(err)
			return
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("200 Account Created"))

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

	mux.Handle("/logs/", &LogsHandler{db: sql})

	mux.Handle("/user/signup/", &SignUpHandler{db: sql})

	mux.Handle("/user/signin/", &SignInHandler{db: sql})

	http.ListenAndServe(":8080", mux)
}
