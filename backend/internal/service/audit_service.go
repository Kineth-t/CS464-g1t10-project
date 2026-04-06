package service

import (
	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
)

// AuditService records significant events for traceability.
type AuditService struct {
	repo repository.AuditLogRepo
}

// NewAuditService creates a new AuditService.
func NewAuditService(repo repository.AuditLogRepo) *AuditService {
	return &AuditService{repo: repo}
}

// Log records an audit event asynchronously so it never slows down a request.
func (s *AuditService) Log(userID *int, action, resourceType, resourceID, ip string, details map[string]any) {
	go func() {
		_ = s.repo.Log(model.AuditLog{
			UserID:       userID,
			Action:       action,
			ResourceType: resourceType,
			ResourceID:   resourceID,
			Details:      details,
			IPAddress:    ip,
		})
	}()
}

// GetRecent returns the most recent audit log entries.
func (s *AuditService) GetRecent(limit, offset int) ([]model.AuditLog, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	entries, err := s.repo.GetRecent(limit, offset)
	if entries == nil {
		entries = []model.AuditLog{}
	}
	return entries, err
}
