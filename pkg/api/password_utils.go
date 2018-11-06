package api

import "golang.org/x/crypto/bcrypt"

func hashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)) == nil
}
