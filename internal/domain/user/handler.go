package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/internal/shared/response"
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
		response.Error(c, http.StatusBadRequest, "invalid request data: "+err.Error(), "VALIDATION_ERROR")
		return
	}

	user, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to register: "+err.Error(), "INTERNAL_ERROR")
		return
	}

	response.Success(c, http.StatusCreated, "registration successful", ToRegisterResponse(user))
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request data: "+err.Error(), "VALIDATION_ERROR")
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), "UNAUTHORIZED")
		return
	}

	response.Success(c, http.StatusOK, "login successful", LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
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

	refreshTokenVal, exists := c.Get("refresh_token")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "refresh token not found", "UNAUTHORIZED")
		return
	}
	refreshTokenStr, ok := refreshTokenVal.(string)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid refresh token format", "UNAUTHORIZED")
		return
	}

	accessToken, refreshToken, err := h.service.RefreshToken(c.Request.Context(), userID, refreshTokenStr)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), "UNAUTHORIZED")
		return
	}

	response.Success(c, http.StatusOK, "token refreshed successfully", RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}


