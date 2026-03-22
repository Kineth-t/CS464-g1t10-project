# Ringr Mobile – developer Makefile
#
# Prerequisites:
#   - Go 1.21+         (https://go.dev/dl)
#   - Node.js 18+      (https://nodejs.org)
#   - PostgreSQL 15+   running and accessible via DATABASE_URL
#
# Quick start:
#   make setup   – copy .env.example → .env (edit it first if needed)
#   make install – install all dependencies
#   make dev     – start backend + frontend together

-include .env
export

.PHONY: help setup install dev backend frontend

## Show this help message
help:
	@echo ""
	@echo "  Ringr Mobile – available targets"
	@echo ""
	@echo "  make setup     Copy .env.example to .env"
	@echo "  make install   Install backend and frontend dependencies"
	@echo "  make dev       Start backend + frontend (runs both together)"
	@echo "  make backend   Start only the Go backend   (http://localhost:8080)"
	@echo "  make frontend  Start only the Vite frontend (http://localhost:5173)"
	@echo ""

## Copy .env.example to .env (skips if .env already exists)
setup:
	@if [ -f .env ]; then \
		echo ".env already exists – skipping. Edit it manually if needed."; \
	else \
		cp .env.example .env; \
		echo "Created .env from .env.example – open it and set your values."; \
	fi

## Install all dependencies
install:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "Tidying backend modules..."
	cd backend && go mod tidy
	@echo "Done."

## Start backend + frontend in parallel
dev:
	@echo "Starting Ringr Mobile (backend + frontend)..."
	@$(MAKE) -j2 backend frontend

## Start only the Go backend
backend:
	@echo "Backend → http://localhost:8080  |  Swagger → http://localhost:8080/swagger/index.html"
	cd backend && go run ./cmd/api

## Start only the Vite frontend
frontend:
	@echo "Frontend → http://localhost:5173"
	cd frontend && npm run dev
