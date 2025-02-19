package service

import (
	"github.com/NerdBow/GrindersAPI/internal/database"
	"github.com/NerdBow/GrindersAPI/internal/model"
)

// The service which is used for user/ endpoint.
type UserService struct {
	db database.Database
}

// Creates a new UserService.
func NewUserService(db database.Database) UserService {
	return UserService{db: db}
}

// Signs up a user to the database.
// Takes in a username and a password.
//
// Returns a bool if the signup was successful or not and an error if unsuccessful.
func (s *UserService) SignUp(username string, password string) error {
	return nil
}

// Signs in a user to the database and creates a session or passes back an API key. //TODO: decide if I want an API KEY or a Session based system
// Takes in a username and a password.
//
// Returns a string of the API Key or session token and if any error an error is returned as well.
func (s *UserService) SignIn(username string, password string) error {
	return nil
}

// The service which is used for the user/{id}/endpoint.
type UserLogService struct {
	db database.Database
}

// Creates a new UserLogService.
func NewUserLogService(db database.Database) UserLogService {
	return UserLogService{db: db}
}

// Adds a new log to the database.
// Takes in a log struct of the log to be inserted in the database.
//
// Returns the id of the inserted log if successful. -1 and an error if unsuccessful.
func (s *UserLogService) AddUserLog(log model.Log) (int, error) {
	return s.db.PostLog(log)
}

// Retrives a log to the database.
// Takes in the userId and the logId
//
// Returns a log struct of the requested log. Empty struct and an error if unsuccessful.
func (s *UserLogService) GetUserLog(userId int, logId int) (model.Log, error) {
	return s.db.GetLog(userId, logId)
}

// Retrives a slice of logs of a user to the database.
// Takes in the userId, age, startTime, timeLength, and category
//
// Returns a slice log struct of the requested log. nil and an error if unsuccessful.
func (s *UserLogService) GetUserLogs(userId int, page int, startTime int64, timeLength string, category string) (*[]model.Log, error) {
	return s.db.GetLogs(userId, page, startTime, timeLength, category)
}

// Updates a log to the database.
// Takes in a log struct
//
// Returns true and nil if update was successful. false and an error if not.
func (s *UserLogService) UpdateUserLog(log model.Log) (bool, error) {
	return s.db.UpdateLog(log)
}

// Deletes a log of the user in the database.
// Takes in a userId and the logId.
//
// Returns true and nil if delete was successful. false and an error if not.
func (s *UserLogService) DeleteUserLog(userId int, logId int) (bool, error) {
	return s.db.DeleteLog(userId, logId)
}
