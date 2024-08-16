package util

import (
	"regexp"
)

var (
	checkCharacters = regexp.MustCompile(`^[a-zA-Z_\-0-9]*$`)
)

func ValidateUsername(username string) (bool, string) {
	if len(username) < 3 || len(username) > 25 {
		return false, "username must be between 3 and 25 characters"
	} else if !checkCharacters.MatchString(username) {
		return false, "username must only contain letters, numbers and - or _"
	}
	return true, ""
}
