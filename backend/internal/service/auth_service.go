package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"     // JWT library
	"golang.org/x/crypto/bcrypt"       // Password hashing

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
)

// AuthService contains business logic for authentication
type AuthService struct {
	repo repository.UserRepo // Interface to interact with user data
}

// Constructor
func NewAuthService(repo repository.UserRepo) *AuthService {
	return &AuthService{repo: repo}
}

// Register creates a new user account
func (s *AuthService) Register(username, password, phoneNumber string, address model.Address) (model.User, error) {

	// Basic validation
	if username == "" || password == "" {
		return model.User{}, errors.New("username and password are required")
	}

	// Hash the password before storing (VERY IMPORTANT for security)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	// Create user in database
	return s.repo.Create(model.User{
		Username:    username,
		Password:    string(hash), // Store hashed password (not plain text)
		PhoneNumber: phoneNumber,
		Address:     address,
		Role:        model.RoleCustomer, // Default role
	})
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(username, password string) (string, error) {

	// Find user by username
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials") // Do not reveal user existence
	}

	// Compare stored hash with provided password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Create JWT token with claims (payload)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24h
	})

	// Sign token using secret key
	return token.SignedString([]byte(jwtSecret()))
}

// jwtSecret retrieves secret key from environment variable
func jwtSecret() string {
	s := os.Getenv("JWT_SECRET")

	// Fallback (NOT safe for production)
	if s == "" {
		return "change-me-in-production"
	}

	return s
}