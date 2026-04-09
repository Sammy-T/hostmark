package auth

import (
	"regexp"
	"slices"
)

func IsValidUsername(username string) bool {
	usernameRe := regexp.MustCompile(`^[\w-]{3,32}$`) // usernames must be 3-32 characters long containing word characters or '-'
	return usernameRe.MatchString(username)
}

func IsValidPassword(password string) bool {
	passwordRe := regexp.MustCompile(`^\S{15,64}$`) // passwords must be 15-64 characters long containing non-space characters
	return passwordRe.MatchString(password)
}

func IsValidRole(role string) bool {
	return slices.Contains([]Role{RoleUser, RoleAdmin}, Role(role))
}

func IsValidVisibility(visibility string) bool {
	return slices.Contains([]string{"public", "protected", "private"}, visibility)
}
