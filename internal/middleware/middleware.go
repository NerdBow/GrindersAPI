package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/NerdBow/GrindersAPI/internal/model"
	"github.com/NerdBow/GrindersAPI/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

type ContextKey uint8

var (
	UserKey        ContextKey = 1
	NoUserInCtxErr            = errors.New("There was no user in the context of the request")
)

// Checks for JWT token in request Authorization header and parses it to a User struct which goes into the request's context.
func CheckAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authString := r.Header.Get("Authorization")
		parsedAuth := strings.Split(authString, " ")
		if len(parsedAuth) != 2 {
			HandleError(w, errors.New("JWT is not formated correctly"), http.StatusBadRequest, "JWT is not formatted correctly")
			return
		}

		bearer, token := parsedAuth[0], parsedAuth[1]

		if bearer != "Bearer" {
			HandleError(w, errors.New("No Bearer in auth header"), http.StatusBadRequest, "JWT is not formatted correctly")
			return
		}

		claims, err := util.GetClaimsFromToken(token)

		if errors.Is(err, util.CanNotAccessJWTSecretErr) {
			HandleError(w, err, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			HandleError(w, nil, http.StatusUnauthorized, "Token Expired")
			return
		}

		if err != nil {
			HandleError(w, err, http.StatusBadRequest, "JWT is not formatted correctly")
			return
		}

		var user model.User
		var ok bool

		userId, ok := claims["userId"].(string)

		user.UserId, err = strconv.Atoi(userId)

		if !ok || err != nil {
			HandleError(w, errors.New("JWT does not have a userId"), http.StatusBadRequest, "JWT is not formatted correctly")
			return
		}

		user.Username, ok = claims["username"].(string)

		if !ok {
			HandleError(w, errors.New("JWT does not have a username"), http.StatusBadRequest, "JWT is not formatted correctly")
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, user)
		log.Printf("User: %d, %s", user.UserId, user.Username)
		handler(w, r.WithContext(ctx))
	}
}

func GetUserFromCtx(ctx context.Context) (model.User, error) {
	value := ctx.Value(UserKey)
	user, ok := value.(model.User)
	if !ok {
		return model.User{}, NoUserInCtxErr
	}
	return user, nil
}

// Logs the IP of where the request is coming from and which endpoint they are requesting.
func SetLog(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("IP: %-22s Endpoint: %s\n", r.RemoteAddr, r.URL)
		handler(w, r)
	}
}

// Sets the header of the writer to application/json.
func SetHeader(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}
}

// Write out the message into the WriteResponse as a JSON
func WriteResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	responseJson := struct {
		Message string `json:"message"`
	}{message}

	messageBytes, err := json.Marshal(responseJson)

	if err != nil {
		log.Printf("Could not write error")
	}

	w.Write(messageBytes)
}

// Logs and writes out message for any error case in handlers.
func HandleError(w http.ResponseWriter, err error, statusCode int, message string) {
	if err != nil {
		log.Printf("Error: %-60s Message: %s\n", err, message)
	}
	WriteResponse(w, statusCode, message)
}
