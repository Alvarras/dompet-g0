package services

import (
	"errors"
	"time"

	"github.com/Alvarras/dompet-g0/internal/dtos/requests"
	"github.com/Alvarras/dompet-g0/internal/dtos/responses"
	"github.com/Alvarras/dompet-g0/internal/models"
	"github.com/Alvarras/dompet-g0/internal/repositories"
	"github.com/google/uuid"
)

type ExpenseService struct {
	expenseRepo *repositories.ExpenseRepository
	budgetRepo  *repositories.BudgetRepository
}

func NewExpenseService(expenseRepo *repositories.ExpenseRepository, budgetRepo *repositories.BudgetRepository) *ExpenseService {
	return &ExpenseService{
		expenseRepo: expenseRepo,
		budgetRepo:  budgetRepo,
	}
}

func (s *ExpenseService) CreateExpense(userID uuid.UUID, req *requests.CreateExpenseRequest) (*responses.ExpenseResponse, error) {
	// Check if budget exists and belongs to user
	budget, err := s.budgetRepo.FindByID(req.BudgetID)
	if err != nil {
		return nil, errors.New("budget not found")
	}

	if budget.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Check if there's enough budget
	if budget.Amount-budget.Spent < req.Amount {
		return nil, errors.New("insufficient budget")
	}

	// Set current time if no date is provided
	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	expense := &models.Expense{
		ID:          uuid.New(),
		UserID:      userID,
		BudgetID:    req.BudgetID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
	}

	if err := s.expenseRepo.Create(expense); err != nil {
		return nil, err
	}

	// Update budget spent amount
	if err := s.budgetRepo.UpdateSpent(req.BudgetID, req.Amount); err != nil {
		return nil, err
	}

	return &responses.ExpenseResponse{
		ID:          expense.ID,
		BudgetID:    expense.BudgetID,
		BudgetName:  budget.Name,
		Amount:      expense.Amount,
		Description: expense.Description,
		Date:        expense.Date,
	}, nil
}

func (s *ExpenseService) GetExpenses(userID uuid.UUID) (*responses.ExpenseListResponse, error) {
	expenses, err := s.expenseRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var expenseResponses []responses.ExpenseResponse
	for _, expense := range expenses {
		// Get budget information
		budget, err := s.budgetRepo.FindByID(expense.BudgetID)
		if err != nil {
			return nil, err
		}

		expenseResponses = append(expenseResponses, responses.ExpenseResponse{
			ID:              expense.ID,
			BudgetID:        expense.BudgetID,
			BudgetName:      expense.Budget.Name,
			Amount:          expense.Amount,
			Description:     expense.Description,
			Date:            expense.Date,
			BudgetRemaining: budget.Amount - budget.Spent,
			BudgetSpent:     budget.Spent,
			BudgetTotal:     budget.Amount,
		})
	}

	return &responses.ExpenseListResponse{
		Expenses: expenseResponses,
		Total:    len(expenseResponses),
	}, nil
}

func (s *ExpenseService) GetExpensesByBudget(userID uuid.UUID, budgetID uuid.UUID) (*responses.ExpenseListResponse, error) {
	// Check if budget belongs to user
	budget, err := s.budgetRepo.FindByID(budgetID)
	if err != nil {
		return nil, errors.New("budget not found")
	}

	if budget.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	expenses, err := s.expenseRepo.FindByBudgetID(budgetID)
	if err != nil {
		return nil, err
	}

	var expenseResponses []responses.ExpenseResponse
	for _, expense := range expenses {
		expenseResponses = append(expenseResponses, responses.ExpenseResponse{
			ID:              expense.ID,
			BudgetID:        expense.BudgetID,
			BudgetName:      expense.Budget.Name,
			Amount:          expense.Amount,
			Description:     expense.Description,
			Date:            expense.Date,
			BudgetRemaining: budget.Amount - budget.Spent,
			BudgetSpent:     budget.Spent,
			BudgetTotal:     budget.Amount,
		})
	}

	return &responses.ExpenseListResponse{
		Expenses: expenseResponses,
		Total:    len(expenseResponses),
	}, nil
}

func (s *ExpenseService) DeleteExpense(userID uuid.UUID, expenseID uuid.UUID) error {
	expense, err := s.expenseRepo.FindByID(expenseID)
	if err != nil {
		return err
	}

	if expense.UserID != userID {
		return errors.New("unauthorized")
	}

	// Update budget spent amount
	if err := s.budgetRepo.UpdateSpent(expense.BudgetID, -expense.Amount); err != nil {
		return err
	}

	return s.expenseRepo.Delete(expenseID)
}

func (s *ExpenseService) UpdateExpense(userID uuid.UUID, expenseID uuid.UUID, req *requests.UpdateExpenseRequest) (*responses.ExpenseResponse, error) {
	// Get existing expense
	expense, err := s.expenseRepo.FindByID(expenseID)
	if err != nil {
		return nil, errors.New("expense not found")
	}

	// Check if expense belongs to user
	if expense.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Check if budget exists and belongs to user
	budget, err := s.budgetRepo.FindByID(req.BudgetID)
	if err != nil {
		return nil, errors.New("budget not found")
	}

	if budget.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Calculate budget adjustment
	budgetAdjustment := req.Amount - expense.Amount

	// Check if there's enough budget for the adjustment
	if budget.Amount-budget.Spent < budgetAdjustment {
		return nil, errors.New("insufficient budget")
	}

	// Update expense
	expense.BudgetID = req.BudgetID
	expense.Amount = req.Amount
	expense.Description = req.Description
	if !req.Date.IsZero() {
		expense.Date = req.Date
	}

	if err := s.expenseRepo.Update(expense); err != nil {
		return nil, err
	}

	// Update budget spent amount
	if err := s.budgetRepo.UpdateSpent(req.BudgetID, budgetAdjustment); err != nil {
		return nil, err
	}

	return &responses.ExpenseResponse{
		ID:          expense.ID,
		BudgetID:    expense.BudgetID,
		BudgetName:  budget.Name,
		Amount:      expense.Amount,
		Description: expense.Description,
		Date:        expense.Date,
	}, nil
}
