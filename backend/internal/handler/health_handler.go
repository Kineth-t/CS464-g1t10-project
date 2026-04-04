package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewHealthHandler(db *pgxpool.Pool, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, rdb: rdb}
}

type healthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Cache    string `json:"cache"`
}

// Health godoc
// @Summary     Health check
// @Description Returns the live status of the API, database, and cache
// @Tags        health
// @Produce     json
// @Success     200 {object} healthResponse
// @Failure     503 {object} healthResponse
// @Router      /health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp := healthResponse{
		Status:   "ok",
		Database: "ok",
		Cache:    "ok",
	}
	statusCode := http.StatusOK

	// Ping database
	if err := h.db.Ping(ctx); err != nil {
		resp.Database = "error"
		resp.Status = "degraded"
		statusCode = http.StatusServiceUnavailable
	}

	// Ping Redis (optional — report ok if not configured)
	if h.rdb != nil {
		if err := h.rdb.Ping(ctx).Err(); err != nil {
			resp.Cache = "error"
			resp.Status = "degraded"
			statusCode = http.StatusServiceUnavailable
		}
	} else {
		resp.Cache = "disabled"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
