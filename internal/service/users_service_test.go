package service

import (
	"errors"
	"fmt"
	"testing"

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

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Println("env could not load")
	}
	m.Run()
}

func TestSignUp(t *testing.T) {
	s := NewUserService(&MockDB{make(map[string]struct{})})

	// Test valid username and password
	err := s.SignUp("TestUser", "Password")

	if err != nil {
		t.Error(err)
	}

	// Test duplicate username
	err = s.SignUp("TestUser", "Password")

	if err == nil {
		t.Error(err)
	}

	var blankFieldsErr *BlankFieldsError
	var invalidPasswordErr *InvalidPasswordError

	// Test blank username and password
	err = s.SignUp("", "")

	if !errors.As(err, &blankFieldsErr) {
		t.Error(err)
	}

	// Test blank password
	err = s.SignUp("a", "")

	if !errors.As(err, &blankFieldsErr) {
		t.Error(err)
	}

	// Test blank username
	err = s.SignUp("", "a")

	if !errors.As(err, &blankFieldsErr) {
		t.Error(err)
	}

	// Test invalid password
	err = s.SignUp("a", "asdfjkl")

	if !errors.As(err, &invalidPasswordErr) {
		t.Error(err)
	}

	// Test one letter username
	err = s.SignUp("a", "asdfjkl;")

	if errors.As(err, &invalidPasswordErr) {
		t.Error(err)
	}
}
