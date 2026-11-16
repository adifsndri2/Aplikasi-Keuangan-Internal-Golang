package model

import (
	"time"
)

type User struct {
	ID string `gorm:"primaryKey"`
	Name string
	Email string `gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-"` 
	Role string `gorm:"not null"` //admin dan user
	IsActive bool `gorm:"default:true"`
	CreatedAt time.Time
}
