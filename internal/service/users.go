package service

import (
	"github.com/NerdBow/GrindersAPI/internal/database"
	"github.com/NerdBow/GrindersAPI/internal/model"
)

type UserService struct {
	db database.Database
}

func NewUserService(db database.Database) UserService {
	return UserService{db: db}
}

func (s *UserService) SignUp(username string, password string) (bool, error) {
	return s.db.SignUp(username, password)
}

func (s *UserService) SignIn(username string, password string) (string, error) {
	return s.db.SignIn(username, password)
}

type UserLogService struct {
	db database.Database
}

func NewUserLogService(db database.Database) UserLogService {
	return UserLogService{db: db}
}

func (s *UserLogService) AddUserLog(log model.Log) (int, error) {
	return s.db.PostLog(log)
}

func (s *UserLogService) GetUserLog(userId int, logId int) (model.Log, error) {
	return s.db.GetLog(userId, logId)
}

func (s *UserLogService) GetUserLogs(userId int, page int, startTime int64, timeLength string, category string) (*[]model.Log, error) {
	return s.db.GetLogs(userId, page, startTime, timeLength, category)
}

func (s *UserLogService) UpdateUserLog(log model.Log) (bool, error) {
	return s.db.UpdateLog(log)
}

func (s *UserLogService) DeleteUserLog(userId int, logId int) (bool, error) {
	return s.db.DeleteLog(userId, logId)
}
