package model

type User struct {
	Username string
	Salt     string
	Hash     string
}
