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
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("database unreachable:", err)
	}
	log.Println("Connected to database")

	// Call your helper to get the Redis client
	rdb := initRedis()
	if rdb != nil {
		defer rdb.Close()
	}

	// Ensure audit_logs table exists (idempotent)
	ensureAuditTable(db)
	// Allow phones to be deleted even when they appear in order history
	ensureOrderItemsNullablePhone(db)

	// Repos
	phoneRepo := pg.NewPhoneRepository(db)
	userRepo := pg.NewUserRepository(db)
	cartRepo := pg.NewCartRepository(db)
	orderRepo := pg.NewOrderRepository(db)
	auditRepo := pg.NewAuditLogRepository(db)

	// Cache
	phoneCache := repository.NewPhoneCache(rdb)

	// Services
	phoneSvc := service.NewPhoneService(phoneRepo, phoneCache)
	authSvc := service.NewAuthService(userRepo)
	cartSvc := service.NewCartService(cartRepo, phoneRepo)
	paymentSvc := service.NewPaymentService(cartRepo, phoneRepo, orderRepo, phoneCache)
	orderSvc := service.NewOrderService(orderRepo)
	auditSvc := service.NewAuditService(auditRepo)

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
	alh := handler.NewAuditHandler(auditSvc)

	// Attach audit service to handlers that log events
	ph.SetAudit(auditSvc)
	ah.SetAudit(auditSvc)
	ch.SetAudit(auditSvc)
	pyh.SetAudit(auditSvc)

	r := router.Setup(ph, ah, ch, pyh, oh, hh, uh, alh, rdb)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
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
		log.Println("Redis not configured, caching disabled")
		return nil
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Printf("Invalid Redis URL: %v", err)
		log.Println("Redis disabled")
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
		log.Printf("Redis unreachable: %v", err)
		log.Println("Redis disabled")
		return nil
	}

	log.Println("Connected to Redis")
	return rdb
}

// ensureOrderItemsNullablePhone makes order_items.phone_id nullable with
// ON DELETE SET NULL so that deleting a phone doesn't destroy order history.
// The phone_name column already snapshots the product name at order time.
func ensureOrderItemsNullablePhone(db *pgxpool.Pool) {
	ctx := context.Background()
	// Drop NOT NULL if still present
	db.Exec(ctx, `ALTER TABLE order_items ALTER COLUMN phone_id DROP NOT NULL`)
	// Re-create FK with SET NULL (drop first so it's idempotent across restarts)
	db.Exec(ctx, `ALTER TABLE order_items DROP CONSTRAINT IF EXISTS order_items_phone_id_fkey`)
	_, err := db.Exec(ctx, `
		ALTER TABLE order_items
		ADD CONSTRAINT order_items_phone_id_fkey
		FOREIGN KEY (phone_id) REFERENCES phones(id) ON DELETE SET NULL
	`)
	if err != nil {
		log.Printf("Warning: could not update order_items FK: %v", err)
	} else {
		log.Println("order_items phone FK updated (ON DELETE SET NULL)")
	}
}

func ensureAuditTable(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS audit_logs (
			id            SERIAL PRIMARY KEY,
			user_id       INT REFERENCES users(id) ON DELETE SET NULL,
			action        VARCHAR(64) NOT NULL,
			resource_type VARCHAR(32),
			resource_id   VARCHAR(64),
			details       JSONB,
			ip_address    VARCHAR(45),
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs (created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id    ON audit_logs (user_id);
	`)
	if err != nil {
		log.Printf("Warning: could not ensure audit_logs table: %v", err)
	} else {
		log.Println("audit_logs table ready")
	}
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
		log.Println("failed to seed admin:", err)
		return
	}
	log.Println("Admin user seeded")
}
