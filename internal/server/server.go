package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NerdBow/GrindersAPI/internal/database"
	"github.com/NerdBow/GrindersAPI/internal/service"
	"github.com/joho/godotenv"
)

// Starts the server for the API
func Run() {
	if os.Getenv("DATABASE") == "" {
		err := godotenv.Load(".env")

		if err != nil {
			log.Println(err)
		}
	}

	mux := http.NewServeMux()

	sqlService := os.Getenv("DATABASE")

	databaseFile := os.Getenv("DATABASEFILE")

	var userService service.UserService
	var userLogService service.UserLogService

	switch sqlService {
	case "sqlite3":
		db, err := database.NewSqlite3DB(databaseFile)

		if err != nil {
			fmt.Println(err)
			return
		}

		err = db.CreateTables()

		if err != nil {
			fmt.Println(err)
			return
		}

		userService = service.NewUserService(db)
		userLogService = service.NewUserLogService(db)
		log.Println("Sqlite3 is being used.")
	default:
		log.Fatal("No database specification found in the .env")
	}

	addRoutes(mux, userService, userLogService)

	server := http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Printf("Server started on port %s", server.Addr)

	log.Fatal(server.ListenAndServe())
}
