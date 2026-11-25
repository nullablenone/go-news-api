package user

import (
	"errors"
	"time"

	"github.com/nullablenone/go-news-api/internal/shared/config"
	"github.com/nullablenone/go-news-api/internal/shared/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req RegisterRequest) (User, error)
	Login(req LoginRequest) (string, string, error)
	RefreshToken(userID uint, oldRefreshToken string) (string, string, error)
}

type userService struct {
	repo UserRepository
	cfg  *config.Config
}

func NewUserService(repo UserRepository, cfg *config.Config) UserService {
	return &userService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *userService) Register(req RegisterRequest) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user := User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err = s.repo.Create(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *userService) Login(req LoginRequest) (string, string, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	accessToken, err := utils.GenerateAccessJWT(user.ID, user.Role, s.cfg.JWT.Secret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshJWT(user.ID, s.cfg.JWT.Secret)
	if err != nil {
		return "", "", err
	}

	refreshTokenModel := RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.repo.CreateRefreshToken(&refreshTokenModel); err != nil {
		return "", "", errors.New("failed to save refresh token")
	}

	return accessToken, refreshToken, nil
}

func (s *userService) RefreshToken(userID uint, oldRefreshToken string) (string, string, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	storedToken, err := s.repo.FindRefreshToken(oldRefreshToken)
	if err != nil {
		return "", "", errors.New("unrecognized or invalid refresh token")
	}

	if storedToken.IsRevoked {
		_ = s.repo.RevokeAllUserTokens(userID)
		return "", "", errors.New("revoked token reuse detected, all sessions have been revoked for security")
	}

	if err := s.repo.RevokeRefreshToken(oldRefreshToken); err != nil {
		return "", "", errors.New("failed to revoke old refresh token")
	}

	accessToken, err := utils.GenerateAccessJWT(user.ID, user.Role, s.cfg.JWT.Secret)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := utils.GenerateRefreshJWT(user.ID, s.cfg.JWT.Secret)
	if err != nil {
		return "", "", err
	}

	newTokenModel := RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.repo.CreateRefreshToken(&newTokenModel); err != nil {
		return "", "", errors.New("failed to save new refresh token")
	}

	return accessToken, newRefreshToken, nil
}
