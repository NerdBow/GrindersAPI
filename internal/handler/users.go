package handler

import (
	"encoding/json"
	"fmt"
	"io"
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

// Takes in the body of request and decodes into a log.
//
// Returns the decoded log from the body and an empty log and an error if the decode is not possible.
func decodeLog(b io.ReadCloser) (model.Log, error) {
	decoder := json.NewDecoder(b)
	decoder.DisallowUnknownFields()

	var requestLog model.Log
	err := decoder.Decode(&requestLog)

	if err != nil {
		return model.Log{}, http.ErrBodyNotAllowed
	}
	return requestLog, nil
}

// Gets the userId from the jwt
// TODO: THIS IS A TEMP FUNCTION. ONCE I ADD IN THE JWT STUFF I WILL NEED TO CHANGE THIS
func getIdFromToken(jwt string) int {
	return 0
}

// Handler for POST request of user/log endpoint.
// Allows user to add a log.
//
// Returns http handler func.
func HandleUserLogPost(s service.UserLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Assume that the ID is valid by middlewear
		requestedLog, err := decodeLog(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := s.AddUserLog(getIdFromToken(""), requestedLog)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dataBytes, err := json.Marshal(struct {
			id int `json:"id"`
		}{id})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(dataBytes)
		return
	}
}

// Handler for GET request of user/log endpoint.
// Allows user to get information of one of their logs.
//
// Returns http handler func.
func HandleUserLogGet(s service.UserLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		queryValues := r.URL.Query()

		logId := queryValues.Get("id")

		if logId == "" { // For the single log request
			logId, err := strconv.Atoi(logId)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			log, err := s.GetUserLog(getIdFromToken(""), logId)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data, err := json.Marshal(log)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(data)
			return
		}

		page := queryValues.Get("p")
		category := queryValues.Get("c")
		startTime := queryValues.Get("s")
		timeFrame := queryValues.Get("f")

		s.GetUserLogs()

		// This can be factored out into a new function
		// requestedId := r.URL.Path[5:]
		// requestedUser := r.PathValue("username")
		//
		// logId, err := strconv.Atoi(requestedId)
		//
		// if err != nil {
		// 	fmt.Println("Get method attempted with:", requestedId)
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	w.Write([]byte("400 Bad Request"))
		// 	return
		// }
		//
		// userId, err := strconv.Atoi(requestedUser)
		//
		// if err != nil {
		// 	fmt.Println("Get method attempted with:", userId)
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	w.Write([]byte("400 Bad Request"))
		// 	return
		// }
		//
		// queriedLog, err := s.GetUserLog(userId, logId)
		//
		// if err != nil {
		// 	fmt.Println(err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	w.Write([]byte("500 Internal Server Error"))
		// 	return
		// }
		//
		// dataBytes, err := json.Marshal(queriedLog)
		//
		// if err != nil {
		// 	fmt.Println(err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	w.Write([]byte("500 Internal Server Error"))
		// 	return
		// }
		//
		// w.Header().Set("Content-Type", "application/json")
		//
		// w.Write(dataBytes)
		//
		// fmt.Printf("Fetch Log %d\n", logId)
		//
		// fmt.Println("Get")
	}
}

// Handler for PUT request of user/{}/log endpoint.
// Allows user to update information of one of their logs.
//
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
//
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
