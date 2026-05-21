package user

import (
	"errors"

	"github.com/nullablenone/go-news-api/config"
	"github.com/nullablenone/go-news-api/pkg/utils/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req RegisterRequest) (User, error)
	Login(req LoginRequest) (string, error)
}

type userService struct {
	repo UserRepository
	env  *config.Env
}

func NewUserService(repo UserRepository, env *config.Env) UserService {
	return &userService{
		repo: repo,
		env:  env,
	}
}

func (s *userService) Register(req RegisterRequest) (User, error) {
	// 1. Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	// Field Role langsung dikunci ke "user" demi keamanan
	user := User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user", 
	}

	// 2. Simpan ke database
	err = s.repo.Create(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *userService) Login(req LoginRequest) (string, error) {
	// 1. Cari user berdasarkan email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// 2. Bandingkan password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// 3. Generate Token JWT jika password cocok
	token, err := jwt.GenerateJWT(user.ID, user.Role, s.env.JWTSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}