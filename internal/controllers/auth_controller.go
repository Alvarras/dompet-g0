package controllers

import (
	"net/http"

	"github.com/Alvarras/dompet-g0/internal/dtos/requests"
	"github.com/Alvarras/dompet-g0/internal/dtos/responses"
	"github.com/Alvarras/dompet-g0/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService *services.AuthService
	validate    *validator.Validate
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
		validate:    validator.New(),
	}
}

func (c *AuthController) Register(ctx echo.Context) error {
	var req requests.RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "AUTH_001"))
	}

	if err := c.validate.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "AUTH_002"))
	}

	response, err := c.authService.Register(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "AUTH_003"))
	}

	return ctx.JSON(http.StatusCreated, responses.NewSuccessResponse(response))
}

func (c *AuthController) Login(ctx echo.Context) error {
	var req requests.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "AUTH_004"))
	}

	if err := c.validate.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(err.Error(), "AUTH_005"))
	}

	response, err := c.authService.Login(&req)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Email atau password salah", "AUTH_006"))
	}

	return ctx.JSON(http.StatusOK, responses.NewSuccessResponse(response))
}
