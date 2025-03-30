package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/NerdBow/GrindersAPI/internal/middleware"
	"github.com/NerdBow/GrindersAPI/internal/model"
	"github.com/NerdBow/GrindersAPI/internal/service"
)

var (
	NoUsernameFieldErr = errors.New("Username field must be provided in the request json or not empty")
	NoPasswordFieldErr = errors.New("Password field must be provided in the request json or not empty")
	NoBodyErr          = errors.New("Request must have a json body")
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
			if err.Error() == "EOF" {
				middleware.HandleError(w, NoBodyErr, http.StatusBadRequest, NoBodyErr.Error())
				return
			}
			middleware.HandleError(w, err, http.StatusBadRequest, err.Error())
			return
		}

		if userInfo.Username == "" {
			middleware.HandleError(w, NoUsernameFieldErr, http.StatusBadRequest, NoUsernameFieldErr.Error())
			return
		}

		if userInfo.Password == "" {
			middleware.HandleError(w, NoPasswordFieldErr, http.StatusBadRequest, NoPasswordFieldErr.Error())
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
			if err.Error() == "EOF" {
				middleware.HandleError(w, NoBodyErr, http.StatusBadRequest, NoBodyErr.Error())
				return
			}
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

// Handler for POST request of user/log endpoint.
// Allows user to add a log.
// Returns http handler func.
func HandleUserLogPost(s service.UserLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := middleware.GetUserFromCtx(r.Context())
		if err != nil {
			middleware.HandleError(w, err, http.StatusUnauthorized, "Unauthorized")
			return
		}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		var requestedLog model.Log
		err = decoder.Decode(&requestedLog)

		if err != nil {
			if err.Error() == "EOF" {
				middleware.HandleError(w, NoBodyErr, http.StatusBadRequest, NoBodyErr.Error())
				return
			}
			middleware.HandleError(w, err, http.StatusBadRequest, err.Error())
			return
		}

		// The id does not matter here as it is not used for the adding of a log
		requestedLog.Id = 1
		requestedLog.UserId = user.UserId

		id, err := s.AddUserLog(requestedLog)

		if err != nil {
			middleware.HandleError(w, err, http.StatusBadRequest, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		responseJson := struct {
			Id int64 `json:"id"`
		}{id}

		messageBytes, err := json.Marshal(responseJson)

		if err != nil {
			log.Printf("Could not write error")
		}

		w.Write(messageBytes)
		return
	}
}

// Handler for GET request of user/{}/log endpoint.
// Allows user to get information of one of their logs.
// Returns http handler func.
func HandleUserLogGet(s service.UserLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

// Handler for PUT request of user/{}/log endpoint.
// Allows user to update information of one of their logs.
// Returns http handler func.
func HandleUserLogUpdate(s service.UserLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

// Handler for DELETE request of user/{}/log endpoint.
// Allows user to delete of one of their logs.
// Returns http handler func.
func HandleUserLogDelete(s service.UserLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
