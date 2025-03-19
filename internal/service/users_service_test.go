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
