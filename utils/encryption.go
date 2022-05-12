package utils

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func verifyPassword(userPassword string, providedPassword string) (bool, string) {
	userPasswordInByte := []byte(userPassword)
	providedPasswordInByte := []byte(providedPassword)
	err := bcrypt.CompareHashAndPassword(userPasswordInByte, providedPasswordInByte)

	check := true
	var msg string

	if err != nil {
		msg = fmt.Sprintf("Invalid credentials")
		check = false
	}

	return check, msg
}
