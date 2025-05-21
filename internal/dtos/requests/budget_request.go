package requests

type CreateBudgetRequest struct {
	Name        string  `json:"name" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description"`
}

type UpdateBudgetRequest struct {
	Name        string  `json:"name" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description"`
}
