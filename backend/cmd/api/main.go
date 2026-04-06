// @title           Ringr Mobile API
// @version         1.0
// @description     REST API for the Ringr Mobile phone store.
//
// @contact.name    Ringr Mobile Team
//
// @host            localhost:8080
// @BasePath        /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token as: Bearer <token>

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/Kineth-t/CS464-g1t10-project/docs"
	"github.com/Kineth-t/CS464-g1t10-project/internal/handler"
	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
	pg "github.com/Kineth-t/CS464-g1t10-project/internal/repository/postgres"
	"github.com/Kineth-t/CS464-g1t10-project/internal/router"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	initLogger()

	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		slog.Error("database unreachable", "error", err)
        os.Exit(1)
	}
	slog.Info("connected to database")

	// Call your helper to get the Redis client
	rdb := initRedis()
	if rdb != nil {
		defer rdb.Close()
	}

	// Repos
	phoneRepo := pg.NewPhoneRepository(db)
	userRepo := pg.NewUserRepository(db)
	cartRepo := pg.NewCartRepository(db)
	orderRepo := pg.NewOrderRepository(db)
	auditRepo := pg.NewAuditRepository(db)

	// Cache
	phoneCache := repository.NewPhoneCache(rdb)

	// Services
	phoneSvc := service.NewPhoneService(phoneRepo, phoneCache, auditRepo)
	authSvc := service.NewAuthService(userRepo)
	cartSvc := service.NewCartService(cartRepo, phoneRepo)
	paymentSvc := service.NewPaymentService(cartRepo, phoneRepo, orderRepo)
	orderSvc := service.NewOrderService(orderRepo)

	// Seed admin if not exists
	seedAdmin(userRepo)

	// Handlers
	ph := handler.NewPhoneHandler(phoneSvc)
	ah := handler.NewAuthHandler(authSvc)
	ch := handler.NewCartHandler(cartSvc)
	pyh := handler.NewPaymentHandler(paymentSvc)
	oh := handler.NewOrderHandler(orderSvc)
	hh := handler.NewHealthHandler(db, rdb)
	uh := handler.NewUploadHandler()

	r := router.Setup(ph, ah, ch, pyh, oh, hh, uh, rdb)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("Ringr Mobile API starting", "port", port, "swagger", "http://localhost:"+port+"/swagger/index.html")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        slog.Error("server failed to start", "error", err)
        os.Exit(1)
    }
}

func initRedis() *redis.Client {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		slog.Warn("Redis not configured, caching disabled")
		return nil
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		slog.Error("Invalid Redis URL", "error", err)
		return nil
	}

	opts.PoolSize = 100
	opts.MinIdleConns = 10
	opts.DialTimeout = 5 * time.Second
	opts.ReadTimeout = 3 * time.Second

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Warn("Redis unreachable", "error", err)
		return nil
	}

	// log.Println("Connected to Redis")

	slog.Info("connected to redis", 
    "url", url, 
    "pool_size", opts.PoolSize,
)
	return rdb
}

func initLogger() {
    // If we are on Railway, use JSON. If local, use readable Text.
    var handler slog.Handler
    if os.Getenv("RAILWAY_ENVIRONMENT") != "" || os.Getenv("PORT") != "" {
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
    } else {
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
    }
    
    slog.SetDefault(slog.New(handler))
}

func seedAdmin(repo *pg.UserRepository) {
	_, err := repo.FindByUsername("admin")
	if err == nil {
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)
	_, err = repo.Create(model.User{
		Username:    "admin",
		Password:    string(hash),
		PhoneNumber: "",
		Role:        model.RoleAdmin,
	})
	if err != nil {
		slog.Error("failed to seed admin", "error", err)
		return
	}
	slog.Info("Admin user seeded")
}
