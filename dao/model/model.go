package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        string `gorm:"primaryKey;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
