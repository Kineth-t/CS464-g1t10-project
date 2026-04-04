package router

import (
	"net/http"
	"os"
	"time"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/redis/go-redis/v9"
	"github.com/Kineth-t/CS464-g1t10-project/internal/handler"
	"github.com/Kineth-t/CS464-g1t10-project/internal/middleware"
)

// Setup initializes all routes and returns the main HTTP handler
func Setup(ph *handler.PhoneHandler, ah *handler.AuthHandler, ch *handler.CartHandler, pyh *handler.PaymentHandler, oh *handler.OrderHandler, rdb *redis.Client) http.Handler {
	mux := http.NewServeMux() // Main router

	// Read allowed origin from env (set to the frontend public URL in production)
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}
	cors := middleware.CORS(allowedOrigin)

// (Global) - Max 500 total requests/min for the whole app
globalSafety := middleware.SlidingWindowThrottle(rdb, 500, 60*time.Second, "api", true)

//(Per-User) Phone  - Max 50 requests/min per person
userPhoneThrottle := middleware.SlidingWindowThrottle(rdb, 50, 60*time.Second, "phones", false)

//(Per-User) Login  Max 5 attempts/min per person
loginThrottle := middleware.SlidingWindowThrottle(rdb, 5, 60*time.Second, "login", false)

// (Global) "Cart Queue" - only 100 people can add to cart per minute globally
globalCartLimit := middleware.SlidingWindowThrottle(rdb, 100, 60*time.Second, "cart_global", true)

// (Per-User) "Anti-Bot" - 2 adds per 10 seconds per person
userCartLimit := middleware.SlidingWindowThrottle(rdb, 2, 10*time.Second, "cart_user", false)

// (Global) Limit to 10 payments per minute globally to ensure database integrity
globalPaymentThrottle := middleware.SlidingWindowThrottle(rdb, 10, 60*time.Second, "payment_global", true)

	// ========================
	// Phone routes
	// ========================
	mux.HandleFunc("/phones", func(w http.ResponseWriter, r *http.Request) {
		// Ensure response is JSON
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			// Public: list all phones
			globalSafety(userPhoneThrottle(http.HandlerFunc(ph.ListPhones))).ServeHTTP(w, r)

		case http.MethodPost:
			// Admin only: create phone
			middleware.RequireAdmin(
				http.HandlerFunc(
					ph.CreatePhone,
				),
			).ServeHTTP(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/phones/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			// Public: get phone by ID
			userPhoneThrottle(http.HandlerFunc(ph.GetPhone)).ServeHTTP(w, r)

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
		loginThrottle(http.HandlerFunc(ah.Login)).ServeHTTP(w, r)
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
			// Auth -> Global Limit -> User Limit -> Handler
			middleware.RequireAuth(
				globalCartLimit(
					userCartLimit(
						http.HandlerFunc(ch.AddToCart),
					),
				),
			).ServeHTTP(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// cartMux.HandleFunc("/cart/checkout", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")

	// 	if r.Method != http.MethodPost {
	// 		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	// 		return
	// 	}

	// 	// Checkout cart
	// 	ch.Checkout(w, r)
	// })

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
	payHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Sequencer: Global Throttle -> Handler
		globalPaymentThrottle(http.HandlerFunc(pyh.Pay)).ServeHTTP(w, r)
	})
	mux.Handle("/pay", middleware.RequireAuth(payHandler))


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

	return cors(mux)
}