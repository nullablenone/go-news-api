package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/internal/domain/article"
	"github.com/nullablenone/go-news-api/internal/domain/user"
	"github.com/nullablenone/go-news-api/internal/shared/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresAndRoutesIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	wd, err := os.Getwd()
	require.NoError(t, err)
	rootDir := filepath.Clean(filepath.Join(wd, "..", ".."))
	require.NoError(t, os.Chdir(rootDir))
	defer func() { _ = os.Chdir(wd) }()

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Skipf("skipping integration test: %v", err)
	}

	db, err := config.ConnectPostgre(cfg.DB)
	if err != nil {
		t.Skipf("skipping integration test: PostgreSQL unavailable: %v", err)
	}

	require.NoError(t, db.AutoMigrate(article.Article{}, user.User{}, user.RefreshToken{}))

	// Run inside a transaction that always rolls back so the database stays clean.
	tx := db.Begin()
	require.NoError(t, tx.Error)
	defer tx.Rollback()

	user.RunAdminSeeder(tx)

	articleRepo := article.NewArticleRepository(tx)
	articleService := article.NewArticleService(articleRepo, nil)
	articleHandler := article.NewArticleHandler(articleService)

	userRepo := user.NewUserRepository(tx)
	userService := user.NewUserService(userRepo, cfg)
	userHandler := user.NewUserHandler(userService)

	router := SetRoutes(cfg, articleHandler, userHandler)

	// 1. Register user
	regBody, _ := json.Marshal(map[string]string{
		"name":     "Tester",
		"email":    "tester@example.com",
		"password": "password123",
	})
	w := performRequest(router, http.MethodPost, "/register", regBody, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Login as seeded admin
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "admin@news.com",
		"password": "#admin123",
	})
	w = performRequest(router, http.MethodPost, "/login", loginBody, "")
	require.Equal(t, http.StatusOK, w.Code)

	var loginResp struct {
		Data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &loginResp))
	assert.NotEmpty(t, loginResp.Data.AccessToken)
	assert.NotEmpty(t, loginResp.Data.RefreshToken)

	// 3. Refresh token rotation
	w = performRequest(router, http.MethodPost, "/auth/refresh", nil, loginResp.Data.RefreshToken)
	require.Equal(t, http.StatusOK, w.Code)

	// 4. Create article as admin
	createBody, _ := json.Marshal(map[string]interface{}{
		"title":     "Integration Test Article",
		"slug":      "integration-test-article",
		"summary":   "Testing DB and routes",
		"category":  "Testing",
		"read_time": "2 min",
		"content":   []map[string]interface{}{{"type": "paragraph", "text": "Body"}},
	})
	w = performRequest(router, http.MethodPost, "/admin/article", createBody, loginResp.Data.AccessToken)
	require.Equal(t, http.StatusCreated, w.Code)

	var createResp struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &createResp))
	require.NotZero(t, createResp.Data.ID)

	// 5. Public list & slug lookup
	w = performRequest(router, http.MethodGet, "/public/articles", nil, "")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "integration-test-article")

	w = performRequest(router, http.MethodGet, "/public/articles/integration-test-article", nil, "")
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. Admin get by ID, update, delete
	w = performRequest(router, http.MethodGet, "/admin/articles", nil, loginResp.Data.AccessToken)
	assert.Equal(t, http.StatusOK, w.Code)

	articleIDPath := fmt.Sprintf("/admin/article/%d", createResp.Data.ID)
	w = performRequest(router, http.MethodGet, articleIDPath, nil, loginResp.Data.AccessToken)
	assert.Equal(t, http.StatusOK, w.Code)

	updateBody, _ := json.Marshal(map[string]interface{}{
		"title":     "Updated Title",
		"slug":      "integration-test-article-updated",
		"summary":   "Updated summary",
		"category":  "Testing",
		"read_time": "3 min",
		"content":   []map[string]interface{}{{"type": "paragraph", "text": "Updated"}},
	})
	w = performRequest(router, http.MethodPut, articleIDPath, updateBody, loginResp.Data.AccessToken)
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequest(router, http.MethodDelete, articleIDPath, nil, loginResp.Data.AccessToken)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. Docs & OpenAPI endpoints
	w = performRequest(router, http.MethodGet, "/docs", nil, "")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SwaggerUIBundle")

	w = performRequest(router, http.MethodGet, "/docs/openapi.yaml", nil, "")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "openapi: 3.0.3")
}

func performRequest(router http.Handler, method, path string, body []byte, bearerToken string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
