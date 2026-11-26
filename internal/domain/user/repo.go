package user

import (
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id uint) (User, error)
	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	FindRefreshToken(ctx context.Context, tokenString string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenString string) error
	RevokeAllUserTokens(ctx context.Context, userID uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return user, err
}

func (r *userRepository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *userRepository) FindRefreshToken(ctx context.Context, tokenString string) (*RefreshToken, error) {
	var token RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", tokenString).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *userRepository) RevokeRefreshToken(ctx context.Context, tokenString string) error {
	return r.db.WithContext(ctx).Model(&RefreshToken{}).Where("token = ?", tokenString).Update("is_revoked", true).Error
}

func (r *userRepository) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&RefreshToken{}).Where("user_id = ? AND is_revoked = ?", userID, false).Update("is_revoked", true).Error
}

