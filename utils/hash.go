package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}
