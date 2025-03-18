package util

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Println("env could not load")
	}
	m.Run()
}

func TestCompareHash(t *testing.T) {
	hash := "$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNhbHQ$SRzrpBkxb+Cwwr5PQJL2pIGh9G59lfzlgOj3RRV73LKQYf2HycaaTY5yHimy7mnlWCY"
	// Test right password
	ok, err := CompareHashToPassword(hash, "password")
	if !ok || err != nil {
		t.Error(err)
	}

	// Test wrong password
	ok, err = CompareHashToPassword(hash, "worng_password")
	if ok || err != nil {
		t.Error(err)
	}

	// Test bad hash
	ok, err = CompareHashToPassword("weird$hash", "worng_password")

	if ok || err == nil {
		t.Error(err)
	}
}

func TestSaltGeneration(t *testing.T) {

	salt1, err := GenerateSalt()

	if err != nil {
		t.Error(err)
	}

	salt2, err := GenerateSalt()

	if err != nil {
		t.Error(err)
	}

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

	hash := GenerateHash("password", salt)

	expectedHash := "$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNhbHQ$SRzrpBkxb+Cwwr5PQJL2pIGh9G59lfzlgOj3RRV73LKQYf2HycaaTY5yHimy7mnlWCY"

	if hash != expectedHash {
		t.Errorf("The two hashes are not the same.\nFuncHash: %s\nConstHash: %s", hash, expectedHash)
	}

	// Test if running the hash again will generate the same hash
	hash2 := GenerateHash("password", salt)

	if hash != hash2 {
		t.Errorf("The regenerated hash is not the same as the first.\nHash1: %s\nHash2: %s", hash, hash2)
	}

	// Test that the hashs will not the be the same with different passwords

	hash3 := GenerateHash("Password", salt)

	if hash3 == hash2 {
		t.Errorf("Hash3 is colliding with hash2 when they are not supposed to be.\nHash3: %s\nHash2: %s", hash3, hash2)
	}

	// Test different salt different hash

	salt2 := []byte("saltSalt")

	hash4 := GenerateHash("password", salt2)

	if hash4 == hash {
		t.Errorf("Hash4 is not supposed to be equal to hash1 as they have different salts\nHash4: %s\nHash1: %s", hash4, hash)
	}

}
