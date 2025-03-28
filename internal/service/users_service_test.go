package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/NerdBow/GrindersAPI/internal/database"
	"github.com/NerdBow/GrindersAPI/internal/model"
	"github.com/joho/godotenv"
)

type MockDB struct {
	usernames map[string]struct{}
}

func (db *MockDB) SignUp(username string, password string) error {
	if _, ok := db.usernames[username]; ok {
		return errors.New("There exist the username " + username)
	}
	db.usernames[username] = struct{}{}
	return nil
}

func (db *MockDB) GetUserInfo(username string) (model.User, error) {
	return model.User{Username: "NerdBow", Hash: "$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNhbHQ$SRzrpBkxb+Cwwr5PQJL2pIGh9G59lfzlgOj3RRV73LKQYf2HycaaTY5yHimy7mnlWCY"}, nil
}

func (db *MockDB) PostLog(log model.Log) (int64, error) {
	return 20, nil
}

func (db *MockDB) GetLog(userId int, id int64) (model.Log, error) {
	var log model.Log
	return log, nil
}

func (db *MockDB) GetLogs(userId int, page uint, startTime int64, endTime int64, category string, order database.LogOrder) ([]model.Log, error) {
	logs := make([]model.Log, 10, 10)
	return logs, nil
}

func (db *MockDB) UpdateLog(log model.Log) (bool, error) {
	return true, nil
}

func (db *MockDB) DeleteLog(userId int, id int64) (bool, error) {
	return true, nil
}

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Println("env could not load")
	}
	m.Run()
}

func TestUpdateUserLog(t *testing.T) {
	s := NewUserLogService(&MockDB{})

	// Test valid input
	log := model.Log{Id: 10, Date: 10000, Duration: 3600, Name: "Calculas", Category: "Study", Goal: "Finish reviewing for test", UserId: 1}

	result, err := s.UpdateUserLog(log)

	if err != nil || !result {
		t.Error(err)
	}

	// Test blank date
	log = model.Log{Id: 10, Date: 0, Duration: 3600, Name: "Calculas", Category: "Study", Goal: "Finish reviewing for test", UserId: 1}

	result, err = s.UpdateUserLog(log)

	if err != nil || !result {
		t.Error(err)
	}

	// Test blank duration
	log = model.Log{Id: 10, Date: 10, Duration: 0, Name: "Calculas", Category: "Study", Goal: "Finish reviewing for test", UserId: 1}

	result, err = s.UpdateUserLog(log)

	if err != nil || !result {
		t.Error(err)
	}

	// Test blank name
	log = model.Log{Id: 10, Date: 10, Duration: 100, Name: "", Category: "Study", Goal: "Finish reviewing for test", UserId: 1}

	result, err = s.UpdateUserLog(log)

	if err != nil || !result {
		t.Error(err)
	}

	// Test blank category
	log = model.Log{Id: 10, Date: 10, Duration: 100, Name: "Calculas", Category: "", Goal: "Finish reviewing for test", UserId: 1}

	result, err = s.UpdateUserLog(log)

	if err != nil || !result {
		t.Error(err)
	}

	// Test blank goal
	log = model.Log{Id: 10, Date: 10, Duration: 100, Name: "Calculas", Category: "Study", Goal: "", UserId: 1}

	result, err = s.UpdateUserLog(log)

	if err != nil || !result {
		t.Error(err)
	}

	// Test invalid id
	log = model.Log{Id: 0, Date: 10000, Duration: 3600, Name: "Calculas", Category: "Study", Goal: "Finish reviewing for test", UserId: 1}

	result, err = s.UpdateUserLog(log)

	if !errors.Is(err, model.InvalidIdErr) || result {
		t.Error(err)
	}

	// Test invalid date
	log = model.Log{Id: 1, Date: -1, Duration: 3600, Name: "Calculas", Category: "Study", Goal: "Finish reviewing for test", UserId: 1}

	result, err = s.UpdateUserLog(log)

	if !errors.Is(err, model.InvalidDateErr) || result {
		t.Error(err)
	}

	// Test invalid duration
	log = model.Log{Id: 1, Date: 100, Duration: -1, Name: "Calculas", Category: "Study", Goal: "Finish reviewing for test", UserId: 1}

	result, err = s.UpdateUserLog(log)

	if !errors.Is(err, model.InvalidDurationErr) || result {
		t.Error(err)
	}

	// Test invalid user id
	log = model.Log{Id: 1, Date: 100, Duration: 3600, Name: "Calculas", Category: "Study", Goal: "Finish reviewing for test", UserId: -1}

	result, err = s.UpdateUserLog(log)

	if !errors.Is(err, model.InvalidUserIdErr) || result {
		t.Error(err)
	}
}

func TestGetUserLogs(t *testing.T) {
	s := NewUserLogService(&MockDB{})

	// Test single log retrival
	logs, err := s.GetUserLogs(1, 1, 1, 0, 0, "", database.ASC_DATE_ASC_DURATION)

	if err != nil {
		t.Error(err)
	}

	if len(logs) != 1 {
		t.Error("There was more than one log returned for the GetLog by id")
	}

	// Test invalid logId
	logs, err = s.GetUserLogs(1, -10, 1, 0, 0, "", database.ASC_DATE_ASC_DURATION)

	if !errors.Is(err, InvalidLogIdQueryErr) || len(logs) != 0 {
		t.Error("There was no error returns when negative logId was inputted")
	}

	// Test multiple log retrival
	logs, err = s.GetUserLogs(1, 0, 1, 0, 0, "", database.ASC_DATE_ASC_DURATION)

	if err != nil {
		t.Error(err)
	}

	if len(logs) != 10 {
		t.Error("There was less than 10 logs returned")
	}

	// Test invalid page
	logs, err = s.GetUserLogs(1, 0, 0, 0, 0, "", database.ASC_DATE_ASC_DURATION)

	if !errors.Is(err, InvalidPageErr) || len(logs) != 0 {
		t.Error(err)
	}

	// Test invalid startTime
	logs, err = s.GetUserLogs(1, 0, 1, -1, 0, "", database.ASC_DATE_ASC_DURATION)

	if !errors.Is(err, InvalidTimeErr) || len(logs) != 0 {
		t.Error(err)
	}

	// Test invalid endTime
	logs, err = s.GetUserLogs(1, 0, 1, 10, -1, "", database.ASC_DATE_ASC_DURATION)

	if !errors.Is(err, InvalidTimeErr) || len(logs) != 0 {
		t.Error(err)
	}
}

func TestAddUserLog(t *testing.T) {
	s := NewUserLogService(&MockDB{})

	// Test correct log
	log := model.Log{Id: 10, Date: 10000, Duration: 3600, Name: "Calculas", Category: "Study", Goal: "Finish reviewing for test", UserId: 1}

	id, err := s.AddUserLog(log)

	if err != nil || id != 20 {
		t.Error(err)
	}

	// Test invalid id
	log = model.Log{Id: -1, Date: 10000, Duration: 3600, Name: "Calculas", Category: "Study", Goal: "Finish reviewing for test", UserId: 1}

	id, err = s.AddUserLog(log)

	if err == nil || id != -1 {
		t.Error(err)
	}

	// Test invalid date
	log = model.Log{Id: 10, Date: -10000, Duration: 3600, Name: "Calculas", Category: "Study", Goal: "Finish reviewing for test", UserId: 1}

	id, err = s.AddUserLog(log)

	if err == nil || id != -1 {
		t.Error(err)
	}

	// Test invalid duration
	log = model.Log{Id: 10, Date: -10000, Duration: 0, Name: "Calculas", Category: "Study", Goal: "Finish reviewing for test", UserId: 1}

	id, err = s.AddUserLog(log)

	if err == nil || id != -1 {
		t.Error(err)
	}

	// Test blank name
	log = model.Log{Id: 10, Date: -10000, Duration: 3600, Name: "", Category: "Study", Goal: "Finish reviewing for test", UserId: 1}

	id, err = s.AddUserLog(log)

	if err == nil || id != -1 {
		t.Error(err)
	}

	// Test blank category
	log = model.Log{Id: 10, Date: -10000, Duration: 3600, Name: "Calculas", Category: "", Goal: "Finish reviewing for test", UserId: 1}

	id, err = s.AddUserLog(log)

	if err == nil || id != -1 {
		t.Error(err)
	}

	// Test blank goal
	log = model.Log{Id: 10, Date: -10000, Duration: 3600, Name: "Calculas", Category: "Study", Goal: "", UserId: 1}

	id, err = s.AddUserLog(log)

	if err == nil || id != -1 {
		t.Error(err)
	}

	// Test blank goal
	log = model.Log{Id: 10, Date: -10000, Duration: 3600, Name: "Calculas", Category: "Study", Goal: "Finish reviewing for test", UserId: 0}

	id, err = s.AddUserLog(log)

	if err == nil || id != -1 {
		t.Error(err)
	}
}

func TestSignin(t *testing.T) {
	s := NewUserService(&MockDB{})

	// Test correct signin info
	token, err := s.SignIn("NerdBow", "password")

	if err != nil {
		t.Error(err)
	}

	// Test wrong signin info
	token, err = s.SignIn("NerdBow", "pass")

	if token != "" || err != nil {
		t.Error(err)
	}
}

func TestSignUp(t *testing.T) {
	s := NewUserService(&MockDB{make(map[string]struct{})})

	// Test valid username and password
	ok, err := s.SignUp("TestUser", "Password")

	if !ok || err != nil {
		t.Error(err)
	}

	// Test duplicate username
	ok, err = s.SignUp("TestUser", "Password")

	if ok || err == nil {
		t.Error(err)
	}

	// Test blank username and password
	ok, err = s.SignUp("", "")

	if ok || !errors.Is(err, BlankFieldsErr) {
		t.Error(err)
	}

	// Test blank password
	ok, err = s.SignUp("a", "")

	if ok || !errors.Is(err, BlankFieldsErr) {
		t.Error(err)
	}

	// Test blank username
	ok, err = s.SignUp("", "a")

	if ok || !errors.Is(err, BlankFieldsErr) {
		t.Error(err)
	}

	// Test invalid password
	ok, err = s.SignUp("a", "asdfjkl")

	if ok || !errors.Is(err, InvalidPasswordErr) {
		t.Error(err)
	}

	// Test one letter username
	ok, err = s.SignUp("a", "asdfjkl;")

	if !ok || errors.Is(err, InvalidPasswordErr) {
		t.Error(err)
	}
}
