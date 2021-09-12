// Fake user repository for testing
package users

import (
	"simple-go-auth/services/mysql"
	"strconv"

	"github.com/pkg/errors"
)

type User struct {
	ID        int
	Username  string
	Name      string
	Password  string
	Email     string
	Role      string
}

func parseUserFromDBData(rows [][]interface{}, columns []string) (User, error) {
	var user User
	for _, row := range rows {
		for columnIndex, columnValue := range row {
			b, ok := columnValue.([]byte)
			if !ok {
				return user, errors.New("Invalid data type retrieving from database")
			}
			switch columns[columnIndex] {
			case "id":
				id, err := strconv.Atoi(string(b))
				if err != nil {
					return user, errors.Wrap(err, "Can't parse id")
				}
				user.ID = id
			case "username":
				user.Username = string(b)
			case "name":
				user.Name = string(b)
			case "password_hash":
				user.Password = string(b)
			case "email":
				user.Email = string(b)
			case "role":
				user.Role = string(b)
			default:
				return user, errors.New("Invalid data type retrieving from database")
			}
		}
	}
	return user, nil
}

func FindUserByID(id int) (User, error) {
	query := "SELECT * from users where id = '" + strconv.Itoa(id) + "';"
	rows, columns, err := mysql.GetClientClient().ExecuteQuery(query)
	if err != nil {
		return User{}, err
	}

	return parseUserFromDBData(rows, columns)
}

func FindUserByUsername(username string) (User, error) {
	query := "SELECT * from users where username = '" + username + "';"
	rows, columns, err := mysql.GetClientClient().ExecuteQuery(query)
	if err != nil {
		return User{}, err
	}
	return parseUserFromDBData(rows, columns)
}

func UpdateUserInfo(id int, email string) error {
	return nil
}