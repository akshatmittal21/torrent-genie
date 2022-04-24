package gormdb

import (
	"time"

	"gorm.io/gorm"
)

type userConfig struct {
	UserID    int64 `gorm:"index,primary_key"`
	FirstName string
	LastName  string
	UserName  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
