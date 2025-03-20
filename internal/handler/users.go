package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/NerdBow/GrindersAPI/internal/middleware"
	"github.com/NerdBow/GrindersAPI/internal/model"
	"github.com/NerdBow/GrindersAPI/internal/service"
)

func HandleUserSignIn(s service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		userInfo := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		err := decoder.Decode(&userInfo)

		if err != nil {
			middleware.HandleError(w, err, http.StatusBadRequest, err.Error())
			return
		}

		token, err := s.SignIn(userInfo.Username, userInfo.Password)

		if err != nil {
			middleware.HandleError(w, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if token == "" {
			middleware.HandleError(w, errors.New("Username or Password is incorrect"), http.StatusForbidden, "Username or Password is incorrect")
			return
		}

		data := struct {
			Token string `json:"token"`
		}{token}

		dataBytes, err := json.Marshal(data)

		if err != nil {
			middleware.HandleError(w, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		w.Write(dataBytes)
		return
	}
}

func HandleUserSignUp(s service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		userInfo := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		err := decoder.Decode(&userInfo)

		if err != nil {
			middleware.HandleError(w, err, http.StatusBadRequest, err.Error())
			return
		}

		ok, err := s.SignUp(userInfo.Username, userInfo.Password)

		if !ok && errors.Is(err, service.InvalidPasswordErr) {
			middleware.HandleError(w, err, http.StatusBadRequest, err.Error())
			return
		}

		if !ok && errors.Is(err, service.BlankFieldsErr) {
			middleware.HandleError(w, err, http.StatusBadRequest, err.Error())
			return
		}

		if !ok || err != nil {
			middleware.HandleError(w, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		middleware.WriteResponse(w, http.StatusOK, "Successfuly created account")
		log.Printf("New account created | Username: %s", userInfo.Username)
	}
}

// Handler for POST request of user/{}/log endpoint.
// Allows user to add a log.
// Returns http handler func.
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

// Handler for GET request of user/{}/log endpoint.
// Allows user to get information of one of their logs.
// Returns http handler func.
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

// Handler for PUT request of user/{}/log endpoint.
// Allows user to update information of one of their logs.
// Returns http handler func.
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

// Handler for DELETE request of user/{}/log endpoint.
// Allows user to delete of one of their logs.
// Returns http handler func.
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
