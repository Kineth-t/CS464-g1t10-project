package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
)

type AuthService struct {
	repo repository.UserRepo
}

func NewAuthService(repo repository.UserRepo) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(username, password, phoneNumber string, address model.Address) (model.User, error) {
	if username == "" || password == "" {
		return model.User{}, errors.New("username and password are required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	return s.repo.Create(model.User{
		Username:    username,
		Password:    string(hash),
		PhoneNumber: phoneNumber,
		Address:     address,
		Role:        model.RoleCustomer,
	})
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(jwtSecret()))
}

func jwtSecret() string {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return "change-me-in-production"
	}
	return s
}