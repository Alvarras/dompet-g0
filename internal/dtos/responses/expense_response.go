package responses

import (
	"time"

	"github.com/google/uuid"
)

type ExpenseResponse struct {
	ID          uuid.UUID `json:"id"`
	BudgetID    uuid.UUID `json:"budget_id"`
	BudgetName  string    `json:"budget_name"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type ExpenseListResponse struct {
	Expenses []ExpenseResponse `json:"expenses"`
	Total    int               `json:"total"`
}
