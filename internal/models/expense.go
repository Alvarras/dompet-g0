package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Expense struct {
	ID          uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:char(36);not null" json:"user_id"`
	BudgetID    uuid.UUID      `gorm:"type:char(36);not null" json:"budget_id"`
	Amount      float64        `gorm:"not null" json:"amount"`
	Description string         `json:"description"`
	Date        time.Time      `gorm:"not null" json:"date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User           `gorm:"foreignKey:UserID" json:"user"`
	Budget      Budget         `gorm:"foreignKey:BudgetID" json:"budget"`
}
