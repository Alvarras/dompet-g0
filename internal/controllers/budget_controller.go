package controllers

import (
	"net/http"

	"github.com/Alvarras/dompet-g0/internal/dtos/requests"
	"github.com/Alvarras/dompet-g0/internal/dtos/responses"
	"github.com/Alvarras/dompet-g0/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BudgetController struct {
	budgetService *services.BudgetService
	validate      *validator.Validate
}

func NewBudgetController(budgetService *services.BudgetService) *BudgetController {
	return &BudgetController{
		budgetService: budgetService,
		validate:      validator.New(),
	}
}

func (c *BudgetController) CreateBudget(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)

	var req requests.CreateBudgetRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "BUDGET_001"))
	}

	if err := c.validate.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "BUDGET_002"))
	}

	response, err := c.budgetService.CreateBudget(userID, &req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "BUDGET_003"))
	}

	return ctx.JSON(http.StatusCreated, responses.NewSuccessResponse(response))
}

func (c *BudgetController) GetBudgets(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)

	response, err := c.budgetService.GetBudgets(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(err.Error(), "BUDGET_004"))
	}

	return ctx.JSON(http.StatusOK, responses.NewSuccessResponse(response))
}

func (c *BudgetController) UpdateBudget(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)
	budgetID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse("invalid budget id", "BUDGET_005"))
	}

	var req requests.UpdateBudgetRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "BUDGET_006"))
	}

	if err := c.validate.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "BUDGET_007"))
	}

	response, err := c.budgetService.UpdateBudget(userID, budgetID, &req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "BUDGET_008"))
	}

	return ctx.JSON(http.StatusOK, responses.NewSuccessResponse(response))
}

func (c *BudgetController) DeleteBudget(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)
	budgetID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse("invalid budget id", "BUDGET_009"))
	}

	if err := c.budgetService.DeleteBudget(userID, budgetID); err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "BUDGET_010"))
	}

	return ctx.NoContent(http.StatusNoContent)
}
