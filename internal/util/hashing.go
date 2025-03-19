package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Generates a randomized byte slice which length is determined by an env SALTLENGTH.
//
// Returns byte slice and nil if no errors.
// Returns nil and error if error occurs.
func GenerateSalt() ([]byte, error) {
	length, err := strconv.Atoi(os.Getenv("SALTLENGTH"))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	salt := make([]byte, length)

	_, err = rand.Read(salt)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return salt, nil
}

// Generates argon2id hash of the saltword with the salt. All parameters for argon2id are from an .env file.
//
// Returns the string of the hashed password.
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

// Parses the argon2id hash and returns the salt, base64 hashed password, and error if any.
func parseHash(hash string) ([]byte, string, error) {
	parsedHash := strings.Split(hash, "$")
	salt, err := base64.RawStdEncoding.DecodeString(parsedHash[len(parsedHash)-2])
	if err != nil {
		return nil, "", err
	}
	return salt, parsedHash[len(parsedHash)-1], nil

}

// Checks if the given password has the same hash as the given argon2id hash.
//
// Returns true if hash and password match.
// Returns false if they do not match.
func CompareHashToPassword(hash string, password string) (bool, error) {
	salt, _, err := parseHash(hash)

	if err != nil {
		return false, err
	}

	if hash != GenerateHash(password, []byte(salt)) {
		return false, nil
	}

	return true, nil
}
