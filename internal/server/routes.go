package server

import (
	"github.com/NerdBow/GrindersAPI/internal/handler"
	"github.com/NerdBow/GrindersAPI/internal/service"
	"net/http"
)

// addRoutes adds all the handlers for each route to the provided mux.
func addRoutes(mux *http.ServeMux, userService service.UserService, userLogService service.UserLogService) {

	mux.HandleFunc("POST /user/{id}/log", handler.HandleUserLogPost(userLogService))
	mux.HandleFunc("GET /user/{id}/log", handler.HandleUserLogGet(userLogService))
	mux.HandleFunc("PUT user/{id}/log", handler.HandleUserLogUpdate(userLogService))
	mux.HandleFunc("DELETE user/{id}/log", handler.HandleUserLogDelete(userLogService))
	//
	// mux.HandleFunc("user/signin", nil)
	// mux.HandleFunc("user/signup", nil)
	//
	// mux.HandleFunc("user/", nil)
	//
	// mux.HandleFunc("logs/", nil)
}
