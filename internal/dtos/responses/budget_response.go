package responses

import "github.com/google/uuid"

type BudgetResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Amount      float64   `json:"amount"`
	Spent       float64   `json:"spent"`
	Remaining   float64   `json:"remaining"`
	Description string    `json:"description"`
}

type BudgetListResponse struct {
	Budgets []BudgetResponse `json:"budgets"`
	Total   int              `json:"total"`
}
