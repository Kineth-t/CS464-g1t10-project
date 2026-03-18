package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

type PhoneHandler struct {
	service *service.PhoneService
}

func NewPhoneHandler(s *service.PhoneService) *PhoneHandler {
	return &PhoneHandler{service: s}
}

func (h *PhoneHandler) ListPhones(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.service.ListPhones())
}

func (h *PhoneHandler) GetPhone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/phones/"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	phone, err := h.service.GetPhone(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(phone)
}

func (h *PhoneHandler) CreatePhone(w http.ResponseWriter, r *http.Request) {
	var p model.Phone
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	created, err := h.service.CreatePhone(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *PhoneHandler) UpdatePhone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/phones/"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var p model.Phone
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	p.ID = id
	if err := h.service.UpdatePhone(p); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func (h *PhoneHandler) DeletePhone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/phones/"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.service.DeletePhone(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PhoneHandler) PurchasePhone(w http.ResponseWriter, r *http.Request) {
	var body struct {
		PhoneID  int `json:"phone_id"`
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := h.service.PurchasePhone(body.PhoneID, body.Quantity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "purchase successful"})
}