package model

import "time"

type user struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-"`
	Role         string `gorm:"not null"` //admin atau user
	IsActive     bool   `gorm:"default:true"`
	CreatedAt    time.Time
}
