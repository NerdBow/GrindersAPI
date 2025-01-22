package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/NerdBow/GrindersAPI/internal/model"
	"github.com/NerdBow/GrindersAPI/internal/service"
)

func handleUserSignIn(s service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func handleUserSignUp(s service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func HandleUserLogPost(s service.UserLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		var requestLog model.Log
		err := decoder.Decode(&requestLog)

		fmt.Println(r.PathValue("id"))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad request"))
			fmt.Println(err)
			return
		}

		if len(requestLog.Validate()) != 0 {
			fmt.Println(requestLog.Validate())
			fmt.Println("Log was not valid")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad request"))
			return
		}

		logId, err := s.AddUserLog(requestLog)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			return
		}

		data := struct {
			Id int `json:"id"`
		}{logId}

		dataBytes, err := json.Marshal(data)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(dataBytes)

		fmt.Println("Post")
	}
}

func HandleUserLogGet(s service.UserLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// This can be factored out into a new function
		requestedId := r.URL.Path[5:]
		requestedUser := r.PathValue("username")

		logId, err := strconv.Atoi(requestedId)

		if err != nil {
			fmt.Println("Get method attempted with:", requestedId)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request"))
			return
		}

		userId, err := strconv.Atoi(requestedUser)

		if err != nil {
			fmt.Println("Get method attempted with:", userId)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request"))
			return
		}

		queriedLog, err := s.GetUserLog(userId, logId)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			return
		}

		dataBytes, err := json.Marshal(queriedLog)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(dataBytes)

		fmt.Printf("Fetch Log %d\n", logId)

		fmt.Println("Get")
	}
}

func HandleUserLogUpdate(s service.UserLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		var requestLog model.Log
		err := decoder.Decode(&requestLog)

		if err != nil || requestLog.Id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad request"))
			fmt.Println(err)
			return
		}

		result, err := s.UpdateUserLog(requestLog)

		if err != nil || !result {
			fmt.Println(result)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			fmt.Println(err)
			return
		}

		data := struct {
			Id int `json:"id"`
		}{requestLog.Id}

		dataBytes, err := json.Marshal(data)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(dataBytes)

		fmt.Printf("Update Log %d", requestLog.Id)

		fmt.Println("Put")
	}
}

func HandleUserLogDelete(s service.UserLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestedId := r.URL.Path[5:]

		logId, err := strconv.Atoi(requestedId)

		if err != nil {
			fmt.Println("Delete method attempted with:", requestedId)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad request"))
			return
		}

		result, err := s.DeleteUserLog(0, logId)

		if err != nil || !result {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			fmt.Println(err)
			return
		}

		data := struct {
			Id int `json:"id"`
		}{logId}

		dataBytes, err := json.Marshal(data)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(dataBytes)

		fmt.Printf("Delete Log %d", logId)

		fmt.Println("Delete")
	}
}
