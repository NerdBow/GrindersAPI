package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func run() error {
	db, err := sql.Open("sqlite3", "/data/logs.db")
	if err != nil {
		return err
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY, date INTEGER, duration INTEGER, name TEXT, category TEXT, userId INTEGER)")

	if err != nil {
		return err
	}

	statement.Exec()

	defer db.Close()
	return nil
}
