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

// GetArticles untuk list semua artikel (Private)
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	articles, err := h.service.GetAllArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data artikel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": articles,
	})
}

// GetArticle untuk mengambil detail satu artikel berdasarkan ID 
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	article, err := h.service.GetArticleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artikel tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": article,
	})
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Data tidak valid",
			"detail": err.Error(),
		})
		return
	}

	article, err := h.service.UpdateArticle(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Artikel berhasil diperbarui",
		"data":    article,
	})
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	err = h.service.DeleteArticle(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus artikel atau data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Artikel berhasil dihapus",
	})
}
