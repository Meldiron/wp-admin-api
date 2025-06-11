package config

import (
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

var users []map[string]string

func init() {
	usersEnv := os.Getenv("USERS")
	for _, user := range strings.Split(usersEnv, ",") {
		credentials := strings.Split(user, ":")
		if len(credentials) >= 2 {
			user := map[string]string{
				"username": credentials[0],
				"password": credentials[1],
			}
			users = append(users, user)
		}
	}
}

// ValidateCredentials checks if the provided username and password match a user in the system
func ValidateCredentials(username, password string) bool {
	for _, user := range users {
		// TODO: Add password hashing
		if user["username"] == username && user["password"] == password {
			return true
		}
	}

	return false
}
