package validator

import (
	"github.com/go-playground/validator/v10"
)

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

var ValidCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		// check currency is supported
		return IsSupportedCurrency(currency)
	}
	return false
}

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
