package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

// newTestAuthHandler wires an AuthHandler backed by an in-memory user repo.
func newTestAuthHandler() *AuthHandler {
	svc := service.NewAuthService(newMockUserRepo())
	return NewAuthHandler(svc)
}

func TestRegister_Success(t *testing.T) {
	h := newTestAuthHandler()
	body, _ := json.Marshal(map[string]any{"username": "alice", "password": "secret123"})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRegister_MissingUsername(t *testing.T) {
	h := newTestAuthHandler()
	body, _ := json.Marshal(map[string]any{"username": "", "password": "secret123"})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty username, got %d", w.Code)
	}
}

func TestRegister_MissingPassword(t *testing.T) {
	h := newTestAuthHandler()
	body, _ := json.Marshal(map[string]any{"username": "alice", "password": ""})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty password, got %d", w.Code)
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	h := newTestAuthHandler()

	// First registration
	body, _ := json.Marshal(map[string]any{"username": "bob", "password": "pass123"})
	h.Register(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body)))

	// Duplicate registration
	body, _ = json.Marshal(map[string]any{"username": "bob", "password": "other123"})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for duplicate username, got %d", w.Code)
	}
}

func TestRegister_BadJSON(t *testing.T) {
	h := newTestAuthHandler()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte("{bad")))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for malformed JSON, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	h := newTestAuthHandler()

	// Register the user first
	regBody, _ := json.Marshal(map[string]any{"username": "carol", "password": "pass1234"})
	h.Register(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(regBody)))

	// Login
	body, _ := json.Marshal(map[string]string{"username": "carol", "password": "pass1234"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if resp["token"] == "" {
		t.Error("expected a non-empty JWT token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	h := newTestAuthHandler()

	regBody, _ := json.Marshal(map[string]any{"username": "dave", "password": "correct"})
	h.Register(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(regBody)))

	body, _ := json.Marshal(map[string]string{"username": "dave", "password": "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for wrong password, got %d", w.Code)
	}
}

func TestLogin_UnknownUser(t *testing.T) {
	h := newTestAuthHandler()
	body, _ := json.Marshal(map[string]string{"username": "nobody", "password": "pass"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for unknown user, got %d", w.Code)
	}
}

func TestLogin_BadJSON(t *testing.T) {
	h := newTestAuthHandler()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte("{invalid}")))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for malformed JSON, got %d", w.Code)
	}
}
