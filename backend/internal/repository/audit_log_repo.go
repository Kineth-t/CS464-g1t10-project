package repository

import "github.com/Kineth-t/CS464-g1t10-project/internal/model"

// AuditLogRepo defines persistence operations for audit log entries.
type AuditLogRepo interface {
	Log(entry model.AuditLog) error
	GetRecent(limit, offset int) ([]model.AuditLog, error)
}
