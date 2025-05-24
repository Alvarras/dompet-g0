package responses

import (
	"time"

	"github.com/google/uuid"
)

// CreateExpenseResponse is used for POST /expenses response
type CreateExpenseResponse struct {
	ID          uuid.UUID `json:"id"`
	BudgetID    uuid.UUID `json:"budget_id"`
	BudgetName  string    `json:"budget_name"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

// UpdateExpenseResponse is used for PUT /expenses/:id response
type UpdateExpenseResponse struct {
	ID          uuid.UUID `json:"id"`
	BudgetID    uuid.UUID `json:"budget_id"`
	BudgetName  string    `json:"budget_name"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

// ExpenseResponse is used for GET responses
type ExpenseResponse struct {
	ID              uuid.UUID `json:"id"`
	BudgetID        uuid.UUID `json:"budget_id"`
	BudgetName      string    `json:"budget_name"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	Date            time.Time `json:"date"`
	BudgetRemaining float64   `json:"budget_remaining"`
	BudgetSpent     float64   `json:"budget_spent"`
	BudgetTotal     float64   `json:"budget_total"`
}

type ExpenseListResponse struct {
	Expenses []ExpenseResponse `json:"expenses"`
	Total    int               `json:"total"`
}
