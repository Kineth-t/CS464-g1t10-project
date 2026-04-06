package model

import "time"

type AuditLog struct {
    ID         int         `json:"id"`
    EntityType string      `json:"entity_type"`
    EntityID   int         `json:"entity_id"`
    Action     string      `json:"action"`
    OldValue   interface{} `json:"old_value"`
    NewValue   interface{} `json:"new_value"`
    ChangedBy  int         `json:"changed_by"`
    CreatedAt  time.Time   `json:"created_at"`
}