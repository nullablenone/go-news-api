package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	service ArticleService
}

func NewArticleHandler(service ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req CreateArticleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Data tidak valid",
			"detail": err.Error(),
		})
		return
	}

	article, err := h.service.CreateArticle(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan artikel",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Artikel berhasil dibuat",
		"data": gin.H{
			"id":         article.ID,
			"title":      article.Title,
			"slug":       article.Slug,
			"content":    article.Content,
			"created_at": article.CreatedAt,
		},
	})
}
