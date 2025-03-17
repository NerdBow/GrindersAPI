package util

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

func GenerateSalt() []byte {

	length, err := strconv.Atoi(os.Getenv("SALTLENGTH"))

	if err != nil {
		log.Fatal(err)
	}

	salt := make([]byte, length)
	rand.Read(salt)
	return salt
}

// Takes in the password string and returns the argonid2 encoded version of it.
func GenerateHash(password string, saltBytes []byte) string {

	hashTime, err := strconv.Atoi(os.Getenv("HASHTIME"))

	if err != nil {
		log.Fatal(err)
	}

	hashMemory, err := strconv.Atoi(os.Getenv("HASHMEMORY"))

	if err != nil {
		log.Fatal(err)
	}

	hashThreads, err := strconv.Atoi(os.Getenv("HASHTHREADS"))

	if err != nil {
		log.Fatal(err)
	}

	hashLength, err := strconv.Atoi(os.Getenv("HASHLENGHT"))

	if err != nil {
		log.Fatal(err)
	}

	hashBytes := argon2.IDKey([]byte(password), saltBytes, uint32(hashTime), uint32(hashMemory*1024), uint8(hashThreads), uint32(hashLength))

	salt := base64.RawStdEncoding.EncodeToString(saltBytes)
	hash := base64.RawStdEncoding.EncodeToString(hashBytes)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, hashMemory*1024, hashTime, hashThreads, salt, hash)
}

func ParseHash(hash string) ([]byte, string, error) {
	parsedHash := strings.Split(hash, "$")
	salt, err := base64.RawStdEncoding.DecodeString(parsedHash[len(parsedHash)-2])
	if err != nil {
		return nil, "", err
	}
	return salt, parsedHash[len(parsedHash)-1], nil

}

func compareHash(user model.User, password string) error {
	salt, _, err := ParseHash(user.Hash)

	if err != nil {
		return err
	}

	if user.Hash != GenerateHash(password, []byte(salt)) {
		return errors.New("Invalid Password")
	}

	return nil
}
