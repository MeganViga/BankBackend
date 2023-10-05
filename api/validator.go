package api

import (
	"github.com/MeganViga/BankBackend/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool{
	if currency , ok := fieldLevel.Field().Interface().(string);ok{
		return util.IsSupportedCurrency(currency)
	}
	return false
}

var validEmail validator.Func = func(fieldLevel validator.FieldLevel) bool{
	if email , ok := fieldLevel.Field().Interface().(string);ok{
		return util.IsEmail(email)
	}
	return false
}