package service

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NerdBow/GrindersAPI/internal/database"
	"github.com/NerdBow/GrindersAPI/internal/model"
	"github.com/NerdBow/GrindersAPI/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

var (
	BlankFieldsErr        = errors.New("Username and Password must not be blank")
	InvalidPasswordErr    = errors.New("Password must be 8 or more characters")
	InvalidPageErr        = errors.New("Page must be greater than 0")
	InvalidTimeErr        = errors.New("Time must be greater than 0 if filtering by time")
	InvalidLogIdQueryErr  = errors.New("LogId must be greater than 0 for single logs or equal to 0 for multiple logs")
	InvalidOrderErr       = errors.New("Order must be DATE_ASC, DATE_DES, DURATION_ASC, or DURATION_DES")
	UnmergableDurationErr = errors.New("Duration must not be negative for a merge log")
	UnmergableDateErr     = errors.New("Date must not be negative for a merge log")
)

// The service which is used for user/ endpoint.
type UserService struct {
	db database.UserDatabase
}

// Creates a new UserService.
func NewUserService(db database.UserDatabase) UserService {
	return UserService{db: db}
}

// Signs up a user to the database.
// Takes in a username and a password.
//
// Returns a bool if the signup was successful or not and an error if unsuccessful.
func (s *UserService) SignUp(username string, password string) (bool, error) {
	if username == "" || password == "" {
		return false, BlankFieldsErr
	}

	if len(password) < 8 {
		return false, InvalidPasswordErr
	}

	salt, err := util.GenerateSalt()

	if err != nil {
		return false, err
	}

	hash := util.GenerateHash(password, salt)

	err = s.db.SignUp(username, hash)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Signs in a user to the database and creates a JWT.
// Takes in a username and a password.
//
// Returns a string of the JWT and nil if successful.
// If passwords is not correct, returns empty string and nil.
// If errors occur, returns empty string and error.
func (s *UserService) SignIn(username string, password string) (string, error) {
	user, err := s.db.GetUserInfo(username)

	if err != nil {
		return "", err
	}

	ok, err := util.CompareHashToPassword(user.Hash, password)

	if err != nil {
		return "", err
	}

	if !ok {
		return "", nil
	}

	expTime, err := strconv.Atoi(os.Getenv("JWTEXP"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"userId":   strconv.Itoa(user.UserId),
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * time.Duration(expTime)).Unix(),
	}

	jwtString, err := util.CreateToken(claims)

	if err != nil {
		return "", err
	}

	return jwtString, nil
}

// The service which is used for the user/{id}/endpoint.
type UserLogService struct {
	db database.UserLogDatabase
}

// Creates a new UserLogService.
func NewUserLogService(db database.UserLogDatabase) UserLogService {
	return UserLogService{db: db}
}

// Adds a new log to the database.
// Takes in a log struct of the log to be inserted in the database.
//
// Returns the id of the inserted log if successful. -1 and an error if unsuccessful.
func (s *UserLogService) AddUserLog(log model.Log) (int64, error) {
	problems := log.Validate()
	for _, err := range problems {
		return -1, err
	}

	id, err := s.db.PostLog(log)

	if err != nil {
		return -1, err
	}

	return id, nil
}

// Retrives a slice of logs of a user to the database.
// Takes in the userId, logId, page, startTime, endTime, category, and order.
// If logId is greater than 0 then function will return a single log that matches the id.
//
// Returns a slice log struct of the requested log. nil and an error if unsuccessful.
func (s *UserLogService) GetUserLogs(userId int, logId int64, page uint, startTime int64, endTime int64, category string, order string) ([]model.Log, error) {
	var logs []model.Log

	if logId < 0 {
		return logs, InvalidLogIdQueryErr
	}

	if page <= 0 {
		return logs, InvalidPageErr
	}

	if startTime < 0 {
		return logs, InvalidTimeErr
	}

	if endTime < 0 {
		return logs, InvalidTimeErr
	}

	if logId > 0 {
		log, err := s.db.GetLog(userId, logId)

		if err != nil {
			return logs, err
		}

		logs = append(logs, log)

		return logs, nil
	}

	var logOrder database.LogOrder

	switch strings.ToUpper(order) {
	case "DATE_ASC":
		logOrder = database.DATE_ASC
	case "DATE_DES":
		logOrder = database.DATE_DES
	case "DURATION_ASC":
		logOrder = database.DURATION_ASC
	case "DURATION_DES":
		logOrder = database.DURATION_DES
	default:
		return logs, InvalidOrderErr

	}

	filteredLogs, err := s.db.GetLogs(userId, page, startTime, endTime, category, logOrder)

	if err != nil {
		return logs, err
	}

	return filteredLogs, nil
}

// Updates a log to the database.
// Takes in a log struct
//
// Returns true and nil if update was successful. false and an error if not.
func (s *UserLogService) UpdateUserLog(log model.Log) (bool, error) {
	problems := log.Validate()
	for field, err := range problems {
		if field == "id" || field == "userId" {
			return false, err
		}
	}

	if log.Date < 0 {
		return false, UnmergableDateErr
	}

	if log.Duration < 0 {
		return false, UnmergableDurationErr
	}

	result, err := s.db.UpdateLog(log)

	if err != nil || !result {
		return false, err
	}

	return result, nil
}

// Deletes a log of the user in the database.
// Takes in a userId and the logId.
//
// Returns true and nil if delete was successful. false and an error if not.
func (s *UserLogService) DeleteUserLog(userId int, logId int64) (bool, error) {
	if userId <= 0 {
		return false, model.InvalidUserIdErr
	}

	if logId <= 0 {
		return false, model.InvalidIdErr
	}

	result, err := s.db.DeleteLog(userId, logId)

	if err != nil || !result {
		return false, err
	}

	return result, nil
}
