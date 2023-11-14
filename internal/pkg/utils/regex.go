package utils

import (
	regexp "github.com/dlclark/regexp2"
)

const (
	userNameRegexPattern = "^[a-zA-Z0-9_-]{8,24}$"
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

// IsUserNameValid check if userName is valid
func IsUserNameValid(userName string) bool {
	userNameExp := regexp.MustCompile(userNameRegexPattern, regexp.None)
	result, _ := userNameExp.MatchString(userName)
	return result
}

// IsEmailValid check if email is valid
func IsEmailValid(email string) bool {
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	result, _ := emailExp.MatchString(email)
	return result
}

// IsPasswordValid check if password is valid
func IsPasswordValid(password string) bool {
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	result, _ := passwordExp.MatchString(password)
	return result
}
