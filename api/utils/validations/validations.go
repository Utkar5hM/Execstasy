package validations

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func UserNameValidator(fl validator.FieldLevel) bool {
	var nameRegex = regexp.MustCompile(`^[a-zA-Z]+([ '-][a-zA-Z]+)*$`)
	return nameRegex.MatchString(fl.Field().String())
}

func LinuxUserValidator(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if val == "*" {
		return true
	}
	var linuxUserRegex = regexp.MustCompile(`^[a-z_]([a-z0-9_-]{0,31}|[a-z0-9_-]{0,30}\$)$`)
	return linuxUserRegex.MatchString(fl.Field().String())
}

func ScopeLinuxUserValidator(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if !strings.HasPrefix(val, "user:") {
		return false
	}
	username := strings.TrimPrefix(val, "user:")
	var linuxUserRegex = regexp.MustCompile(`^[a-z_]([a-z0-9_-]{0,31}|[a-z0-9_-]{0,30}\$)$`)
	return linuxUserRegex.MatchString(username)
}

func NameValidator(fl validator.FieldLevel) bool {
	var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+([ '-][a-zA-Z0-9]+)*$`)
	return nameRegex.MatchString(fl.Field().String())
}
