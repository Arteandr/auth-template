package dto

import "mzhn/auth/internal/entity"

type CreateUser struct {
	LastName   *string
	FirstName  *string
	MiddleName *string
	Email      string
	Password   string
	Roles      []entity.Role
}
