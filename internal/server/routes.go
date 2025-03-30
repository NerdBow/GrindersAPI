package server

import (
	"net/http"

	"github.com/NerdBow/GrindersAPI/internal/handler"
	"github.com/NerdBow/GrindersAPI/internal/middleware"
	"github.com/NerdBow/GrindersAPI/internal/service"
)

// addRoutes adds all the handlers for each route to the provided mux.
func addRoutes(mux *http.ServeMux, userService service.UserService, userLogService service.UserLogService) {
	mux.HandleFunc("POST /user/signup", middleware.SetHeader(middleware.SetLog(handler.HandleUserSignUp(userService))))
	mux.HandleFunc("POST /user/signin", middleware.SetHeader(middleware.SetLog(handler.HandleUserSignIn(userService))))

	mux.HandleFunc("POST /user/log", middleware.SetLog(middleware.CheckAuth(handler.HandleUserLogPost(userLogService))))
	mux.HandleFunc("GET /user/log", middleware.SetLog(middleware.CheckAuth(handler.HandleUserLogGet(userLogService))))
	mux.HandleFunc("PUT /user/log", middleware.SetLog(middleware.CheckAuth(handler.HandleUserLogUpdate(userLogService))))
	mux.HandleFunc("DELETE /user/log", middleware.SetLog(middleware.CheckAuth(handler.HandleUserLogDelete(userLogService))))
}
