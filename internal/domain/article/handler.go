package article

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/internal/shared/response"
)

type ArticleHandler struct {
	service ArticleService
}

func NewArticleHandler(service ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized user", "UNAUTHORIZED")
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid user ID format", "UNAUTHORIZED")
		return
	}

	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request data: "+err.Error(), "VALIDATION_ERROR")
		return
	}

	article, err := h.service.CreateArticle(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to save article", "INTERNAL_ERROR")
		return
	}

	response.Success(c, http.StatusCreated, "article created successfully", ToArticleResponse(article))
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid ID", "BAD_REQUEST")
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request data: "+err.Error(), "VALIDATION_ERROR")
		return
	}

	article, err := h.service.UpdateArticle(c.Request.Context(), uint(id), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	response.Success(c, http.StatusOK, "article updated successfully", ToArticleResponse(article))
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
	articles, err := h.service.GetAllArticles(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch articles", "INTERNAL_ERROR")
		return
	}

	items := make([]ArticleResponse, 0, len(articles))
	for _, a := range articles {
		items = append(items, ToArticleResponse(a))
	}

	response.Success(c, http.StatusOK, "articles fetched successfully", items)
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid ID", "BAD_REQUEST")
		return
	}

	article, err := h.service.GetArticleByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "article not found", "NOT_FOUND")
		return
	}

	response.Success(c, http.StatusOK, "article fetched successfully", ToArticleResponse(article))
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid ID", "BAD_REQUEST")
		return
	}

	err = h.service.DeleteArticle(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to delete article", "INTERNAL_ERROR")
		return
	}

	response.Success(c, http.StatusOK, "article deleted successfully", nil)
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	articles, err := h.service.GetAllArticles(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch articles", "INTERNAL_ERROR")
		return
	}

	items := make([]ArticleResponse, 0, len(articles))
	for _, a := range articles {
		items = append(items, ToArticleResponse(a))
	}

	response.Success(c, http.StatusOK, "articles fetched successfully", items)
}

func (h *ArticleHandler) GetArticleBySlug(c *gin.Context) {
	slug := c.Param("slug")

	article, err := h.service.GetArticleBySlug(c.Request.Context(), slug)
	if err != nil {
		response.Error(c, http.StatusNotFound, "article not found", "NOT_FOUND")
		return
	}

	response.Success(c, http.StatusOK, "article fetched successfully", ToArticleResponse(article))
}
