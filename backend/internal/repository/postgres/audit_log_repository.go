package postgres

import (
	"context"
	"encoding/json"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuditLogRepository persists audit events to the audit_logs table.
type AuditLogRepository struct {
	db *pgxpool.Pool
}

// NewAuditLogRepository creates a new AuditLogRepository.
func NewAuditLogRepository(db *pgxpool.Pool) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

// Log inserts a single audit event.
func (r *AuditLogRepository) Log(entry model.AuditLog) error {
	detailsJSON, _ := json.Marshal(entry.Details)
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO audit_logs (user_id, action, resource_type, resource_id, details, ip_address)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		entry.UserID, entry.Action, entry.ResourceType, entry.ResourceID, detailsJSON, entry.IPAddress,
	)
	return err
}

// GetRecent returns audit log entries ordered newest-first.
func (r *AuditLogRepository) GetRecent(limit, offset int) ([]model.AuditLog, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, user_id, action, resource_type, resource_id, details, ip_address, created_at
		 FROM audit_logs
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.AuditLog
	for rows.Next() {
		var l model.AuditLog
		var resourceType, resourceID, ipAddress *string
		var detailsJSON []byte
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &resourceType, &resourceID, &detailsJSON, &ipAddress, &l.CreatedAt); err != nil {
			return nil, err
		}
		if resourceType != nil {
			l.ResourceType = *resourceType
		}
		if resourceID != nil {
			l.ResourceID = *resourceID
		}
		if ipAddress != nil {
			l.IPAddress = *ipAddress
		}
		if len(detailsJSON) > 0 {
			_ = json.Unmarshal(detailsJSON, &l.Details)
		}
		logs = append(logs, l)
	}
	return logs, nil
}
