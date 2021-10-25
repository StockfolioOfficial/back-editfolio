package di

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func newValidator() (v *validator.Validate) {
	v = validator.New()
	v.RegisterValidation("sf_mobile", mobileValidation)
	v.RegisterValidation("sf_password", passwordValidation)
	return
}

var (
	mobileRegex    = regexp.MustCompile("^010\\d{8}$")
	passwordRegex  = regexp.MustCompile("[A-Za-z]")
	passwordRegex1 = regexp.MustCompile("[0-9]")
)

func mobileValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	return mobileRegex.MatchString(field.String())
}

func passwordValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	if len(field.String()) < 8 || 32 < len(field.String()) {
		return false
	}

	return passwordRegex.MatchString(field.String()) && passwordRegex1.MatchString(field.String())
}
