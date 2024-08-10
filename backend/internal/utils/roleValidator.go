package utils

import "github.com/go-playground/validator/v10"

const (
	AdminRole   = "admin"
	SalesRole   = "sales"
	ManagerRole = "manager"
)

func RoleValidator(fl validator.FieldLevel) bool {
	if role, ok := fl.Field().Interface().(string); ok {
		return isSupportedRole(role)
	}
	return false
}

func isSupportedRole(role string) bool {
	switch role {
	case AdminRole, SalesRole, ManagerRole:
		return true
	}
	return false
}
