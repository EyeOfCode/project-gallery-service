package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func SetupValidator() {
    validate.RegisterValidation("password_validator", PasswordValidator)
}

func PasswordValidator(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*()+-_=\[\]{}|;:,.<>?]`).MatchString(password)

    return hasUpper && hasNumber && hasSpecial && hasLower
}

func ValidateStruct(payload interface{}) error {
    return validate.Struct(payload)
}