package api

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(pwd string) (string, error) {
	pwd = strings.TrimSpace(strings.Trim(pwd, "\n"))
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	plainPwd = strings.TrimSpace(strings.Trim(plainPwd, "\n"))
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)) == nil
}
