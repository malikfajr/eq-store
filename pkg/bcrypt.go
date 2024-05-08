package pkg

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

var BCRYPT_SALT int

func init() {
	bcrypt_salt_str := os.Getenv("BCRYPT_SALT")
	if salt, err := strconv.Atoi(bcrypt_salt_str); err != nil {
		BCRYPT_SALT = 8
	} else {
		BCRYPT_SALT = salt
	}
}

func HashPassword(password string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(password), BCRYPT_SALT)
	if err != nil {
		panic(err)
	}

	return string(h)
}

func ValidPassword(hash string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
