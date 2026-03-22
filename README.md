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
  - [1. Configure environment variables](#1-configure-environment-variables)
  - [2. Set up the database](#2-set-up-the-database)
  - [3. Install dependencies](#3-install-dependencies)
  - [4. Start both servers](#4-start-both-servers)
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
│   │   │   ├── phone_handler.go     # Phone CRUD
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

**Interactive docs:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) — available whenever the backend is running. Click **Authorize** and paste `Bearer <your_jwt>` to test protected endpoints directly from the browser.

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

| Method | Path           | Auth       | Description              |
|--------|----------------|------------|--------------------------|
| GET    | `/phones`      | None       | List all phones          |
| GET    | `/phones/{id}` | None       | Get a single phone       |
| POST   | `/phones`      | Admin only | Create a phone listing   |
| PUT    | `/phones/{id}` | Admin only | Update a phone listing   |
| DELETE | `/phones/{id}` | Admin only | Delete a phone listing   |

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

- **Go** 1.21 or later — [go.dev/dl](https://go.dev/dl)
- **Node.js** 18 or later + npm — [nodejs.org](https://nodejs.org)
- **PostgreSQL** 14 or later — running and accessible
- **make** — pre-installed on Mac/Linux; on Windows use [Git Bash](https://git-scm.com) or [WSL](https://learn.microsoft.com/en-us/windows/wsl/)

### 1. Configure environment variables

```bash
cp .env.example .env
```

Open `.env` and set your values. The defaults in `.env.example` work if your local Postgres matches those credentials. See [Environment Variables](#environment-variables) for details.

### 2. Set up the database

Create a database and run the migration:

```bash
createdb phonestore
psql -d phonestore -f backend/migration/001_init.sql
```

> If your Postgres is in Docker or has a different user/host, adjust the `createdb`/`psql` commands accordingly and update `DATABASE_URL` in `.env`.

### 3. Install dependencies

```bash
make install
```

### 4. Start both servers

```bash
make dev
```

This starts the Go backend and Vite frontend in parallel:

| Service  | URL                                              |
|----------|--------------------------------------------------|
| Frontend | http://localhost:5173                            |
| Backend  | http://localhost:8080                            |
| Swagger  | http://localhost:8080/swagger/index.html         |

On first startup the backend automatically seeds an `admin` account using the `ADMIN_PASSWORD` from your `.env`.

### Running servers individually

```bash
make backend    # Go API only
make frontend   # Vite only
```

### Production build

```bash
cd frontend && npm run build
# Output → frontend/dist/
```

### Without make (Windows)

If `make` is unavailable, run each command manually in separate terminals:

```bash
# Terminal 1 – backend
cd backend
set DATABASE_URL=postgres://root:mysecretpassword@localhost:5432/phonestore
set JWT_SECRET=phonestore-jwt-secret-cs464
set ADMIN_PASSWORD=Admin1234!
go run ./cmd/api

# Terminal 2 – frontend
cd frontend
npm run dev
```

---

## Environment Variables

Copy `.env.example` to `.env` and fill in your values. The `make dev` / `make backend` targets load `.env` automatically.

| Variable         | Required | Description                                                        |
|------------------|----------|--------------------------------------------------------------------|
| `DATABASE_URL`   | Yes      | PostgreSQL DSN, e.g. `postgres://user:pass@localhost:5432/phonestore` |
| `JWT_SECRET`     | Yes      | Secret used to sign JWT tokens — change this in production         |
| `ADMIN_PASSWORD` | Yes      | Password for the seeded `admin` account                            |

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
  - Checkout cart (deducts stock on completion)

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
