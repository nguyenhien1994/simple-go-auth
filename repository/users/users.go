// Fake user repository for testing
package users

import (
	"errors"
	"log"
)

// TODO:
// - Change ID to uint
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

// Example user for testing
var users = []User{
	{
		ID:       "1",
		Username: "alice",
		Password: "alice",
		Email:    "alice@gmail.com",
	},
	{
		ID:       "2",
		Username: "bob",
		Password: "bob",
		Email:    "bob@gmail.com",
	},
}

func FindUserByID(id string) (User, error) {
	for _, u := range users {
		if id == u.ID {
			return u, nil
		}
	}

	return User{}, errors.New("User not found")
}

func FindUserByUsername(username string) (User, error) {
	for _, u := range users {
		if username == u.Username {
			return u, nil
		}
	}

	return User{}, errors.New("User not found")
}

func UpdateUserEmail(id string, email string) error {
	for i, u := range users {
		if id == u.ID {
			users[i].Email = email
			return nil
		}
	}

	return errors.New("User not found")
}