package server

import (
	"github.com/NerdBow/GrindersAPI/internal/handler"
	"github.com/NerdBow/GrindersAPI/internal/service"
	"net/http"
)

// addRoutes adds all the handlers for each route to the provided mux.
func addRoutes(mux *http.ServeMux, userService service.UserService, userLogService service.UserLogService) {
	mux.HandleFunc("POST /users/logs", handler.HandleUserLogPost(userLogService))
	mux.HandleFunc("GET /users/logs", handler.HandleUserLogGet(userLogService))
	mux.HandleFunc("PUT /users/logs", handler.HandleUserLogUpdate(userLogService))
	mux.HandleFunc("DELETE /users/logs", handler.HandleUserLogDelete(userLogService))
}
