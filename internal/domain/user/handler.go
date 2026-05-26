package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	user, err := h.service.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melakukan registrasi", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil",
		"data":    ToRegisterResponse(user),
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	token, err := h.service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"data": LoginResponse{
			Token: token,
		},
	})
}
