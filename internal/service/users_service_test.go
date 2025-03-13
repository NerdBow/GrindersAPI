package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

type MockDB struct {
	tempHash string
}

func (db MockDB) SignUp(username string, password string) error {
	if password == db.tempHash {
		return errors.New("The hash is the same!!")
	}
	return nil
}

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Println("env could not load")
	}
	m.Run()
}

func TestSaltGeneration(t *testing.T) {

	s := NewUserService(MockDB{})
	salt1 := s.generateSalt()
	salt2 := s.generateSalt()

	same := true

	for i := range salt1 {
		if salt1[i] != salt2[i] {
			same = false
			break
		}
	}
	if same {
		t.Errorf("The two generated salts are the same. HOW IS THIS POSSIBLE???\nsalt1: %x\nsalt2: %x", salt1, salt2)
	}

}

func TestHashGeneration(t *testing.T) {
	s := NewUserService(MockDB{})
	salt := []byte("saltsalt")

	hash := s.generateHash("password", salt)

	expectedHash := "$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNhbHQ$SRzrpBkxb+Cwwr5PQJL2pIGh9G59lfzlgOj3RRV73LKQYf2HycaaTY5yHimy7mnlWCY"

	if hash != expectedHash {
		t.Errorf("The two hashes are not the same.\nFuncHash: %s\nConstHash: %s", hash, expectedHash)
	}

}

func TestSignUp(t *testing.T) {
	s := NewUserService(MockDB{})

	err := s.SignUp("TestUser", "Password")

	if err != nil {
		t.Error(err)
	}

	err = s.SignUp("TestUser", "Password")

	if err != nil {
		t.Error(err)
	}
}
