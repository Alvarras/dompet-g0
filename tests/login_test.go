package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Alvarras/dompet-g0/internal/controllers"
	"github.com/Alvarras/dompet-g0/internal/dtos/requests"
	"github.com/Alvarras/dompet-g0/internal/dtos/responses"
	"github.com/Alvarras/dompet-g0/internal/models"
	"github.com/Alvarras/dompet-g0/internal/repositories"
	"github.com/Alvarras/dompet-g0/internal/services"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	// Memuat variabel lingkungan dari .env.test.
	if err := godotenv.Load("../.env.test"); err != nil {
		fmt.Printf("Peringatan: .env.test tidak ditemukan: %v\n", err)
	}
}

func TestLogin(t *testing.T) {
	e, cleanupFunc := setupTestServer()
	defer cleanupFunc()

	t.Run("Success Login", func(t *testing.T) {
		loginReq := requests.LoginRequest{
			Email:    os.Getenv("TEST_USER_EMAIL"),
			Password: os.Getenv("TEST_USER_PASSWORD"),
		}
		reqBody, _ := json.Marshal(loginReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response responses.StandardResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)

		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.NotEmpty(t, data["token"])

		user, ok := data["user"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, os.Getenv("TEST_USER_EMAIL"), user["email"])
		assert.Equal(t, "Test User", user["name"])
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		loginReq := requests.LoginRequest{
			Email:    os.Getenv("TEST_USER_EMAIL"),
			Password: "wrongpassword",
		}
		reqBody, _ := json.Marshal(loginReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response responses.StandardResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Email atau password salah", response.Message)
		assert.Equal(t, "AUTH_006", response.Code)
	})

	t.Run("Missing Email", func(t *testing.T) {
		loginReq := requests.LoginRequest{
			Password: os.Getenv("TEST_USER_PASSWORD"),
		}
		reqBody, _ := json.Marshal(loginReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Missing Password", func(t *testing.T) {
		loginReq := requests.LoginRequest{
			Email: os.Getenv("TEST_USER_EMAIL"),
		}
		reqBody, _ := json.Marshal(loginReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		loginReq := requests.LoginRequest{
			Email:    "invalid-email",
			Password: os.Getenv("TEST_USER_PASSWORD"),
		}
		reqBody, _ := json.Marshal(loginReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// setupTestServer menginisialisasi server Echo untuk tes E2E dan fungsi cleanup.
func setupTestServer() (*echo.Echo, func()) {
	db, testUserID := setupTestDB()

	userRepo := repositories.NewUserRepository(db)
	jwtDuration, _ := time.ParseDuration(getEnv("JWT_EXPIRATION_TEST", "15m"))
	authService := services.NewAuthService(userRepo, getEnv("JWT_SECRET_TEST", "test-secret-key-e2e"), jwtDuration)
	authController := controllers.NewAuthController(authService)

	e := echo.New()
	e.POST("/api/v1/login", authController.Login)

	cleanup := func() {
		// Hapus pengguna tes untuk menjaga kebersihan antar test run.
		if err := db.Unscoped().Delete(&models.User{}, "id = ?", testUserID).Error; err != nil {
			fmt.Printf("Peringatan: Gagal menghapus pengguna tes (ID: %s): %v\n", testUserID, err)
		}
	}

	return e, cleanup
}

// setupTestDB mengkonfigurasi koneksi database tes, menjalankan migrasi, dan membuat pengguna tes.
// Mengembalikan instance GORM DB dan ID pengguna tes yang baru dibuat.
func setupTestDB() (*gorm.DB, uuid.UUID) {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "go_budget_test")
	testUserEmail := getEnv("TEST_USER_EMAIL", "test@example.com")
	testUserPassword := getEnv("TEST_USER_PASSWORD", "password123")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Kritis: Gagal terhubung ke database tes '%s': %v", dbName, err))
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic(fmt.Sprintf("Kritis: Gagal AutoMigrate User: %v", err))
	}

	testUserID := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testUserPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Sprintf("Kritis: Gagal hash password tes: %v", err))
	}

	// Hapus pengguna tes lama jika ada untuk konsistensi (mencegah error unique constraint).
	db.Unscoped().Where("email = ?", testUserEmail).Delete(&models.User{})

	testUser := models.User{
		ID:       testUserID,
		Email:    testUserEmail,
		Password: string(hashedPassword),
		Name:     "Test User",
	}
	if err := db.Create(&testUser).Error; err != nil {
		panic(fmt.Sprintf("Kritis: Gagal membuat pengguna tes: %v", err))
	}

	return db, testUserID
}

// getEnv mengambil variabel lingkungan atau mengembalikan nilai default.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
