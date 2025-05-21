package repositories

import (
	"github.com/Alvarras/dompet-g0/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

func (r *ExpenseRepository) FindByID(id uuid.UUID) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.Preload("Budget").First(&expense, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepository) FindByUserID(userID uuid.UUID) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Preload("Budget").Where("user_id = ?", userID).Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepository) FindByBudgetID(budgetID uuid.UUID) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Preload("Budget").Where("budget_id = ?", budgetID).Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepository) Update(expense *models.Expense) error {
	return r.db.Save(expense).Error
}

func (r *ExpenseRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Expense{}, "id = ?", id).Error
}
