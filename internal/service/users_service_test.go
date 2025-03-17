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

func TestCompareHash(t *testing.T) {
	db := MockDB{}

	user, _ := db.GetUserInfo("")

	// Test right password
	err := compareHash(user, "password")
	if err != nil {
		t.Error(err)
	}

	// Test wrong password
	err = compareHash(user, "worng_password")
	if err == nil {
		t.Error(err)
	}
}

func TestSaltGeneration(t *testing.T) {

	salt1 := generateSalt()
	salt2 := generateSalt()

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
	// Test for expected hash
	salt := []byte("saltsalt")

	hash := generateHash("password", salt)

	expectedHash := "$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNhbHQ$SRzrpBkxb+Cwwr5PQJL2pIGh9G59lfzlgOj3RRV73LKQYf2HycaaTY5yHimy7mnlWCY"

	if hash != expectedHash {
		t.Errorf("The two hashes are not the same.\nFuncHash: %s\nConstHash: %s", hash, expectedHash)
	}

	// Test if running the hash again will generate the same hash
	hash2 := generateHash("password", salt)

	if hash != hash2 {
		t.Errorf("The regenerated hash is not the same as the first.\nHash1: %s\nHash2: %s", hash, hash2)
	}

	// Test that the hashs will not the be the same with different passwords

	hash3 := generateHash("Password", salt)

	if hash3 == hash2 {
		t.Errorf("Hash3 is colliding with hash2 when they are not supposed to be.\nHash3: %s\nHash2: %s", hash3, hash2)
	}

	// Test different salt different hash

	salt2 := []byte("saltSalt")

	hash4 := generateHash("password", salt2)

	if hash4 == hash {
		t.Errorf("Hash4 is not supposed to be equal to hash1 as they have different salts\nHash4: %s\nHash1: %s", hash4, hash)
	}

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
