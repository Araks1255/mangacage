package utils

import "golang.org/x/crypto/bcrypt"

func GenerateHashPassword(password string) (hash string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}
