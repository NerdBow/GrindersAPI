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

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username STRING, salt STRING, hash STRING)")

	if err != nil {
		return nil, err
	}

	statement.Exec()

	realDb := Sqlite3DB{db}

	return realDb, nil
}

type Sqlite3DB struct{ *sql.DB }

func (db Sqlite3DB) PostLog(log logs.Log) (int, error) {
	statement, err := db.Prepare("INSERT INTO 'logs' (date, duration, name, category, userId) VALUES(?, ?, ?, ?, ?);")

	if err != nil {
		return -1, err
	}

	result, err := statement.Exec(log.Date, log.Duration, log.Name, log.Category, log.UserId)

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

func (db Sqlite3DB) GetLog(id int) (logs.Log, error) {
	log := logs.Log{}

	row, err := db.Query("SELECT id, date, duration, name, category, userId FROM 'logs' WHERE id = ?;", id)

	if err != nil {
		return log, err
	}

	defer row.Close()

	if !row.Next() {
		return log, errors.New("There was no row with that primary key.")
	}

	err = row.Scan(&log.Id, &log.Date, &log.Duration, &log.Name, &log.Category, &log.UserId)

	if err != nil {
		return log, err
	}

	return log, nil
}

func (db Sqlite3DB) GetLogs(page int, category string) (*[]logs.Log, error) {
	logsList := make([]logs.Log, 0, 10)

	var query string

	if category == "" {
		query = "SELECT id, date, duration, name, category, userId FROM 'logs' WHERE category IS NOT ? ORDER BY date DESC LIMIT ?, ?;"
	} else {
		query = "SELECT id, date, duration, name, category, userId FROM 'logs' WHERE category = ? ORDER BY date DESC LIMIT ?, ?;"
	}

	rows, err := db.Query(query, category, page*10, (page+1)*10)

	if err != nil {
		return &logsList, err
	}

	defer rows.Close()

	for rows.Next() {
		rowLog := logs.Log{}
		err = rows.Scan(&rowLog.Id, &rowLog.Date, &rowLog.Duration, &rowLog.Name, &rowLog.Category, &rowLog.UserId)
		if err != nil {
			return &logsList, err
		}
		logsList = append(logsList, rowLog)
	}

	return &logsList, nil
}

func (db Sqlite3DB) UpdateLog(newLogAddr *logs.Log) (bool, error) {
	newLog := *newLogAddr
	log, err := db.GetLog(newLog.Id)

	if err != nil {
		return false, err
	}

	newLog.Merge(log)

	statement, err := db.Prepare("UPDATE 'logs' SET date = ?, duration = ?, name = ?, category = ? WHERE id = ?;")

	if err != nil {
		return false, err
	}

	result, err := statement.Exec(newLog.Date, newLog.Duration, newLog.Name, newLog.Category, newLog.Id)

	if err != nil {
		fmt.Println(result)
		return false, err
	}

	return true, nil
}

func (db Sqlite3DB) DeleteLog(id int) (bool, error) {
	statement, err := db.Prepare("DELETE FROM 'logs' WHERE id = ?")

	if err != nil {
		return false, err
	}

	result, err := statement.Exec(id)

	if err != nil {
		fmt.Println(result)
		return false, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		fmt.Println(result)
		return false, err
	}

	if rowsUpdated == 0 {
		return false, errors.New("There was no log deleted")
	}

	return true, nil
}

func (db Sqlite3DB) Close() error {
	return db.Close()
}

type Database interface {
	PostLog(logs.Log) (int, error)
	GetLog(int) (logs.Log, error)
	GetLogs(int, string) (*[]logs.Log, error)
	UpdateLog(*logs.Log) (bool, error)
	DeleteLog(int) (bool, error)
	Close() error
}
