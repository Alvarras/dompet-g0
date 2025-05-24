package routes

import (
	"github.com/Alvarras/dompet-g0/internal/controllers"
	"github.com/Alvarras/dompet-g0/internal/middlewares"
	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(e *echo.Echo, jwtSecret string, authController *controllers.AuthController, budgetController *controllers.BudgetController, expenseController *controllers.ExpenseController) {
	// API version group
	v1 := e.Group("/api/v1")
	{
		// Public routes
		v1.POST("/register", authController.Register)
		v1.POST("/login", authController.Login)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middlewares.AuthMiddleware(jwtSecret))
		{
			// Budget routes
			budgets := protected.Group("/budgets")
			budgets.POST("", budgetController.CreateBudget)
			budgets.GET("", budgetController.GetBudgets)
			budgets.PUT("/:id", budgetController.UpdateBudget)
			budgets.DELETE("/:id", budgetController.DeleteBudget)

			// Expense routes
			expenses := protected.Group("/expenses")
			expenses.POST("", expenseController.CreateExpense)
			expenses.GET("", expenseController.GetExpenses)
			expenses.GET("/budget/:budget_id", expenseController.GetExpensesByBudget)
			expenses.PUT("/:id", expenseController.UpdateExpense)
			expenses.DELETE("/:id", expenseController.DeleteExpense)
		}
	}
}
