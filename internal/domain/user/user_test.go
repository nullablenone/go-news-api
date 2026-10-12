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
	users  map[uint]User
	tokens map[string]RefreshToken
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

func (m *mockUserRepository) CreateRefreshToken(token *RefreshToken) error {
	if m.tokens == nil {
		m.tokens = make(map[string]RefreshToken)
	}
	m.tokens[token.Token] = *token
	return nil
}

func (m *mockUserRepository) FindRefreshToken(tokenString string) (*RefreshToken, error) {
	t, ok := m.tokens[tokenString]
	if !ok {
		return nil, errors.New("token not found")
	}
	return &t, nil
}

func (m *mockUserRepository) RevokeRefreshToken(tokenString string) error {
	t, ok := m.tokens[tokenString]
	if !ok {
		return errors.New("token not found")
	}
	t.IsRevoked = true
	m.tokens[tokenString] = t
	return nil
}

func (m *mockUserRepository) RevokeAllUserTokens(userID uint) error {
	for k, t := range m.tokens {
		if t.UserID == userID {
			t.IsRevoked = true
			m.tokens[k] = t
		}
	}
	return nil
}

func TestRefreshTokenService(t *testing.T) {
	secret := "test-secret"
	env := &config.Env{JWTSecret: secret}
	mockRepo := &mockUserRepository{
		users: map[uint]User{
			1: {ID: 1, Name: "Admin", Email: "admin@test.com", Role: "admin"},
		},
		tokens: make(map[string]RefreshToken),
	}
	service := NewUserService(mockRepo, env)

	// Persiapkan refresh token awal di repository
	initialRefreshToken, err := utils.GenerateRefreshJWT(1, secret)
	assert.NoError(t, err)
	_ = mockRepo.CreateRefreshToken(&RefreshToken{
		UserID:    1,
		Token:     initialRefreshToken,
		IsRevoked: false,
	})

	t.Run("Success Refresh Token with RTR", func(t *testing.T) {
		newAccess, newRefresh, err := service.RefreshToken(1, initialRefreshToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, newAccess)
		assert.NotEmpty(t, newRefresh)
		assert.NotEqual(t, initialRefreshToken, newRefresh)

		// 1. Cek token lama sudah di-revoke
		oldToken, err := mockRepo.FindRefreshToken(initialRefreshToken)
		assert.NoError(t, err)
		assert.True(t, oldToken.IsRevoked, "token lama harus sudah di-revoke")

		// 2. Cek token baru tersimpan dan aktif
		newToken, err := mockRepo.FindRefreshToken(newRefresh)
		assert.NoError(t, err)
		assert.False(t, newToken.IsRevoked, "token baru harus aktif")

		// 3. Verifikasi claims
		accessClaims, err := utils.ValidateToken(newAccess, secret)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), accessClaims.UserID)
		assert.Equal(t, "admin", accessClaims.Role)
		assert.Equal(t, "access", accessClaims.TokenType)
	})

	t.Run("Trigger Reuse Detection when reusing revoked old token", func(t *testing.T) {
		// Coba gunakan kembali token lama yang sudah di-revoke
		_, _, err := service.RefreshToken(1, initialRefreshToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "terdeteksi penggunaan ulang")

		// Pastikan semua token milik user 1 sekarang di-revoke
		for _, tok := range mockRepo.tokens {
			if tok.UserID == 1 {
				assert.True(t, tok.IsRevoked, "seluruh sesi user harus di-revoke saat reuse terdeteksi")
			}
		}
	})

	t.Run("Fail when token not found in repository", func(t *testing.T) {
		_, _, err := service.RefreshToken(1, "unregistered-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "refresh token tidak dikenali")
	})

	t.Run("Fail when user not found in repository", func(t *testing.T) {
		_, _, err := service.RefreshToken(999, initialRefreshToken)
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
		tokens: make(map[string]RefreshToken),
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

	t.Run("Valid Header and Refresh Token with Rotation", func(t *testing.T) {
		refreshToken, err := utils.GenerateRefreshJWT(1, secret)
		assert.NoError(t, err)

		// Simpan token ke mock DB
		_ = mockRepo.CreateRefreshToken(&RefreshToken{
			UserID:    1,
			Token:     refreshToken,
			IsRevoked: false,
		})

		// Request pertama: Sukses diperbarui
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
		assert.NotEqual(t, refreshToken, resp.Data.RefreshToken)

		// Request kedua menggunakan token lama yang sama: Harus GAGAL (401) karena sudah di-revoke!
		reqReuse, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
		reqReuse.Header.Set("Authorization", "Bearer "+refreshToken)
		recReuse := httptest.NewRecorder()
		router.ServeHTTP(recReuse, reqReuse)

		assert.Equal(t, http.StatusUnauthorized, recReuse.Code)
		assert.Contains(t, recReuse.Body.String(), "terdeteksi penggunaan ulang")
	})
}
