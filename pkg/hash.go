package pkg

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(token string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error GenerateHash: %v\n", err)
		return "", err
	}
	return string(bytes), nil
}

func CompareHashAndData(token, hashTokenInDB string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashTokenInDB), []byte(token))
	if err != nil {
		log.Printf("Error CompareHashAndData: %v\n", err)
		return false
	}
	return true
}
