package crypt

import "golang.org/x/crypto/bcrypt"

func Hash(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func Check(str, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}
