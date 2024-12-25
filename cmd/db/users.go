package db

import (
	"crypto/rand"
	"crypto/sha256"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Username string
	Salt     string
	Hash     string
}

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, 10)
	_, err := rand.Read(salt)

	if err != nil {
		return salt, err
	}

	return salt, nil
}

func HashPassword(password []byte) ([]byte, error) {
	hashFunc := sha256.New()
	_, err := hashFunc.Write(password)

	if err != nil {
		return make([]byte, 0), nil
	}

	return hashFunc.Sum(nil), nil
}
