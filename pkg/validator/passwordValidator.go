package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var ValidPassword validator.Func = func(fl validator.FieldLevel) bool {
	if password, ok := fl.Field().Interface().(string); ok {
		// 密碼長度
		lengthRegex := regexp.MustCompile(`^.{6,16}$`)
		if !lengthRegex.MatchString(password) {
			return false
		}

		// 至少一小寫字母
		lowerRegex := regexp.MustCompile(`[a-z]`)
		if !lowerRegex.MatchString(password) {
			return false
		}

		// 至少一大寫字母
		upperRegex := regexp.MustCompile(`[A-Z]`)
		if !upperRegex.MatchString(password) {
			return false
		}

		// 至少一特殊符號
		specialRegex := regexp.MustCompile(`\W`)
		if !specialRegex.MatchString(password) {
			return false
		}
	}
	return true
}
