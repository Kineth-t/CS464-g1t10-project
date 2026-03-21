# Ringr Mobile — Backend

Go REST API for the Ringr Mobile phone e-commerce application.

---

## Tech Stack

| Tool                   | Version | Purpose                          |
|------------------------|---------|----------------------------------|
| Go                     | 1.25    | Language                         |
| `net/http`             | stdlib  | HTTP server and routing          |
| `pgx/v5`               | 5.8     | PostgreSQL driver + pool         |
| `golang-jwt/jwt`       | v5      | JWT creation and validation      |
| `golang.org/x/crypto`  | —       | bcrypt password hashing          |

---

## Project Structure

```
backend/
├── cmd/api/
│   └── main.go                      # Entry point; wires dependencies, seeds admin
├── internal/
│   ├── handler/                     # HTTP layer — parse input, write response
│   │   ├── auth_handler.go          # POST /auth/register, POST /auth/login
│   │   ├── phone_handler.go         # GET/POST/PUT/DELETE /phones, POST /purchase
│   │   └── cart_handler.go          # GET/POST/DELETE /cart, POST /cart/checkout
│   ├── middleware/
│   │   └── auth.go                  # RequireAuth and RequireAdmin middleware
│   ├── model/                       # Plain data structs (no business logic)
│   │   ├── user.go
│   │   ├── phone.go
│   │   └── cart.go
│   ├── repository/                  # Data access interfaces
│   │   ├── user_repo.go
│   │   ├── phone_repo.go
│   │   ├── cart_repo.go
│   │   └── postgres/                # PostgreSQL implementations
│   │       ├── user_repository.go
│   │       ├── phone_repository.go
│   │       └── cart_repository.go
│   ├── router/
│   │   └── router.go                # Registers all routes with middleware
│   └── service/                     # Business logic layer
│       ├── auth_service.go
│       ├── phone_service.go
│       └── cart_service.go
├── migration/
│   └── 001_init.sql                 # Database schema (run once)
├── go.mod
├── go.sum
└── phones-api.postman_collection.json
```

---

## Architecture

```
HTTP Request
    ↓
Handler       — Parses input, writes JSON response
    ↓
Service       — Business logic, validation
    ↓
Repository    — Interface (decouples DB from logic)
    ↓
Postgres impl — pgx queries against PostgreSQL
```

Dependencies flow inward: handlers depend on services, services depend on repository interfaces. The concrete PostgreSQL implementations are injected at startup in `main.go`.

---

## Getting Started

### 1. Create the database

```bash
createdb ringr
psql -d ringr -f migration/001_init.sql
```

### 2. Set environment variables

| Variable         | Required | Default                       | Description                              |
|------------------|----------|-------------------------------|------------------------------------------|
| `DATABASE_URL`   | Yes      | —                             | PostgreSQL DSN                           |
| `JWT_SECRET`     | No       | `change-me-in-production`     | HMAC secret for signing tokens           |
| `ADMIN_PASSWORD` | No       | —                             | Password for the seeded `admin` account  |

```bash
export DATABASE_URL="postgres://postgres:password@localhost:5432/ringr"
export JWT_SECRET="your-strong-random-secret"
export ADMIN_PASSWORD="AdminPassword123!"
```

### 3. Run the server

```bash
go run ./cmd/api
```

Server starts on `http://localhost:8080`. On first startup, an `admin` user is automatically seeded using the value of `ADMIN_PASSWORD`.

---

## API Endpoints

### Auth

| Method | Path             | Auth     | Description            |
|--------|------------------|----------|------------------------|
| POST   | `/auth/register` | None     | Register a new user    |
| POST   | `/auth/login`    | None     | Login and receive JWT  |

### Phones

| Method | Path           | Auth          | Description              |
|--------|----------------|---------------|--------------------------|
| GET    | `/phones`      | None          | List all phones          |
| GET    | `/phones/{id}` | None          | Get a single phone       |
| POST   | `/phones`      | Admin only    | Create a phone listing   |
| PUT    | `/phones/{id}` | Admin only    | Update a phone listing   |
| DELETE | `/phones/{id}` | Admin only    | Delete a phone listing   |
| POST   | `/purchase`    | Authenticated | Buy a phone directly     |

### Cart

| Method | Path              | Auth          | Description              |
|--------|-------------------|---------------|--------------------------|
| GET    | `/cart`           | Authenticated | Get current user's cart  |
| POST   | `/cart`           | Authenticated | Add item to cart         |
| DELETE | `/cart/{itemId}`  | Authenticated | Remove item from cart    |
| POST   | `/cart/checkout`  | Authenticated | Checkout the cart        |

Protected routes require the header:
```
Authorization: Bearer <jwt_token>
```

---

## Authentication

- Tokens are signed with **HMAC-SHA256** using `JWT_SECRET`.
- Token lifetime is **24 hours**.
- Token payload includes `user_id`, `username`, `role`, and `exp`.
- `RequireAuth` middleware validates the token and injects `user_id` and `role` into the request context.
- `RequireAdmin` additionally checks that `role == "admin"`.

---

## Database Schema

See [`migration/001_init.sql`](migration/001_init.sql) for the full schema.

Key design decisions:
- `cart_items` has a `UNIQUE(cart_id, phone_id)` constraint — adding the same phone again uses `ON CONFLICT DO UPDATE SET quantity = quantity + EXCLUDED.quantity` to stack quantities instead of inserting a duplicate row.
- Passwords are never stored in plaintext; bcrypt is applied in the auth service before writing to the database.
- The `role` column on `users` defaults to `'customer'`; only the seeded admin account has `role = 'admin'`.

---

## Postman Collection

Import `phones-api.postman_collection.json` into Postman to test all endpoints interactively.

Set the `base_url` collection variable to `http://localhost:8080`, then after calling `/auth/login` copy the returned token into the `token` collection variable — all protected requests will pick it up automatically.
