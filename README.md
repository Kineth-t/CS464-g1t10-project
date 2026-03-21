# Ringr Mobile

A full-stack mobile phone e-commerce application built for CS464. Customers can browse phones, manage a shopping cart, and checkout. Admins can manage the product catalog through a dedicated admin panel.

---

## Table of Contents

- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [API Reference](#api-reference)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Database Setup](#database-setup)
  - [Backend](#backend)
  - [Frontend](#frontend)
- [Environment Variables](#environment-variables)
- [Authentication](#authentication)
- [Features](#features)

---

## Tech Stack

| Layer      | Technology                                      |
|------------|-------------------------------------------------|
| Backend    | Go 1.25 · `net/http` · `golang-jwt/jwt v5`     |
| Database   | PostgreSQL · `pgx/v5` connection pool           |
| Frontend   | React 19 · Vite 8 · React Router 7             |
| Styling    | Tailwind CSS v4 · shadcn/ui (Base UI variant)   |
| Icons      | Lucide React                                    |
| Toasts     | Sonner                                          |

---

## Architecture

The backend follows a **layered architecture**:

```
HTTP Request
    ↓
Handler       (parses input, writes response)
    ↓
Service       (business logic, validation)
    ↓
Repository    (interface — abstracts DB)
    ↓
Postgres impl (pgx queries against PostgreSQL)
```

The frontend uses **React Context** for global auth state and a thin **API client** (`src/api/client.js`) that proxies all requests to the backend through Vite's dev-server proxy.

---

## Project Structure

```
CS464-g1t10-project/
├── backend/
│   ├── cmd/api/
│   │   └── main.go                  # Entry point; seeds admin on startup
│   ├── internal/
│   │   ├── handler/
│   │   │   ├── auth_handler.go      # Register / Login
│   │   │   ├── phone_handler.go     # Phone CRUD + purchase
│   │   │   └── cart_handler.go      # Cart get/add/remove/checkout
│   │   ├── middleware/
│   │   │   └── auth.go              # JWT RequireAuth / RequireAdmin
│   │   ├── model/
│   │   │   ├── user.go
│   │   │   ├── phone.go
│   │   │   └── cart.go
│   │   ├── repository/
│   │   │   ├── user_repo.go         # Repository interfaces
│   │   │   ├── phone_repo.go
│   │   │   ├── cart_repo.go
│   │   │   └── postgres/            # PostgreSQL implementations
│   │   │       ├── user_repository.go
│   │   │       ├── phone_repository.go
│   │   │       └── cart_repository.go
│   │   ├── router/
│   │   │   └── router.go            # Route registration
│   │   └── service/
│   │       ├── auth_service.go
│   │       ├── phone_service.go
│   │       └── cart_service.go
│   ├── migration/
│   │   └── 001_init.sql             # Schema (run once)
│   ├── go.mod
│   └── phones-api.postman_collection.json
│
└── frontend/
    ├── index.html
    ├── vite.config.js
    ├── src/
    │   ├── App.jsx                  # Router + top-level layout
    │   ├── main.jsx
    │   ├── api/
    │   │   └── client.js            # Typed API wrappers
    │   ├── components/
    │   │   ├── Navbar.jsx
    │   │   ├── ProtectedRoute.jsx   # Auth + admin route guards
    │   │   └── ui/                  # shadcn components
    │   ├── context/
    │   │   └── AuthContext.jsx      # JWT storage + auth state
    │   └── pages/
    │       ├── Home.jsx             # Phone listing + search
    │       ├── PhoneDetail.jsx      # Single phone + add to cart
    │       ├── Login.jsx
    │       ├── Register.jsx
    │       ├── Cart.jsx             # Cart management + checkout
    │       └── Admin.jsx            # Admin CRUD panel
    └── package.json
```

---

## Database Schema

Run `backend/migration/001_init.sql` against your PostgreSQL database to create all tables.

### `users`

| Column         | Type           | Notes                          |
|----------------|----------------|--------------------------------|
| `id`           | SERIAL PK      |                                |
| `username`     | VARCHAR(100)   | Unique, required               |
| `password`     | TEXT           | bcrypt-hashed                  |
| `phone_number` | VARCHAR(20)    |                                |
| `street`       | TEXT           | Address fields                 |
| `city`         | VARCHAR(100)   |                                |
| `state`        | VARCHAR(100)   |                                |
| `country`      | VARCHAR(100)   |                                |
| `zip_code`     | VARCHAR(20)    |                                |
| `role`         | VARCHAR(20)    | `'customer'` or `'admin'`      |

### `phones`

| Column        | Type           | Notes            |
|---------------|----------------|------------------|
| `id`          | SERIAL PK      |                  |
| `brand`       | VARCHAR(100)   | Required         |
| `model`       | VARCHAR(100)   | Required         |
| `price`       | NUMERIC(10,2)  | Required         |
| `stock`       | INT            | Default 0        |
| `description` | TEXT           |                  |

### `carts`

| Column    | Type         | Notes                            |
|-----------|--------------|----------------------------------|
| `id`      | SERIAL PK    |                                  |
| `user_id` | INT FK       | References `users(id)`           |
| `status`  | VARCHAR(20)  | `'active'` or `'checked_out'`    |

### `cart_items`

| Column     | Type          | Notes                                      |
|------------|---------------|--------------------------------------------|
| `id`       | SERIAL PK     |                                            |
| `cart_id`  | INT FK        | References `carts(id)` — cascades on delete|
| `phone_id` | INT FK        | References `phones(id)`                    |
| `quantity` | INT           | Required                                   |
| `price`    | NUMERIC(10,2) | Price at time of add                       |

`UNIQUE(cart_id, phone_id)` — adding the same phone again increments quantity instead of inserting a duplicate row.

---

## API Reference

Base URL: `http://localhost:8080`

All protected routes require the header:
```
Authorization: Bearer <jwt_token>
```

### Auth

| Method | Path             | Auth     | Description                  |
|--------|------------------|----------|------------------------------|
| POST   | `/auth/register` | None     | Create a new customer account|
| POST   | `/auth/login`    | None     | Login and receive a JWT token|

**POST `/auth/register`**
```json
// Request body
{
  "username": "john",
  "password": "secret123",
  "phone_number": "+1-555-0100",
  "address": {
    "street": "123 Main St",
    "city": "Springfield",
    "state": "IL",
    "country": "US",
    "zip_code": "62701"
  }
}
// Response 201
{ "id": 2, "username": "john", "role": "customer", ... }
```

**POST `/auth/login`**
```json
// Request body
{ "username": "john", "password": "secret123" }
// Response 200
{ "token": "<jwt>" }
```

---

### Phones

| Method | Path           | Auth         | Description              |
|--------|----------------|--------------|--------------------------|
| GET    | `/phones`      | None         | List all phones          |
| GET    | `/phones/{id}` | None         | Get a single phone       |
| POST   | `/phones`      | Admin only   | Create a phone listing   |
| PUT    | `/phones/{id}` | Admin only   | Update a phone listing   |
| DELETE | `/phones/{id}` | Admin only   | Delete a phone listing   |
| POST   | `/purchase`    | Authenticated| Buy a phone directly     |

**GET `/phones`**
```json
// Response 200
[
  { "id": 1, "brand": "Apple", "model": "iPhone 15", "price": 799.99, "stock": 10, "description": "..." }
]
```

**POST `/phones`** (Admin)
```json
// Request body
{ "brand": "Samsung", "model": "Galaxy S24", "price": 699.99, "stock": 5, "description": "..." }
// Response 201 — created phone object
```

**POST `/purchase`** (Authenticated)
```json
// Request body
{ "phone_id": 1, "quantity": 2 }
// Response 200
{ "message": "purchase successful" }
```

---

### Cart

All cart routes require authentication.

| Method | Path                | Auth          | Description            |
|--------|---------------------|---------------|------------------------|
| GET    | `/cart`             | Authenticated | Get current user's cart|
| POST   | `/cart`             | Authenticated | Add item to cart       |
| DELETE | `/cart/{itemId}`    | Authenticated | Remove item from cart  |
| POST   | `/cart/checkout`    | Authenticated | Checkout the cart      |

**GET `/cart`**
```json
// Response 200
{
  "id": 3,
  "user_id": 2,
  "status": "active",
  "items": [
    { "id": 7, "cart_id": 3, "phone_id": 1, "quantity": 2, "price": 799.99 }
  ]
}
```

**POST `/cart`**
```json
// Request body
{ "phone_id": 1, "quantity": 1 }
// Response 201 — cart item object
```

---

## Getting Started

### Prerequisites

- **Go** 1.21 or later
- **Node.js** 18 or later + npm
- **PostgreSQL** 14 or later

### Database Setup

1. Create a database:
   ```bash
   createdb ringr
   ```

2. Run the migration:
   ```bash
   psql -d ringr -f backend/migration/001_init.sql
   ```

### Backend

1. Set environment variables (see [Environment Variables](#environment-variables) below).

2. Run the server:
   ```bash
   cd backend
   go run ./cmd/api
   ```

   The server starts on `http://localhost:8080`. On first startup it automatically seeds an `admin` user with the password from `ADMIN_PASSWORD`.

### Frontend

1. Install dependencies:
   ```bash
   cd frontend
   npm install
   ```

2. Start the dev server:
   ```bash
   npm run dev
   ```

   Vite starts on `http://localhost:5173` and proxies all `/api/*` requests to the Go backend at `http://127.0.0.1:8080`.

3. Build for production:
   ```bash
   npm run build
   # Output → frontend/dist/
   ```

---

## Environment Variables

| Variable         | Required | Default                      | Description                                   |
|------------------|----------|------------------------------|-----------------------------------------------|
| `DATABASE_URL`   | Yes      | —                            | PostgreSQL DSN, e.g. `postgres://user:pass@localhost:5432/ringr` |
| `JWT_SECRET`     | No       | `change-me-in-production`    | Secret used to sign JWT tokens                |
| `ADMIN_PASSWORD` | No       | —                            | Password for the seeded `admin` account       |

Set them in your shell before running the backend:

```bash
export DATABASE_URL="postgres://postgres:password@localhost:5432/ringr"
export JWT_SECRET="your-strong-secret"
export ADMIN_PASSWORD="YourAdminPassword!"

cd backend && go run ./cmd/api
```

---

## Authentication

Ringr Mobile uses **JWT (JSON Web Tokens)** with HMAC-SHA256 signing.

**Token lifetime:** 24 hours from login.

**Token payload:**
```json
{
  "user_id": 1,
  "username": "john",
  "role": "customer",
  "exp": 1750000000
}
```

**Using the token:**
Include it in the `Authorization` header of any protected request:
```
Authorization: Bearer eyJhbGci...
```

The frontend stores the token in `localStorage` and automatically attaches it to every API request via the `client.js` wrapper.

---

## Features

- **Customer**
  - Register and log in
  - Browse the phone catalog with live search by brand or model
  - View phone details, stock availability, and description
  - Add phones to cart (re-adding the same phone stacks quantity)
  - Remove individual items from cart
  - Checkout cart
  - Direct purchase without using the cart

- **Admin**
  - All customer features
  - Create, edit, and delete phone listings from the Admin panel
  - Role assigned at account seed time (admin) or via direct database update

- **General**
  - JWT-based stateless authentication
  - Role-based access control enforced on both backend middleware and frontend route guards
  - Password hashing with bcrypt
  - Responsive UI with dark-mode CSS variables

---

## Postman Collection

A ready-to-import Postman collection is provided at:

```
backend/phones-api.postman_collection.json
```

Import it into Postman to explore and test all API endpoints interactively. Set the `base_url` collection variable to `http://localhost:8080` and the `token` variable to your JWT after logging in.
