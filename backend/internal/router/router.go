package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/Kineth-t/CS464-g1t10-project/internal/handler"
	"github.com/Kineth-t/CS464-g1t10-project/internal/middleware"
)

// Setup initializes all routes and returns the main HTTP handler
func Setup(ph *handler.PhoneHandler, ah *handler.AuthHandler, ch *handler.CartHandler, pyh *handler.PaymentHandler, oh *handler.OrderHandler) http.Handler {
	mux := http.NewServeMux() // Main router

	// ========================
	// Phone routes
	// ========================
	mux.HandleFunc("/phones", func(w http.ResponseWriter, r *http.Request) {
		// Ensure response is JSON
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			// Public: list all phones
			ph.ListPhones(w, r)

		case http.MethodPost:
			// Admin only: create phone
			middleware.RequireAdmin(http.HandlerFunc(ph.CreatePhone)).ServeHTTP(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/phones/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			// Public: get phone by ID
			ph.GetPhone(w, r)

		case http.MethodPut:
			// Admin only: update phone
			middleware.RequireAdmin(http.HandlerFunc(ph.UpdatePhone)).ServeHTTP(w, r)

		case http.MethodDelete:
			// Admin only: delete phone
			middleware.RequireAdmin(http.HandlerFunc(ph.DeletePhone)).ServeHTTP(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// ========================
	// Auth routes (PUBLIC)
	// ========================
	mux.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Register new user
		ah.Register(w, r)
	})

	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Login user → returns JWT
		ah.Login(w, r)
	})

	// ========================
	// Cart routes (PROTECTED)
	// ========================
	cartMux := http.NewServeMux()

	cartMux.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			// Get current user's cart
			ch.GetCart(w, r)

		case http.MethodPost:
			// Add item to cart
			ch.AddToCart(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	cartMux.HandleFunc("/cart/checkout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Checkout cart
		ch.Checkout(w, r)
	})

	cartMux.HandleFunc("/cart/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Remove item from cart
		ch.RemoveFromCart(w, r)
	})

	// Apply authentication middleware to ALL /cart routes
	mux.Handle("/cart", middleware.RequireAuth(cartMux))
	mux.Handle("/cart/", middleware.RequireAuth(cartMux))

	// ========================
	// Payment routes (PROTECTED)
	// ========================
	mux.Handle("/pay", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Process Stripe payment for the user's active cart
		pyh.Pay(w, r)
	})))


	// ========================
	// Order routes (PROTECTED)
	// ========================
	// Order routes (protected) — read-only, no mutations
	mux.Handle("/orders", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// List all orders for the logged in user
		oh.GetOrders(w, r)
	})))

	mux.Handle("/orders/", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Get a single order by Stripe payment intent ID
		oh.GetOrder(w, r)
	})))

	// Swagger UI
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}