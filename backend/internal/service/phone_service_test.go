package service

import (
	"testing"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

func TestCreatePhone_ZeroPrice(t *testing.T) {
	svc := NewPhoneService(newMockPhoneRepo())
	_, err := svc.CreatePhone(model.Phone{Brand: "Apple", Model: "iPhone 15", Price: 0, Stock: 10})
	if err == nil {
		t.Fatal("expected error for zero price")
	}
}

func TestCreatePhone_NegativePrice(t *testing.T) {
	svc := NewPhoneService(newMockPhoneRepo())
	_, err := svc.CreatePhone(model.Phone{Brand: "Apple", Model: "iPhone 15", Price: -100, Stock: 10})
	if err == nil {
		t.Fatal("expected error for negative price")
	}
}

func TestCreatePhone_NegativeStock(t *testing.T) {
	svc := NewPhoneService(newMockPhoneRepo())
	_, err := svc.CreatePhone(model.Phone{Brand: "Apple", Model: "iPhone 15", Price: 999, Stock: -1})
	if err == nil {
		t.Fatal("expected error for negative stock")
	}
}

func TestCreatePhone_Success(t *testing.T) {
	svc := NewPhoneService(newMockPhoneRepo())
	phone, err := svc.CreatePhone(model.Phone{Brand: "Samsung", Model: "Galaxy S24", Price: 799.99, Stock: 50})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if phone.ID == 0 {
		t.Error("expected non-zero phone ID after creation")
	}
	if phone.Brand != "Samsung" {
		t.Errorf("expected brand Samsung, got %s", phone.Brand)
	}
}

func TestListPhones_Empty(t *testing.T) {
	svc := NewPhoneService(newMockPhoneRepo())
	phones := svc.ListPhones()
	if len(phones) != 0 {
		t.Errorf("expected 0 phones, got %d", len(phones))
	}
}

func TestListPhones_ReturnsAll(t *testing.T) {
	repo := newMockPhoneRepo()
	repo.Create(model.Phone{Brand: "Apple", Price: 999, Stock: 10})
	repo.Create(model.Phone{Brand: "Samsung", Price: 799, Stock: 5})
	svc := NewPhoneService(repo)
	phones := svc.ListPhones()
	if len(phones) != 2 {
		t.Errorf("expected 2 phones, got %d", len(phones))
	}
}

func TestGetPhone_NotFound(t *testing.T) {
	svc := NewPhoneService(newMockPhoneRepo())
	_, err := svc.GetPhone(999)
	if err == nil {
		t.Fatal("expected error for non-existent phone ID")
	}
}

func TestGetPhone_Found(t *testing.T) {
	repo := newMockPhoneRepo()
	created := repo.Create(model.Phone{Brand: "Apple", Model: "iPhone 15", Price: 999, Stock: 10})
	svc := NewPhoneService(repo)

	phone, err := svc.GetPhone(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if phone.Brand != "Apple" {
		t.Errorf("expected brand Apple, got %s", phone.Brand)
	}
}

func TestDeletePhone_NotFound(t *testing.T) {
	svc := NewPhoneService(newMockPhoneRepo())
	err := svc.DeletePhone(999)
	if err == nil {
		t.Fatal("expected error when deleting non-existent phone")
	}
}

func TestDeletePhone_Success(t *testing.T) {
	repo := newMockPhoneRepo()
	created := repo.Create(model.Phone{Brand: "Apple", Price: 999, Stock: 10})
	svc := NewPhoneService(repo)

	if err := svc.DeletePhone(created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := svc.GetPhone(created.ID); err == nil {
		t.Fatal("expected phone to be gone after deletion")
	}
}