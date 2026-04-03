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
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"github.com/redis/go-redis/v9"
	_ "github.com/Kineth-t/CS464-g1t10-project/docs"
	"github.com/Kineth-t/CS464-g1t10-project/internal/handler"
	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
	pg "github.com/Kineth-t/CS464-g1t10-project/internal/repository/postgres"
	"github.com/Kineth-t/CS464-g1t10-project/internal/router"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

func main() {
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("database unreachable:", err)
	}
	log.Println("Connected to database")

	//Call your helper to get the Redis client
    rdb := initRedis()
    defer rdb.Close()

	// Repos
	phoneRepo := pg.NewPhoneRepository(db)
	userRepo  := pg.NewUserRepository(db)
	cartRepo  := pg.NewCartRepository(db)
	orderRepo := pg.NewOrderRepository(db)

	// Cache 
	phoneCache := repository.NewPhoneCache(rdb)

	// Services
	phoneSvc := service.NewPhoneService(phoneRepo, phoneCache)
	authSvc  := service.NewAuthService(userRepo)
	cartSvc  := service.NewCartService(cartRepo, phoneRepo)
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

	r := router.Setup(ph, ah, ch, pyh, oh, rdb)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,                // Your router
		ReadTimeout:  5 * time.Second,  // Max time to read the request
		WriteTimeout: 10 * time.Second, // Max time to write the response
		IdleTimeout:  120 * time.Second, // Max time to keep idle connections alive
	}

	log.Println("----------------------------------------")
	log.Println("  Ringr Mobile API")
	log.Println("  http://localhost:8080")
	log.Println("  Swagger UI: http://localhost:8080/swagger/index.html")
	log.Println("----------------------------------------")
	log.Fatal(server.ListenAndServe())
}

func initRedis() *redis.Client {
    url := os.Getenv("REDIS_URL")
    if url == "" {
        // Use the Docker service name 'redis' if REDIS_URL isn't set
        url = "redis://redis:6379" 
    }

    opts, err := redis.ParseURL(url)
    if err != nil {
        log.Println("Invalid Redis URL: %v", err)
		opts = &redis.Options{Addr: "redis:6379"}
    }

	// Your optimized pool settings
	opts.PoolSize = 100
	opts.MinIdleConns = 10

	// Connection timeouts are vital for cloud stability
	opts.DialTimeout = 5 * time.Second
    opts.ReadTimeout = 3 * time.Second

    rdb := redis.NewClient(opts)

	// This allows Railway to finish starting the app while Redis warms up.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Check if Redis is alive [cite: 388]
    if err := rdb.Ping(context.Background()).Err(); err != nil {
        log.Printnln("Redis unreachable: %v", err)
    }
    
    log.Println("Connected to Redis")
    return rdb
}

func seedAdmin(repo *pg.UserRepository) {
	_, err := repo.FindByUsername("admin")
	if err == nil {
		return // already exists
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)
	_, err = repo.Create(model.User{
		Username:    "admin",
		Password:    string(hash),
		PhoneNumber: "",
		Role:        model.RoleAdmin,
	})
	if err != nil {
		log.Println("failed to seed admin:", err)
		return
	}
	log.Println("Admin user seeded")
}