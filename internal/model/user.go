package model

// Struct which represents the user's information
type User struct {
	UserId   int
	Username string
	Salt     string
	Hash     string
}
