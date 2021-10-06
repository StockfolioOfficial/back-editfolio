package di

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
)

func newValidator() (v *validator.Validate) {
	v = validator.New()
	v.RegisterValidation("sf_mobile", mobileValidation)
	return
}

var (
	mobileRegex = regexp.MustCompile("^010\\d{8}$")
)

func mobileValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	return mobileRegex.MatchString(field.String())
}