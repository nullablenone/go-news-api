package user

import (
	"errors"
	"testing"

	"github.com/nullablenone/go-news-api/internal/shared/config"
	"github.com/nullablenone/go-news-api/internal/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRepo struct {
	users       map[uint]User
	usersByMail map[string]User
	tokens      map[string]*RefreshToken
	nextID      uint
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:       make(map[uint]User),
		usersByMail: make(map[string]User),
		tokens:      make(map[string]*RefreshToken),
		nextID:      1,
	}
}

func (m *mockUserRepo) Create(u *User) error {
	u.ID = m.nextID
	m.nextID++
	m.users[u.ID] = *u
	m.usersByMail[u.Email] = *u
	return nil
}

func (m *mockUserRepo) FindByEmail(email string) (User, error) {
	u, ok := m.usersByMail[email]
	if !ok {
		return User{}, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) FindByID(id uint) (User, error) {
	u, ok := m.users[id]
	if !ok {
		return User{}, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) CreateRefreshToken(token *RefreshToken) error {
	cp := *token
	m.tokens[token.Token] = &cp
	return nil
}

func (m *mockUserRepo) FindRefreshToken(tokenString string) (*RefreshToken, error) {
	t, ok := m.tokens[tokenString]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}

func (m *mockUserRepo) RevokeRefreshToken(tokenString string) error {
	if t, ok := m.tokens[tokenString]; ok {
		t.IsRevoked = true
	}
	return nil
}

func (m *mockUserRepo) RevokeAllUserTokens(userID uint) error {
	for _, t := range m.tokens {
		if t.UserID == userID {
			t.IsRevoked = true
		}
	}
	return nil
}

func TestUserServiceRegisterLoginAndRefreshRotation(t *testing.T) {
	repo := newMockUserRepo()
	cfg := &config.Config{JWT: config.JWTConfig{Secret: "test-secret"}}
	svc := NewUserService(repo, cfg)

	u, err := svc.Register(RegisterRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	})
	require.NoError(t, err)
	assert.Equal(t, "user", u.Role)
	assert.NotEqual(t, "password123", u.Password)

	_, _, err = svc.Login(LoginRequest{Email: "alice@example.com", Password: "wrong"})
	assert.Error(t, err)

	access1, refresh1, err := svc.Login(LoginRequest{Email: "alice@example.com", Password: "password123"})
	require.NoError(t, err)
	assert.NotEmpty(t, access1)
	assert.NotEmpty(t, refresh1)

	access2, refresh2, err := svc.RefreshToken(u.ID, refresh1)
	require.NoError(t, err)
	assert.NotEqual(t, refresh1, refresh2)

	claims, err := utils.ValidateToken(access2, cfg.JWT.Secret)
	require.NoError(t, err)
	assert.Equal(t, u.ID, claims.UserID)

	// Reusing an already-rotated refresh token must revoke all user tokens.
	_, _, err = svc.RefreshToken(u.ID, refresh1)
	assert.Error(t, err)
	assert.True(t, repo.tokens[refresh2].IsRevoked)
}
