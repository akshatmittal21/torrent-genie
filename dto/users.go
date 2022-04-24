package dto

import "time"

type APIUser struct {
	UserID    int64
	FirstName string
}

type DBUser struct {
	FirstName string
	LastName  string
	UserName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
