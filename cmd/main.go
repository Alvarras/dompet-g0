package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Alvarras/dompet-g0/internal/config"
	"github.com/Alvarras/dompet-g0/internal/controllers"
	"github.com/Alvarras/dompet-g0/internal/models"
	"github.com/Alvarras/dompet-g0/internal/repositories"
	"github.com/Alvarras/dompet-g0/internal/routes"
	"github.com/Alvarras/dompet-g0/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.Charset,
		cfg.Database.ParseTime,
		cfg.Database.Loc,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database
	if err := db.AutoMigrate(&models.User{}, &models.Budget{}, &models.Expense{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	budgetRepo := repositories.NewBudgetRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)

	// Initialize services
	jwtDuration, _ := time.ParseDuration(cfg.JWT.Expiration)
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, jwtDuration)
	budgetService := services.NewBudgetService(budgetRepo)
	expenseService := services.NewExpenseService(expenseRepo, budgetRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	budgetController := controllers.NewBudgetController(budgetService)
	expenseController := controllers.NewExpenseController(expenseService)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Setup routes
	routes.SetupRoutes(e, cfg.JWT.Secret, authController, budgetController, expenseController)

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := e.Start(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
