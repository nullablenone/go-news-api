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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan, silakan login terlebih dahulu"})
			return
		}

		// Memisahkan kata "Bearer " dengan string token asli
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format token harus 'Bearer <token>'"})
			return
		}

		// Parsing dan validasi
		token, err := jwt.ParseWithClaims(tokenString, &utils.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau telah kedaluwarsa"})
			return
		}

		claims, ok := token.Claims.(*utils.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Gagal membaca claims pada token"})
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
