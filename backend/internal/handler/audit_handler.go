package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

// AuditHandler exposes the audit log to admins.
type AuditHandler struct {
	auditSvc *service.AuditService
}

// NewAuditHandler creates a new AuditHandler.
func NewAuditHandler(s *service.AuditService) *AuditHandler {
	return &AuditHandler{auditSvc: s}
}

// List returns recent audit log entries.
//
// @Summary      List audit log entries (admin only)
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Param        limit   query  int  false  "Max entries to return (default 100)"
// @Param        offset  query  int  false  "Pagination offset"
// @Success      200  {array}   model.AuditLog
// @Router       /audit-logs [get]
func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 100
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			offset = n
		}
	}

	logs, err := h.auditSvc.GetRecent(limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch audit logs"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
