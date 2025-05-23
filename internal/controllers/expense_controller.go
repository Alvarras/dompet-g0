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

type ExpenseController struct {
	expenseService *services.ExpenseService
	validate       *validator.Validate
}

func NewExpenseController(expenseService *services.ExpenseService) *ExpenseController {
	return &ExpenseController{
		expenseService: expenseService,
		validate:       validator.New(),
	}
}

func (c *ExpenseController) CreateExpense(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)

	var req requests.CreateExpenseRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "EXPENSE_001"))
	}

	if err := c.validate.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "EXPENSE_002"))
	}

	response, err := c.expenseService.CreateExpense(userID, &req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "EXPENSE_003"))
	}

	return ctx.JSON(http.StatusCreated, responses.NewSuccessResponse(response))
}

func (c *ExpenseController) GetExpenses(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)

	response, err := c.expenseService.GetExpenses(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(err.Error(), "EXPENSE_004"))
	}

	return ctx.JSON(http.StatusOK, responses.NewSuccessResponse(response))
}

func (c *ExpenseController) GetExpensesByBudget(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)
	budgetID, err := uuid.Parse(ctx.Param("budget_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse("invalid budget id", "EXPENSE_005"))
	}

	response, err := c.expenseService.GetExpensesByBudget(userID, budgetID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "EXPENSE_006"))
	}

	return ctx.JSON(http.StatusOK, responses.NewSuccessResponse(response))
}

func (c *ExpenseController) DeleteExpense(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)
	expenseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse("invalid expense id", "EXPENSE_007"))
	}

	if err := c.expenseService.DeleteExpense(userID, expenseID); err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "EXPENSE_008"))
	}

	return ctx.NoContent(http.StatusNoContent)
}
