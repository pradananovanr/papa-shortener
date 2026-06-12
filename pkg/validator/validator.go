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

// IsValidCustomURL checks if custom URL is alphanumeric (case-insensitive), with hyphen or underscore, and at least 3 characters
func IsValidCustomURL(url string) bool {
	if len(url) < 3 {
		return false
	}
	for _, c := range url {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return true
}