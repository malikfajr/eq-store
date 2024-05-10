package pkg

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func IsValidPhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()

	// Nomor telepon harus diawali dengan "+" dan diikuti dengan angka atau dash
	return strings.HasPrefix(phoneNumber, "+") && (strings.ContainsAny(phoneNumber[1:], "0123456789-"))
}

func ValidateURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()

	// Ekspresi reguler untuk memvalidasi URL
	regex := `^(https?://)?([a-zA-Z0-9-]+\.){1,}[a-zA-Z]{2,}(/[a-zA-Z0-9-._~:/?#[\]@!$&'()*+,;=]*)?$`
	match, _ := regexp.MatchString(regex, url)
	return match
}
