package services

import (
	"errors"

	"github.com/Alvarras/dompet-g0/internal/dtos/requests"
	"github.com/Alvarras/dompet-g0/internal/dtos/responses"
	"github.com/Alvarras/dompet-g0/internal/models"
	"github.com/Alvarras/dompet-g0/internal/repositories"
	"github.com/google/uuid"
)

type BudgetService struct {
	budgetRepo *repositories.BudgetRepository
}

func NewBudgetService(budgetRepo *repositories.BudgetRepository) *BudgetService {
	return &BudgetService{
		budgetRepo: budgetRepo,
	}
}

func (s *BudgetService) CreateBudget(userID uuid.UUID, req *requests.CreateBudgetRequest) (*responses.BudgetResponse, error) {
	budget := &models.Budget{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		Amount:      req.Amount,
		Description: req.Description,
	}

	if err := s.budgetRepo.Create(budget); err != nil {
		return nil, err
	}

	return &responses.BudgetResponse{
		ID:          budget.ID,
		Name:        budget.Name,
		Amount:      budget.Amount,
		Spent:       budget.Spent,
		Remaining:   budget.Amount - budget.Spent,
		Description: budget.Description,
	}, nil
}

func (s *BudgetService) GetBudgets(userID uuid.UUID) (*responses.BudgetListResponse, error) {
	budgets, err := s.budgetRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var budgetResponses []responses.BudgetResponse
	for _, budget := range budgets {
		budgetResponses = append(budgetResponses, responses.BudgetResponse{
			ID:          budget.ID,
			Name:        budget.Name,
			Amount:      budget.Amount,
			Spent:       budget.Spent,
			Remaining:   budget.Amount - budget.Spent,
			Description: budget.Description,
		})
	}

	return &responses.BudgetListResponse{
		Budgets: budgetResponses,
		Total:   len(budgetResponses),
	}, nil
}

func (s *BudgetService) UpdateBudget(userID uuid.UUID, budgetID uuid.UUID, req *requests.UpdateBudgetRequest) (*responses.BudgetResponse, error) {
	budget, err := s.budgetRepo.FindByID(budgetID)
	if err != nil {
		return nil, err
	}

	if budget.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	budget.Name = req.Name
	budget.Amount = req.Amount
	budget.Description = req.Description

	if err := s.budgetRepo.Update(budget); err != nil {
		return nil, err
	}

	return &responses.BudgetResponse{
		ID:          budget.ID,
		Name:        budget.Name,
		Amount:      budget.Amount,
		Spent:       budget.Spent,
		Remaining:   budget.Amount - budget.Spent,
		Description: budget.Description,
	}, nil
}

func (s *BudgetService) DeleteBudget(userID uuid.UUID, budgetID uuid.UUID) error {
	budget, err := s.budgetRepo.FindByID(budgetID)
	if err != nil {
		return err
	}

	if budget.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.budgetRepo.Delete(budgetID)
}
