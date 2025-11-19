package validator

import "regexp"

var emailRegex = regexp.MustCompile(`^[\w._%+\-]+@[\w.\-]+\.[A-Za-z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
