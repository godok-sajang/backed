package util

import (
	"net/mail"
	"strings"
	"unicode"
)

func Contain(s string, strArray []string) bool {
	for _, val := range strArray {
		if strings.Contains(s, val) {
			return true
		}
	}
	return false
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidatePassword(password string) bool {
	if len(password) < 8 || len(password) > 50 {
		return false
	}
	var number, letter, special bool
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c) || c == ' ':
			special = true
		case unicode.IsLetter(c):
			letter = true
		}
	}
	return (number && letter) || (number && special) || (letter && special)
}
