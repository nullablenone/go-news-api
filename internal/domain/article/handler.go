package article

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID format"})
		return
	}

	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data", "detail": err.Error()})
		return
	}

	article, err := h.service.CreateArticle(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save article"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "article created successfully",
		"data":    ToArticleResponse(article),
	})
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data", "detail": err.Error()})
		return
	}

	article, err := h.service.UpdateArticle(c.Request.Context(), uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "article updated successfully",
		"data":    ToArticleResponse(article),
	})
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
	articles, err := h.service.GetAllArticles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch articles"})
		return
	}

	var response []ArticleResponse
	for _, a := range articles {
		response = append(response, ToArticleResponse(a))
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	article, err := h.service.GetArticleByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ToArticleResponse(article)})
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = h.service.DeleteArticle(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "article deleted successfully"})
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	articles, err := h.service.GetAllArticles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch articles"})
		return
	}

	var response []ArticleResponse
	for _, a := range articles {
		response = append(response, ToArticleResponse(a))
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *ArticleHandler) GetArticleBySlug(c *gin.Context) {
	slug := c.Param("slug")

	article, err := h.service.GetArticleBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ToArticleResponse(article)})
}
