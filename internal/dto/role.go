package dto

import "mzhn/auth/internal/entity"

type AddRoles struct {
	UserId string
	Roles  []entity.Role
}

type RemoveRoles struct {
	UserId string
	Roles  []entity.Role
}

type CheckRoles struct {
	UserId string
	Roles  []entity.Role
}
