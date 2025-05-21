package services

import (
	"errors"
	"time"

	"github.com/Alvarras/dompet-g0/internal/dtos/requests"
	"github.com/Alvarras/dompet-g0/internal/dtos/responses"
	"github.com/Alvarras/dompet-g0/internal/models"
	"github.com/Alvarras/dompet-g0/internal/repositories"
	"github.com/Alvarras/dompet-g0/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repositories.UserRepository
	jwtSecret   string
	jwtDuration time.Duration
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string, jwtDuration time.Duration) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		jwtDuration: jwtDuration,
	}
}

func (s *AuthService) Register(req *requests.RegisterRequest) (*responses.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret, s.jwtDuration)
	if err != nil {
		return nil, err
	}

	return &responses.AuthResponse{
		Token: token,
		User: responses.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (s *AuthService) Login(req *requests.LoginRequest) (*responses.AuthResponse, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret, s.jwtDuration)
	if err != nil {
		return nil, err
	}

	return &responses.AuthResponse{
		Token: token,
		User: responses.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}
