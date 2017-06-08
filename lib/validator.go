package lib

import "regexp"

var emailRegexp = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

//ValidateEmail func validates an email
func ValidateEmail(email string) bool {
	return emailRegexp.MatchString(email)
}
