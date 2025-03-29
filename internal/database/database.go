package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NerdBow/GrindersAPI/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

type LogOrder uint8

const (
	PAGE_ROW_COUNT        int      = 20
	ASC_DATE_ASC_DURATION LogOrder = 1
	ASC_DATE_DES_DURATION LogOrder = 2
	DES_DATE_ASC_DURATION LogOrder = 3
	DES_DATE_DES_DURATION LogOrder = 4
)

type UserLogDatabase interface {
	// Adds the given log into the database.
	//
	// Returns an int of the log's id in the database and an sql error if one occurs.
	PostLog(model.Log) (int64, error)

	// Retrives a specificed log from the database of a user by userId and logId.
	//
	// Returns log struct of the specificed log if successful.
	// Else, it returns an empty log and the error.
	GetLog(int, int64) (model.Log, error)

	// Retrives a slice of 10 logs specified by the parameters of:
	//
	// userId, page, startTime(epoch of the lower bound of logs), timeLength(DAY, WEEK, MONTH since the start time), category
	//
	// Returns a pointer of the slice of logs if successful.
	// Else, it returns a pointer with the an empty slice of logs and the error.
	GetLogs(int, uint, int64, int64, string, LogOrder) ([]model.Log, error)

	// Update the specified log with new information in the given log.
	//
	// Returns true of the operation was successful.
	// Else, it returns false and an error if it was unsuccessful.
	UpdateLog(model.Log) (bool, error)

	// Deletes the specified log from the userId and logId.
	//
	// Returns true of the operation was successful.
	// Else, it returns false and an error if it was unsuccessful.
	DeleteLog(int, int64) (bool, error)
}

type UserDatabase interface {
	// Adds user to the database provided their username and hashed password.
	//
	// Returns an error or nil depending if an error happened.
	SignUp(string, string) error

	// Get the information of the user from the database with the given username
	//
	// Returns the user struct of the specified user and nil if successful else empty user and error
	GetUserInfo(string) (model.User, error)
}

// Shuts off the connection to the database.
//
// Returns an error if any occur.
func Close(db *sql.DB) error {
	return db.Close()
}

type Sqlite3DB struct{ *sql.DB }

func NewSqlite3DB(file string) (*Sqlite3DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL;")

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

	statement, err = db.Prepare("CREATE INDEX IF NOT EXISTS logs_index_0 ON logs (userId, date);")

	if err != nil {
		return err
	}

	statement.Exec()

	statement, err = db.Prepare("CREATE INDEX IF NOT EXISTS logs_index_1 ON logs (userId, category, date);")

	if err != nil {
		return err
	}

	statement.Exec()

	statement, err = db.Prepare("CREATE INDEX IF NOT EXISTS logs_index_2 ON logs (userId, name, date);")

	if err != nil {
		return err
	}

	statement.Exec()
	return nil
}

func (db Sqlite3DB) PostLog(log model.Log) (int64, error) {
	statement, err := db.Prepare("INSERT INTO 'logs' (date, duration, name, category, goal, userId) VALUES(?, ?, ?, ?, ?, ?);")

	if err != nil {
		return -1, err
	}

	result, err := statement.Exec(log.Date, log.Duration, log.Name, log.Category, log.Goal, log.UserId)

	if err != nil {
		return -2, err
	}

	logId, err := result.LastInsertId()

	if err != nil {
		return -3, err
	}

	return logId, nil
}

func (db Sqlite3DB) GetLog(userId int, id int64) (model.Log, error) {
	log := model.Log{}

	row, err := db.Query("SELECT id, date, duration, name, category, goal, userId FROM 'logs' WHERE id = ? AND userId = ?;", id, userId)

	if err != nil {
		return log, err
	}

	defer row.Close()

	if !row.Next() {
		return log, errors.New("Unable to find log with the id for the user")
	}

	err = row.Scan(&log.Id, &log.Date, &log.Duration, &log.Name, &log.Category, &log.Goal, &log.UserId)

	if err != nil {
		return log, err
	}

	return log, nil
}

func (db Sqlite3DB) GetLogs(userId int, page uint, startTime int64, endTime int64, category string, ordering LogOrder) ([]model.Log, error) {
	logsList := make([]model.Log, 0, PAGE_ROW_COUNT)

	query := "SELECT id, date, duration, name, category, goal, userId FROM 'logs' WHERE userId = ? AND (? = 0 OR date >= ?) AND (? = 0 OR date <= ?) AND (? = '' OR category = ?)"

	switch ordering {
	case ASC_DATE_ASC_DURATION:
		query += " ORDER BY date ASC, duration ASC"

	case ASC_DATE_DES_DURATION:
		query += " ORDER BY date ASC, duration DESC"

	case DES_DATE_ASC_DURATION:
		query += " ORDER BY date DESC, duration ASC"

	case DES_DATE_DES_DURATION:
		query += " ORDER BY date DESC, duration DESC"
	}

	query += " LIMIT ?, ?;"

	rows, err := db.Query(query, userId, startTime, startTime, endTime, endTime, category, category, ((page - 1) * 10), page*10)

	if err != nil {
		return logsList, err
	}

	defer rows.Close()

	for rows.Next() {
		rowLog := model.Log{}
		err = rows.Scan(&rowLog.Id, &rowLog.Date, &rowLog.Duration, &rowLog.Name, &rowLog.Category, &rowLog.Goal, &rowLog.UserId)
		if err != nil {
			return logsList, err
		}
		logsList = append(logsList, rowLog)
	}

	return logsList, nil
}

func (db Sqlite3DB) UpdateLog(newLog model.Log) (bool, error) {
	log, err := db.GetLog(newLog.UserId, newLog.Id)

	if err != nil {
		return false, err
	}

	newLog.Merge(log)

	statement, err := db.Prepare("UPDATE 'logs' SET date = ?, duration = ?, name = ?, category = ?, goal = ? WHERE id = ? AND userId = ?;")

	if err != nil {
		return false, err
	}

	_, err = statement.Exec(newLog.Date, newLog.Duration, newLog.Name, newLog.Category, newLog.Goal, newLog.Id, newLog.UserId)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (db Sqlite3DB) DeleteLog(userId int, id int64) (bool, error) {
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
		return false, errors.New("There is no log that belongs to the user with that id")
	}

	return true, nil
}

func (db Sqlite3DB) SignUp(username string, hash string) error {
	statement, err := db.Prepare("INSERT INTO 'users' (username, hash) VALUES(?, ?)")

	if err != nil {
		return err
	}

	_, err = statement.Exec(username, hash)

	if err != nil {
		return err
	}

	return nil
}

func (db Sqlite3DB) GetUserInfo(username string) (model.User, error) {
	var user model.User

	row, err := db.Query("SELECT id, username, hash FROM 'users' WHERE username = ?;", username)

	if err != nil {
		return user, err
	}

	defer row.Close()

	if !row.Next() {
		return user, errors.New("No User Exist")
	}

	err = row.Scan(&user.UserId, &user.Username, &user.Hash)

	if err != nil {
		return user, err
	}

	return user, nil
}
