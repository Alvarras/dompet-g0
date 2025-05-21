package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Budget struct {
	ID          uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:char(36);not null" json:"user_id"`
	Name        string         `gorm:"not null" json:"name"`
	Amount      float64        `gorm:"not null" json:"amount"`
	Spent       float64        `gorm:"default:0" json:"spent"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User           `gorm:"foreignKey:UserID" json:"user"`
}
