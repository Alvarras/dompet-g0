package controllers

import (
	"net/http"

	"github.com/Alvarras/dompet-g0/internal/dtos/requests"
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response, err := c.expenseService.CreateExpense(userID, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (c *ExpenseController) GetExpenses(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)

	response, err := c.expenseService.GetExpenses(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *ExpenseController) GetExpensesByBudget(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)
	budgetID, err := uuid.Parse(ctx.Param("budget_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid budget id")
	}

	response, err := c.expenseService.GetExpensesByBudget(userID, budgetID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *ExpenseController) DeleteExpense(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)
	expenseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid expense id")
	}

	if err := c.expenseService.DeleteExpense(userID, expenseID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.NoContent(http.StatusNoContent)
}
