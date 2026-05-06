package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
}

func Get() *validator.Validate {
	return validate
}

// IsValidCustomURL checks if custom URL is alphanumeric, lowercase, and within length constraints
func IsValidCustomURL(url string) bool {
	if len(url) < 3 || len(url) > 20 {
		return false
	}
	for _, c := range url {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}