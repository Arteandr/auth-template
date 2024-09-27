package entity

import "time"

type User struct {
	Id             string     `json:"id" db:"id"`
	LastName       *string    `json:"lastName" db:"last_name"`
	FirstName      *string    `json:"firstName" db:"first_name"`
	MiddleName     *string    `json:"middleName" db:"middle_name"`
	Email          string     `json:"email" db:"email"`
	HashedPassword string     `json:"password" db:"hashed_password"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      *time.Time `json:"updatedAt" db:"updated_at"`
}

type UserClaims struct {
	Id    string
	Email string
}
