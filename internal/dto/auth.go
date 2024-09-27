package dto

import "mzhn/auth/internal/entity"

type Authenticate struct {
	AccessToken string
	Roles       []entity.Role
}

type Login struct {
	Email    string
	Password string
}

type Refresh struct {
	RefreshToken string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
