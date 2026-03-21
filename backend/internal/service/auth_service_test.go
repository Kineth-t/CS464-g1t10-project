package service

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

func TestRegister_EmptyUsername(t *testing.T) {
	svc := NewAuthService(newMockUserRepo())
	_, err := svc.Register("", "password123", "", model.Address{})
	if err == nil {
		t.Fatal("expected error for empty username")
	}
}

func TestRegister_EmptyPassword(t *testing.T) {
	svc := NewAuthService(newMockUserRepo())
	_, err := svc.Register("alice", "", "", model.Address{})
	if err == nil {
		t.Fatal("expected error for empty password")
	}
}

func TestRegister_Success(t *testing.T) {
	svc := NewAuthService(newMockUserRepo())
	user, err := svc.Register("alice", "password123", "555-1234", model.Address{City: "Manila"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID == 0 {
		t.Error("expected non-zero user ID")
	}
	if user.Username != "alice" {
		t.Errorf("expected username alice, got %s", user.Username)
	}
	if user.Role != model.RoleCustomer {
		t.Errorf("expected role customer, got %s", user.Role)
	}
	if user.Password == "password123" {
		t.Error("password should be hashed, not stored as plain text")
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	svc := NewAuthService(newMockUserRepo())
	if _, err := svc.Register("alice", "password123", "", model.Address{}); err != nil {
		t.Fatalf("unexpected error on first register: %v", err)
	}
	_, err := svc.Register("alice", "otherpass", "", model.Address{})
	if err == nil {
		t.Fatal("expected error for duplicate username")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	svc := NewAuthService(newMockUserRepo())
	_, err := svc.Login("nobody", "password")
	if err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("expected 'invalid credentials', got %v", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := newMockUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	repo.users["alice"] = model.User{ID: 1, Username: "alice", Password: string(hash), Role: model.RoleCustomer}

	svc := NewAuthService(repo)
	_, err := svc.Login("alice", "wrong")
	if err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("expected 'invalid credentials', got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	repo := newMockUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	repo.users["alice"] = model.User{ID: 1, Username: "alice", Password: string(hash), Role: model.RoleCustomer}

	svc := NewAuthService(repo)
	token, err := svc.Login("alice", "correct")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty JWT token")
	}
}
