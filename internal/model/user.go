package model

// Struct which represents the user's information
type User struct {
	UserId   int
	Username string
	Hash     string
}

func (user *User) Validate() map[string]string {
	problems := make(map[string]string)
	if user.UserId == 0 {
		problems["userId"] = "UserId is 0 which is invalid"
	}
	if user.Username == "" {
		problems["username"] = "Username is blank"
	}
	if user.Hash == "" {
		problems["hash"] = "There is no hash"
	}
	return problems
}
