package user

import "gorm.io/gorm"

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (User, error)
	FindByID(id uint) (User, error)
	CreateRefreshToken(token *RefreshToken) error
	FindRefreshToken(tokenString string) (*RefreshToken, error)
	RevokeRefreshToken(tokenString string) error
	RevokeAllUserTokens(userID uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) FindByID(id uint) (User, error) {
	var user User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *userRepository) CreateRefreshToken(token *RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *userRepository) FindRefreshToken(tokenString string) (*RefreshToken, error) {
	var token RefreshToken
	err := r.db.Where("token = ?", tokenString).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *userRepository) RevokeRefreshToken(tokenString string) error {
	return r.db.Model(&RefreshToken{}).Where("token = ?", tokenString).Update("is_revoked", true).Error
}

func (r *userRepository) RevokeAllUserTokens(userID uint) error {
	return r.db.Model(&RefreshToken{}).Where("user_id = ? AND is_revoked = ?", userID, false).Update("is_revoked", true).Error
}

