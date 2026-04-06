package handler

import (
	"encoding/json"
	"net/http"
	"log/slog"
	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

// AuthHandler handles authentication-related HTTP requests (register, login)
type AuthHandler struct {
	service *service.AuthService // Reference to business logic layer
}

// Constructor function to create a new AuthHandler
func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Register handles user registration requests
//
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body object{username=string,password=string,phone_number=string,address=object} true "Registration payload"
// @Success      201  {object}  model.User
// @Failure      400  {string}  string "invalid body / username already taken"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Define a temporary struct to store incoming JSON request body
	var body struct {
		Username    string        `json:"username"`
		Password    string        `json:"password"`
		PhoneNumber string        `json:"phone_number"`
		Address     model.Address `json:"address"`
	}

	// Decode JSON request body into the struct
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		// If JSON is invalid, return 400 Bad Request
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// Call service layer to register the user
	user, err := h.service.Register(
		body.Username,
		body.Password,
		body.PhoneNumber,
		body.Address,
	)
	if err != nil {
		// If registration fails (e.g., validation error), return 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If successful, return HTTP 201 Created
	w.WriteHeader(http.StatusCreated)

	slog.Info("new user registered", 
	"username", body.Username, 
	"phone", body.PhoneNumber,
	)

	// Encode the created user as JSON in the response
	json.NewEncoder(w).Encode(user)
}

// Login handles user login requests
//
// @Summary      Login and receive a JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body object{username=string,password=string} true "Login credentials"
// @Success      200  {object}  object{token=string}
// @Failure      400  {string}  string "invalid body"
// @Failure      401  {string}  string "invalid credentials"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var body struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "invalid body", http.StatusBadRequest)
        return
    }

    token, err := h.service.Login(body.Username, body.Password)
    if err != nil {
        // 1. Get the real IP address
        userIP := r.Header.Get("X-Forwarded-For")
        if userIP == "" {
            userIP = r.RemoteAddr
        }

        // 2. Log the failure (Fixed: use 'body.Username' and 'userIP')
        slog.Warn("login failed",
            "username",  body.Username, // Changed from req.Username
            "ip_address", userIP,        // Used the variable we defined above
            "error",      err.Error(),
        )

        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    // Optional: Log successful logins too!
    slog.Info("login successful", "username", body.Username)

    json.NewEncoder(w).Encode(map[string]string{
        "token": token,
    })
}