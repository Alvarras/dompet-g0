package repositories

import (
	"github.com/Alvarras/dompet-g0/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BudgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) Create(budget *models.Budget) error {
	return r.db.Create(budget).Error
}

func (r *BudgetRepository) FindByID(id uuid.UUID) (*models.Budget, error) {
	var budget models.Budget
	err := r.db.First(&budget, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *BudgetRepository) FindByUserID(userID uuid.UUID) ([]models.Budget, error) {
	var budgets []models.Budget
	err := r.db.Where("user_id = ?", userID).Find(&budgets).Error
	if err != nil {
		return nil, err
	}
	return budgets, nil
}

func (r *BudgetRepository) Update(budget *models.Budget) error {
	return r.db.Save(budget).Error
}

func (r *BudgetRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Budget{}, "id = ?", id).Error
}

func (r *BudgetRepository) UpdateSpent(id uuid.UUID, amount float64) error {
	return r.db.Model(&models.Budget{}).Where("id = ?", id).
		UpdateColumn("spent", gorm.Expr("spent + ?", amount)).Error
}
