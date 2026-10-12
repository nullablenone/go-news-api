package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	secret := "my-secret-key"
	userID := uint(1)
	role := "admin"

	// 1. Access Token
	accessToken, err := GenerateAccessJWT(userID, role, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)

	claims, err := ValidateToken(accessToken, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, "access", claims.TokenType)

	// 2. Refresh Token
	refreshToken, err := GenerateRefreshJWT(userID, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)

	refreshClaims, err := ValidateToken(refreshToken, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, refreshClaims.UserID)
	assert.Equal(t, "refresh", refreshClaims.TokenType)

	// 3. Wrong Secret
	_, err = ValidateToken(accessToken, "wrong-secret")
	assert.Error(t, err)

	// 4. Expired Token
	expiredClaims := JWTClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString([]byte(secret))
	assert.NoError(t, err)

	_, err = ValidateToken(expiredToken, secret)
	assert.Error(t, err)
}
