package validator

import (
	"net/mail"
	"regexp"
	"strings"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-z0-9_]{3,20}$`)
)

func Email(email string) bool {
	if email == "" || len(email) > 255 {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

func Username(username string) bool {
	return usernameRegex.MatchString(username)
}

func Password(password string) bool {
	return len(password) >= 8 && len(password) <= 24
}

func MaxLength(value string, max int) bool {
	return len(strings.TrimSpace(value)) <= max
}

func MinLength(value string, min int) bool {
	return len(strings.TrimSpace(value)) >= min
}
