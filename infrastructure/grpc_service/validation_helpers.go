package grpcservice

import (
	"regexp"
)

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func validatePasswordMatch(password, confirmPassword string) error {
	if password != confirmPassword {
		return ErrPasswordMatchNotMatch
	}
	return nil
}
