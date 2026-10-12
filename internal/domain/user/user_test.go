package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/config"
	"github.com/nullablenone/go-news-api/internal/middleware"
	"github.com/nullablenone/go-news-api/internal/utils"
	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct {
	users map[uint]User
}

func (m *mockUserRepository) Create(user *User) error {
	m.users[user.ID] = *user
	return nil
}

func (m *mockUserRepository) FindByEmail(email string) (User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return User{}, errors.New("user not found")
}

func (m *mockUserRepository) FindByID(id uint) (User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return User{}, errors.New("user not found")
}

func TestRefreshTokenService(t *testing.T) {
	secret := "test-secret"
	env := &config.Env{JWTSecret: secret}
	mockRepo := &mockUserRepository{
		users: map[uint]User{
			1: {ID: 1, Name: "Admin", Email: "admin@test.com", Role: "admin"},
		},
	}
	service := NewUserService(mockRepo, env)

	t.Run("Success Refresh Token", func(t *testing.T) {
		newAccess, newRefresh, err := service.RefreshToken(1)
		assert.NoError(t, err)
		assert.NotEmpty(t, newAccess)
		assert.NotEmpty(t, newRefresh)

		// Verify new access token
		accessClaims, err := utils.ValidateToken(newAccess, secret)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), accessClaims.UserID)
		assert.Equal(t, "admin", accessClaims.Role)
		assert.Equal(t, "access", accessClaims.TokenType)

		// Verify new refresh token
		refreshClaims, err := utils.ValidateToken(newRefresh, secret)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), refreshClaims.UserID)
		assert.Equal(t, "refresh", refreshClaims.TokenType)
	})

	t.Run("Fail when user not found in repository", func(t *testing.T) {
		_, _, err := service.RefreshToken(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user tidak ditemukan")
	})
}

func TestRefreshTokenHandlerWithMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"
	env := &config.Env{JWTSecret: secret}
	mockRepo := &mockUserRepository{
		users: map[uint]User{
			1: {ID: 1, Name: "User One", Email: "user@test.com", Role: "user"},
		},
	}
	service := NewUserService(mockRepo, env)
	handler := NewUserHandler(service)

	router := gin.New()
	authRoutes := router.Group("/auth")
	authRoutes.Use(middleware.RefreshTokenMiddleware(secret))
	{
		authRoutes.POST("/refresh", handler.RefreshToken)
	}

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Refresh token tidak ditemukan")
	})

	t.Run("Invalid Format Header (No Bearer)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
		req.Header.Set("Authorization", "InvalidFormatToken")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		var resp map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Format refresh token harus 'Bearer <refresh_token>'", resp["error"])
	})

	t.Run("Token is access token instead of refresh token", func(t *testing.T) {
		accessToken, err := utils.GenerateAccessJWT(1, "user", secret)
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "token yang digunakan bukan refresh token")
	})

	t.Run("Valid Header and Refresh Token", func(t *testing.T) {
		refreshToken, err := utils.GenerateRefreshJWT(1, secret)
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
		req.Header.Set("Authorization", "Bearer "+refreshToken)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp struct {
			Message string               `json:"message"`
			Data    RefreshTokenResponse `json:"data"`
		}
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Token berhasil diperbarui", resp.Message)
		assert.NotEmpty(t, resp.Data.AccessToken)
		assert.NotEmpty(t, resp.Data.RefreshToken)
	})
}
