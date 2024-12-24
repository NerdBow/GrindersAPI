package db

import (
	"errors"
	"fmt"

	"crypto/rand"
	"crypto/sha256"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Username string
	Salt     string
	Hash     string
}

type UserService interface {
	SignUp(string, string) (bool, error)
	SignIn(string, string) (bool, error)
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

func (db *Sqlite3DB) SignUp(username string, password string) (bool, error) {
	const saltLength int = 10
	row := db.QueryRow("SELECT username FROM 'users' WHERE username = ?", username)

	var queryUsername string
	err := row.Scan(&queryUsername)

	if err == nil {
		fmt.Println("There already exist that username")
		return false, err
	}

	salt, err := GenerateSalt(saltLength)

	if err != nil {
		fmt.Println("Salt failed to generate")
		return false, err
	}

	hash, err := HashPassword(append([]byte(password), salt...))

	if err != nil {
		return false, err
	}

	statement, err := db.Prepare("INSERT INTO 'users' (username, salt, hash) VALUES(?, ?, ?);")

	if err != nil {
		return false, err
	}

	_, err = statement.Exec(username, string(salt), string(hash))

	if err != nil {
		return false, err
	}

	return true, nil
}

func (db *Sqlite3DB) SignIn(username string, password string) (string, error) {
	row := db.QueryRow("SELECT username, salt, hash FROM 'users' WHERE username = ?", username)

	var queriedUser User
	err := row.Scan(&queriedUser.Username, &queriedUser.Salt, &queriedUser.Hash)

	if err != nil {
		fmt.Println("This user does not exist")
		return "", err
	}

	hash, err := HashPassword([]byte(password + queriedUser.Salt))

	if err != nil {
		fmt.Println("Hash went wrong")
		return "", err
	}

	if string(hash) != queriedUser.Hash {
		return "", errors.New("Password was not correct")
	}

	// This should return a JWT or some sort of session token
	return "Success", nil

}
