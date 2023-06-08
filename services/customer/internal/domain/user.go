package domain

import "time"

type User struct {
	ID           int64
	Details      UserDetails
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type UserDetails struct {
	FirstName string
	LastName  string
	Phone     string
}
