package router

import (
	"net/http"

	"github.com/Kineth-t/CS464-g1t10-project/internal/handler"
	"github.com/Kineth-t/CS464-g1t10-project/internal/middleware"
)

func Setup(ph *handler.PhoneHandler, ah *handler.AuthHandler, ch *handler.CartHandler) http.Handler {
	mux := http.NewServeMux()

	// Phone routes
	mux.HandleFunc("/phones", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			ph.ListPhones(w, r)
		case http.MethodPost:
			middleware.RequireAdmin(http.HandlerFunc(ph.CreatePhone)).ServeHTTP(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/phones/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			ph.GetPhone(w, r)
		case http.MethodPut:
			middleware.RequireAdmin(http.HandlerFunc(ph.UpdatePhone)).ServeHTTP(w, r)
		case http.MethodDelete:
			middleware.RequireAdmin(http.HandlerFunc(ph.DeletePhone)).ServeHTTP(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/purchase", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ph.PurchasePhone(w, r)
	})

	// Auth routes (public)
	mux.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ah.Register(w, r)
	})
	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ah.Login(w, r)
	})

	// Cart routes (protected)
	cartMux := http.NewServeMux()
	cartMux.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			ch.GetCart(w, r)
		case http.MethodPost:
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
		ch.Checkout(w, r)
	})
	cartMux.HandleFunc("/cart/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ch.RemoveFromCart(w, r)
	})

	mux.Handle("/cart", middleware.RequireAuth(cartMux))
	mux.Handle("/cart/", middleware.RequireAuth(cartMux))

	return mux
}