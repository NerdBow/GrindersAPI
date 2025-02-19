package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/NerdBow/GrindersAPI/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

// Defines the methods a database needs in order to be usable for the API.
type Database interface {
	// Adds the given log into the database.
	//
	// Returns an int of the log's id in the database and an sql error if one occurs.
	PostLog(model.Log) (int, error)

	// Retrives a specificed log from the database of a user by userId and logId.
	//
	// Returns log struct of the specificed log if successful.
	// Else, it returns an empty log and the error.
	GetLog(int, int) (model.Log, error)

	// Retrives a slice of 10 logs specified by the parameters of:
	//
	// userId, page, startTime(epoch of the lower bound of logs), timeLength(DAY, WEEK, MONTH since the start time), category
	//
	// Returns a pointer of the slice of logs if successful.
	// Else, it returns a pointer with the an empty slice of logs and the error.
	GetLogs(int, int, int64, string, string) (*[]model.Log, error)

	// Update the specified log with new information in the given log.
	//
	// Returns true of the operation was successful.
	// Else, it returns false and an error if it was unsuccessful.
	UpdateLog(model.Log) (bool, error)

	// Deletes the specified log from the userId and logId.
	//
	// Returns true of the operation was successful.
	// Else, it returns false and an error if it was unsuccessful.
	DeleteLog(int, int) (bool, error)

	SignUp(string, string) error

	// Shuts off the connection to the database.
	//
	// Returns an error if any occur.
	Close() error
}

func Start() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "data/logs.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

type Sqlite3DB struct{ *sql.DB }

func NewSqlite3DB() (*Sqlite3DB, error) {
	db, err := sql.Open("sqlite3", "data/logs.db")
	if err != nil {
		return nil, err
	}
	return &Sqlite3DB{db}, nil
}

func (db Sqlite3DB) CreateTables() error {

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT NOT NULL UNIQUE, hash TEXT NOT NULL UNIQUE)")

	if err != nil {
		return err
	}

	statement.Exec()

	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS logs " +
		"(id INTEGER PRIMARY KEY, date INTEGER NOT NULL, duration INTEGER NOT NULL, name TEXT NOT NULL, category TEXT NOT NULL, goal TEXT NOT NULL, userId INTEGER NOT NULL, " +
		"FOREIGN KEY (userId) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE NO ACTION);")

	if err != nil {
		return err
	}

	statement.Exec()

	return nil
}

func (db Sqlite3DB) PostLog(log model.Log) (int, error) {
	statement, err := db.Prepare("INSERT INTO 'logs' (date, duration, name, category, userId) VALUES(?, ?, ?, ?, ?);")

	if err != nil {
		return -1, err
	}

	result, err := statement.Exec(log.Date, log.Duration, log.Name, log.Category, log.UserId)

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

func (db Sqlite3DB) GetLog(userId int, id int) (model.Log, error) {
	log := model.Log{}

	row, err := db.Query("SELECT id, date, duration, name, category, userId FROM 'logs' WHERE id = ? AND userId = ?;", id, userId)

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

func (db Sqlite3DB) GetLogs(userId int, page int, startTime int64, timeLength string, category string) (*[]model.Log, error) {
	logsList := make([]model.Log, 0, 10)

	query := "SELECT id, date, duration, name, category, userId FROM 'logs' WHERE userId = ? AND (? = 0 OR date >= ?) AND (? = 0 OR date <= ?) AND (? = '' OR category = ?) ORDER BY date DESC LIMIT ?, ?;"

	var endTime int64
	switch timeLength {
	case "DAY":
		endTime = startTime + int64(time.Hour*24)
	case "WEEK":
		endTime = startTime + int64(time.Hour*24*7)
	case "MONTH":
		endTime = startTime + int64(time.Hour*24*7*30)
	}

	rows, err := db.Query(query, userId, startTime, startTime, endTime, endTime, category, category, ((page - 1) * 10), page*10)

	if err != nil {
		return &logsList, err
	}

	defer rows.Close()

	for rows.Next() {
		rowLog := model.Log{}
		err = rows.Scan(&rowLog.Id, &rowLog.Date, &rowLog.Duration, &rowLog.Name, &rowLog.Category, &rowLog.UserId)
		if err != nil {
			return &logsList, err
		}
		logsList = append(logsList, rowLog)
	}

	return &logsList, nil
}

func (db Sqlite3DB) UpdateLog(newLog model.Log) (bool, error) {
	log, err := db.GetLog(newLog.UserId, newLog.Id)

	if err != nil {
		return false, err
	}

	newLog.Merge(log)

	statement, err := db.Prepare("UPDATE 'logs' SET date = ?, duration = ?, name = ?, category = ? WHERE id = ? AND userId = ?;")

	if err != nil {
		return false, err
	}

	result, err := statement.Exec(newLog.Date, newLog.Duration, newLog.Name, newLog.Category, newLog.Id, newLog.UserId)

	if err != nil {
		fmt.Println(result)
		return false, err
	}

	return true, nil
}

func (db Sqlite3DB) DeleteLog(userId int, id int) (bool, error) {
	statement, err := db.Prepare("DELETE FROM 'logs' WHERE id = ? AND userId = ?")

	if err != nil {
		return false, err
	}

	result, err := statement.Exec(id, userId)

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

func (db Sqlite3DB) SignUp(string, string) error {
	return nil
}

func (db Sqlite3DB) Close() error {
	return db.Close()
}
