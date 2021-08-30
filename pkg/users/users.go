package users

import (
	"errors"
)

// TODO:
// - Change ID to uint
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Example user for test
var users = []User{
	{
		ID:       "1",
		Username: "alice",
		Password: "alice",
	},
	{
		ID:       "2",
		Username: "bob",
		Password: "bob",
	},
}

func FindUserByID(id string) (User, error) {
	for _, u := range users {
		if id == u.ID {
			return u, nil
		}
	}

	return User{}, errors.New("Not found")
}

func FindUserByUsername(username string) (User, error) {
	for _, u := range users {
		if username == u.Username {
			return u, nil
		}
	}

	return User{}, errors.New("Not found")
}
