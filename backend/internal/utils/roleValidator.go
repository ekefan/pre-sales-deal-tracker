package utils

import "github.com/go-playground/validator/v10"

const (
	Admin   = "admin"
	Sales   = "sales"
	Manager = "manager"
)

func RoleValidator(fl validator.FieldLevel) bool {
	if role, ok := fl.Field().Interface().(string); ok {
		return isSupportedRole(role)
	}
	return false
}

func isSupportedRole(role string) bool {
	switch role {
	case Admin, Sales, Manager:
		return true
	}
	return false
}
