package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NerdBow/GrindersAPI/internal/database"
	"github.com/NerdBow/GrindersAPI/internal/service"
)

// Starts the server for the API
func Run() {
	mux := http.NewServeMux()

	db, err := database.NewSqlite3DB()

	if err != nil {
		fmt.Println(err)
		return
	}

	err = db.CreateTables()

	if err != nil {
		fmt.Println(err)
		return
	}

	userService := service.NewUserService(db)
	userLogService := service.NewUserLogService(db)

	addRoutes(mux, userService, userLogService)

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Println("Server started.")

	log.Fatal(server.ListenAndServe())
}
