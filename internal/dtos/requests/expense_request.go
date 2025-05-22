package requests

import (
	"time"

	"github.com/google/uuid"
)

type CreateExpenseRequest struct {
	BudgetID    uuid.UUID `json:"budget_id" validate:"required"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	Description string    `json:"description"`
	Date        time.Time `json:"date,omitempty"`
}

type UpdateExpenseRequest struct {
	BudgetID    uuid.UUID `json:"budget_id" validate:"required"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	Description string    `json:"description"`
	Date        time.Time `json:"date,omitempty"`
}
