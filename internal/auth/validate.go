package auth

import "regexp"

func IsValidUsername(username string) bool {
	usernameRe := regexp.MustCompile(`^[\w-]{3,32}$`) // usernames must be 3-32 characters long containing word characters or '-'
	return usernameRe.MatchString(username)
}

func IsValidPassword(password string) bool {
	passwordRe := regexp.MustCompile(`^\S{15,64}$`) // passwords must be 15-64 characters long containing non-space characters
	return passwordRe.MatchString(password)
}
