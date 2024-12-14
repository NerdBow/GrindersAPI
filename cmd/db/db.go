package db

import (
	"database/sql"
	"errors"
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

	temp, err := result.LastInsertId()

	id := int(temp)

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

	row, err := db.db.Query(fmt.Sprintf("SELECT id, date, duration, name, category, userId FROM 'logs' WHERE id = %d;", id))

	defer row.Close()

	if err != nil {
		return log, err
	}

	if !row.Next() {
		return log, errors.New("There was no row with that primary key.")
	}

	err = row.Scan(&log.Id, &log.Date, &log.Duration, &log.Name, &log.Category, &log.UserId)

	if err != nil {
		return log, err
	}

	return log, nil
}

func (db Sqlite3DB) UpdateLog(newLogAddr *logs.Log) (bool, error) {
	newLog := *newLogAddr
	log, err := db.GetLog(newLog.Id)

	if err != nil {
		return false, err
	}

	newLog.Merge(log)

	statement, err := db.db.Prepare(
		fmt.Sprintf("UPDATE 'logs' SET date = %d, duration = %d, name = \"%s\", category = \"%s\" WHERE id = %d;",
			newLog.Date, newLog.Duration, newLog.Name, newLog.Category, newLog.Id))

	if err != nil {
		return false, err
	}

	result, err := statement.Exec()

	if err != nil {
		fmt.Println(result)
		return false, err
	}

	return true, nil
}

func (db Sqlite3DB) DeleteLog(id int) error {
	statement, err := db.db.Prepare(fmt.Sprintf("DELETE FROM 'logs' WHERE id = %d", id))
	if err != nil {
		return err
	}

	result, err := statement.Exec()

	if err != nil {
		fmt.Println(result)
		return err
	}

	return nil
}

func (db Sqlite3DB) Close() error {
	return db.db.Close()
}

type Database interface {
	PostLog(logs.Log) (int, error)
	GetRecentLog() (int, error)
	GetLog(int) (logs.Log, error)
	UpdateLog(*logs.Log) (bool, error)
	DeleteLog(int) error
	Close() error
}
