package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pswd string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
