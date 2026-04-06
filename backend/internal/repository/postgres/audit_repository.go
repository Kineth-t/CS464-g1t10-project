package postgres

import (
	"context"
	"encoding/json"
	
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

type AuditRepository struct {
	db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(log model.AuditLog) error {
	// Convert interfaces to JSON for the JSONB columns
	oldVal, _ := json.Marshal(log.OldValue)
	newVal, _ := json.Marshal(log.NewValue)

	query := `
		INSERT INTO audit_logs (entity_type, entity_id, action, old_value, new_value, changed_by)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(context.Background(), query, 
		log.EntityType, log.EntityID, log.Action, oldVal, newVal, log.ChangedBy)
	return err
}