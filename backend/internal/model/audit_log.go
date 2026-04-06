package model

import "time"

// AuditLog records a significant user or admin action for traceability.
type AuditLog struct {
	ID           int            `json:"id"`
	UserID       *int           `json:"user_id"`
	Action       string         `json:"action"`
	ResourceType string         `json:"resource_type,omitempty"`
	ResourceID   string         `json:"resource_id,omitempty"`
	Details      map[string]any `json:"details,omitempty"`
	IPAddress    string         `json:"ip_address,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}
