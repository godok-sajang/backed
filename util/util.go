package util

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Contain(s string, strArray []string) bool {
	for _, val := range strArray {
		if strings.Contains(s, val) {
			return true
		}
	}
	return false
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
