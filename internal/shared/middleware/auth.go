package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nullablenone/go-news-api/internal/shared/response"
	"github.com/nullablenone/go-news-api/internal/shared/utils"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.AbortError(c, http.StatusUnauthorized, "access token not found, please login first", "UNAUTHORIZED")
			return
		}

		accessTokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if accessTokenString == authHeader {
			response.AbortError(c, http.StatusUnauthorized, "access token format must be 'Bearer <access_token>'", "UNAUTHORIZED")
			return
		}

		accessToken, err := jwt.ParseWithClaims(accessTokenString, &utils.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !accessToken.Valid {
			response.AbortError(c, http.StatusUnauthorized, "invalid or expired access token", "UNAUTHORIZED")
			return
		}

		claims, ok := accessToken.Claims.(*utils.JWTClaims)
		if !ok {
			response.AbortError(c, http.StatusUnauthorized, "failed to parse access token claims", "UNAUTHORIZED")
			return
		}

		if claims.TokenType != "access" {
			response.AbortError(c, http.StatusUnauthorized, "provided token is not an access token", "UNAUTHORIZED")
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RefreshTokenMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.AbortError(c, http.StatusUnauthorized, "refresh token not found, please login first", "UNAUTHORIZED")
			return
		}

		refreshTokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if refreshTokenString == authHeader {
			response.AbortError(c, http.StatusUnauthorized, "refresh token format must be 'Bearer <refresh_token>'", "UNAUTHORIZED")
			return
		}

		refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &utils.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !refreshToken.Valid {
			response.AbortError(c, http.StatusUnauthorized, "invalid or expired refresh token", "UNAUTHORIZED")
			return
		}

		claims, ok := refreshToken.Claims.(*utils.JWTClaims)
		if !ok {
			response.AbortError(c, http.StatusUnauthorized, "failed to parse refresh token claims", "UNAUTHORIZED")
			return
		}

		if claims.TokenType != "refresh" {
			response.AbortError(c, http.StatusUnauthorized, "provided token is not a refresh token", "UNAUTHORIZED")
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("refresh_token", refreshTokenString)

		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.AbortError(c, http.StatusForbidden, "access denied: role not found", "FORBIDDEN")
			return
		}

		userRole, ok := role.(string)
		if !ok {
			response.AbortError(c, http.StatusForbidden, "access denied: invalid role format", "FORBIDDEN")
			return
		}
		isAllowed := false

		for _, r := range allowedRoles {
			if userRole == r {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			response.AbortError(c, http.StatusForbidden, "access denied: insufficient permissions", "FORBIDDEN")
			return
		}

		c.Next()
	}
}
