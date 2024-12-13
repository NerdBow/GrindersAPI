package db

import (
	"database/sql"
	"fmt"
	"github.com/NerdBow/GrindersAPI/api/logs"

	_ "github.com/mattn/go-sqlite3"
)

func Start() (Database, error) {
	db, err := sql.Open("sqlite3", "data/logs.db")
	if err != nil {
		return nil, err
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY, date INTEGER, duration INTEGER, name TEXT, category TEXT, userId INTEGER)")

	if err != nil {
		return nil, err
	}

	statement.Exec()

	realDb := Sqlite3DB{db: db}

	return realDb, nil
}

type Sqlite3DB struct {
	db *sql.DB
}

func (db Sqlite3DB) PostLog(log logs.Log) (int, error) {
	statement, err := db.db.Prepare(
		fmt.Sprintf("INSERT INTO 'logs' (date, duration, name, category, userId) VALUES(%d, %d, \"%s\", \"%s\", %d);", log.Date, log.Duration, log.Name, log.Category, log.UserId))

	if err != nil {
		return -1, err
	}

	result, err := statement.Exec()

	fmt.Println(result)

	if err != nil {
		return -2, err
	}

	id, err := db.GetRecentLog()

	if err != nil {
		return -3, err
	}

	return id, nil
}

func (db Sqlite3DB) GetRecentLog() (int, error) {

	var rowId int

	row, err := db.db.Query("SELECT id FROM 'logs' ORDER BY id DESC LIMIT 1;")

	defer row.Close()

	if err != nil {
		return -1, err
	}

	if !row.Next() {
		return -1, nil
	}

	err = row.Scan(&rowId)

	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (db Sqlite3DB) GetLog(id int) (logs.Log, error) {
	log := logs.Log{}
	return log, nil
}

func (db Sqlite3DB) UpdateLog(id int) (logs.Log, error) {
	log := logs.Log{}
	return log, nil
}

func (db Sqlite3DB) DeleteLog(id int) error {
	return nil
}

func (db Sqlite3DB) Close() error {
	return db.db.Close()
}

type Database interface {
	PostLog(logs.Log) (int, error)
	GetRecentLog() (int, error)
	GetLog(int) (logs.Log, error)
	UpdateLog(int) (logs.Log, error)
	DeleteLog(int) error
	Close() error
}
