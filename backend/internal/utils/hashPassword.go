package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}
