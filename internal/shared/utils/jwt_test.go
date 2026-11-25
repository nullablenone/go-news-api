package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidateAccessJWT(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateAccessJWT(42, "admin", secret)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateToken(token, secret)
	require.NoError(t, err)
	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, "admin", claims.Role)
	assert.Equal(t, "access", claims.TokenType)
}

func TestGenerateAndValidateRefreshJWT(t *testing.T) {
	secret := "test-secret"
	t1, err := GenerateRefreshJWT(7, secret)
	require.NoError(t, err)
	t2, err := GenerateRefreshJWT(7, secret)
	require.NoError(t, err)
	assert.NotEqual(t, t1, t2)

	claims, err := ValidateToken(t1, secret)
	require.NoError(t, err)
	assert.Equal(t, uint(7), claims.UserID)
	assert.Equal(t, "refresh", claims.TokenType)
	assert.NotEmpty(t, claims.ID)
}

func TestValidateTokenErrors(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateAccessJWT(1, "user", secret)
	require.NoError(t, err)

	_, err = ValidateToken(token, "wrong-secret")
	assert.Error(t, err)

	expiredClaims := JWTClaims{
		UserID:    1,
		Role:      "user",
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
		},
	}
	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString([]byte(secret))
	require.NoError(t, err)

	_, err = ValidateToken(expiredToken, secret)
	assert.Error(t, err)
}
