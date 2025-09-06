package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashedPassword(hashedPassword string, password string) bool {
	fmt.Println("hashed Password -> ", hashedPassword, " \n password -> ", password)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	fmt.Println("COMPARE -> ", err)
	return err == nil
}
