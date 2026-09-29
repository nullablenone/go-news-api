package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nullablenone/go-news-api/internal/utils"
)

// AuthMiddleware memvalidasi token JWT dari header Authorization
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access token tidak ditemukan, silakan login terlebih dahulu"})
			return
		}

		// Memisahkan kata "Bearer " dengan string token asli
		accessTokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if accessTokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format access token harus 'Bearer <access_token>'"})
			return
		}

		// Parsing dan validasi
		accessToken, err := jwt.ParseWithClaims(accessTokenString, &utils.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !accessToken.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access token tidak valid atau telah kedaluwarsa"})
			return
		}

		claims, ok := accessToken.Claims.(*utils.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Gagal membaca claims pada access token"})
			return
		}

		// Menyimpan data user_id dan role ke dalam konteks Gin
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware membatasi akses berdasarkan role tertentu
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: role tidak ditemukan"})
			return
		}

		userRole := role.(string)
		isAllowed := false

		// Cek apakah role user saat ini ada di dalam daftar role yang diizinkan
		for _, r := range allowedRoles {
			if userRole == r {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: Anda tidak memiliki hak akses"})
			return
		}

		c.Next()
	}
}
