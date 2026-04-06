# Ringr Mobile

A full-stack mobile phone e-commerce application built for CS464. Customers can browse phones, manage a shopping cart, and checkout with Stripe. Admins can manage the product catalog, upload product images, and monitor system activity via a live audit log.

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
  - [2. Start with Docker](#2-start-with-docker)
  - [3. Running without Docker](#3-running-without-docker)
- [Environment Variables](#environment-variables)
- [Authentication](#authentication)
- [Features](#features)

---

## Tech Stack

| Layer            | Technology                                      |
| ---------------- | ----------------------------------------------- |
| Backend          | Go 1.25 · `net/http` · `golang-jwt/jwt v5`    |
| Database         | PostgreSQL 16 · `pgx/v5` connection pool       |
| Cache            | Redis · sliding-window rate limiting + phone cache |
| Payments         | Stripe (`stripe-go/v76`)                       |
| Image Storage    | Cloudinary (`cloudinary-go/v2`)                |
| Frontend         | React 19 · Vite 8 · React Router 7            |
| Styling          | Tailwind CSS v4 · shadcn/ui (Base UI variant)  |
| Icons            | Lucide React                                    |
| Toasts           | Sonner                                          |
| Containerisation | Docker · Docker Compose                        |

---

## Architecture

### Deployment Topology

```mermaid
graph TD
    Browser("🌐 Browser")

    subgraph GitHub
        Repo["Source Code\n(GitHub)"]
        CI["GitHub Actions CI\ngo test · eslint · vite build"]
    end

    subgraph Railway["Railway (Production)"]
        direction TB
        FE["Frontend Service\nnginx + React SPA\nringr.up.railway.app\n(public HTTPS)"]
        BE["Backend Service\nGo API :8080\n(private network only)"]
        PG[("PostgreSQL\n(Railway managed)")]
        RD[("Redis\n(Railway managed)")]
    end

    Stripe["Stripe API\n(external)"]
    Cloudinary["Cloudinary\n(image storage)"]

    Browser -- "HTTPS" --> FE
    FE -- "static files" --> Browser
    FE -- "proxy /api/*\nbackend.railway.internal:8080" --> BE
    BE -- "pgx/v5 SQL queries" --> PG
    BE -- "cache · rate limiting" --> RD
    BE -- "payment processing" --> Stripe
    BE -- "image uploads" --> Cloudinary
    Repo -- "push to main\ntriggers redeploy" --> Railway
    Repo --> CI
```

> The backend has **no public URL** — all browser traffic enters through nginx on the frontend service and is forwarded over Railway's private network. nginx reloads every 30 seconds to pick up new backend IPs after redeployments.

---

### Backend Layer

| Layer | Responsibility |
|---|---|
| **Handler** | Decode request, call service, write JSON response, emit audit event |
| **Service** | Business rules, validation, orchestration |
| **Repository (interface)** | Decouples business logic from storage |
| **postgres/ impl** | pgx queries against PostgreSQL |
| **PhoneCache** | Optional Redis cache; falls back gracefully if unavailable |
| **AuditService** | Async fire-and-forget event logging to `audit_logs` table |

The frontend uses **React Context** for global auth state and a thin **API client** (`src/api/client.js`) that proxies all requests through nginx in production (or Vite's dev-server proxy locally).

---

### Backend Request Flow

```mermaid
flowchart TD
    Req["Incoming HTTP Request"]

    subgraph Middleware
        CORS["CORS"]
        RL["Rate Limiter\n(Redis sliding window)"]
        Auth["JWT Auth\n(RequireAuth / RequireAdmin)"]
    end

    subgraph Handler["Handler Layer"]
        H["Decode JSON\nValidate input\nWrite HTTP response\n+ fire audit event"]
    end

    subgraph Service["Service Layer"]
        S["Business logic\nValidation rules"]
    end

    subgraph Repository["Repository Layer"]
        I["Interface\n(decouples DB from logic)"]
        PG["postgres/ impl\npgx queries"]
        Cache["PhoneCache\n(Redis — optional)"]
    end

    DB[("PostgreSQL")]

    Req --> CORS --> RL --> Auth --> H --> S --> I
    I --> Cache
    I --> PG --> DB
```
---

### Cart Race Condition Prevention
```mermaid
sequenceDiagram
    participant A as User A
    participant B as User B
    participant DB as PostgreSQL

    A->>DB: SELECT FOR UPDATE (stock = 1)
    B->>DB: SELECT FOR UPDATE (stock = 1)
    Note over DB: User A gets lock — User B blocked, waiting

    DB-->>A: Lock acquired, stock = 1
    A->>A: Stock OK (1 ≥ 1)
    A->>DB: INSERT cart item + COMMIT
    Note over DB: Stock → 0, lock released

    DB-->>B: Lock acquired, stock = 0
    B->>B: Stock = 0, validation fails
    B-->>B: Return out-of-stock error
```

---

## Project Structure

```
CS464-g1t10-project/
├── tests/
│   └── load_test.js                  # k6 load & rate-limit tests
├── backend/
│   ├── cmd/api/
│   │   └── main.go                   # Entry point; seeds admin, ensures audit table
│   ├── internal/
│   │   ├── handler/
│   │   │   ├── auth_handler.go       # Register / Login
│   │   │   ├── phone_handler.go      # Phone CRUD
│   │   │   ├── cart_handler.go       # Cart get/add/remove
│   │   │   ├── payment_handler.go    # Stripe payment processing
│   │   │   ├── order_handler.go      # Order history
│   │   │   ├── upload_handler.go     # Cloudinary image upload
│   │   │   ├── audit_handler.go      # Audit log viewer (admin)
│   │   │   └── helpers.go            # clientIP helper
│   │   ├── middleware/
│   │   │   └── auth.go               # JWT RequireAuth / RequireAdmin
│   │   ├── model/
│   │   │   ├── user.go
│   │   │   ├── phone.go
│   │   │   ├── cart.go
│   │   │   ├── order.go
│   │   │   └── audit_log.go
│   │   ├── repository/
│   │   │   ├── user_repo.go          # Repository interfaces
│   │   │   ├── phone_repo.go
│   │   │   ├── cart_repo.go
│   │   │   ├── order_repo.go
│   │   │   ├── audit_log_repo.go
│   │   │   ├── phone_cache.go        # Redis phone cache
│   │   │   └── postgres/             # PostgreSQL implementations
│   │   │       ├── user_repository.go
│   │   │       ├── phone_repository.go
│   │   │       ├── cart_repository.go
│   │   │       ├── order_repository.go
│   │   │       └── audit_log_repository.go
│   │   ├── router/
│   │   │   └── router.go             # Route registration
│   │   └── service/
│   │       ├── auth_service.go
│   │       ├── phone_service.go
│   │       ├── cart_service.go
│   │       ├── payment_service.go
│   │       ├── order_service.go
│   │       └── audit_service.go
│   ├── migration/
│   │   └── 001_init.sql              # Schema + seed phones
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── go.mod
│
└── frontend/
    ├── nginx.conf.template            # nginx config template
    ├── start.sh                       # Startup script — periodic nginx reload
    ├── Dockerfile
    ├── vite.config.js
    └── src/
        ├── App.jsx                    # Router + top-level layout
        ├── main.jsx
        ├── api/
        │   └── client.js              # Typed API wrappers
        ├── components/
        │   ├── Navbar.jsx
        │   ├── ProtectedRoute.jsx     # Auth + admin route guards
        │   └── ui/                    # shadcn components
        ├── context/
        │   └── AuthContext.jsx        # JWT storage + auth state
        └── pages/
            ├── Home.jsx               # Phone listing + search
            ├── PhoneDetail.jsx        # Single phone + add to cart
            ├── Login.jsx
            ├── Register.jsx
            ├── Cart.jsx               # Cart management
            ├── Checkout.jsx           # Stripe payment checkout
            ├── Orders.jsx             # Order history
            ├── OrderDetail.jsx        # Order details
            ├── Admin.jsx              # Admin CRUD panel + image upload
            └── AuditLog.jsx           # Admin audit log viewer
```

---

## Database Schema

The base schema is in `backend/migration/001_init.sql`. The `audit_logs` table is created automatically at startup (`CREATE TABLE IF NOT EXISTS`).

### `users`

| Column           | Type         | Notes                         |
| ---------------- | ------------ | ----------------------------- |
| `id`           | SERIAL PK    |                               |
| `username`     | VARCHAR(100) | Unique, required              |
| `password`     | TEXT         | bcrypt-hashed                 |
| `phone_number` | VARCHAR(20)  |                               |
| `street`       | TEXT         | Address fields                |
| `city`         | VARCHAR(100) |                               |
| `state`        | VARCHAR(100) |                               |
| `country`      | VARCHAR(100) |                               |
| `zip_code`     | VARCHAR(20)  |                               |
| `role`         | VARCHAR(20)  | `'customer'` or `'admin'`   |

### `phones`

| Column          | Type          | Notes     |
| --------------- | ------------- | --------- |
| `id`          | SERIAL PK     |           |
| `brand`       | VARCHAR(100)  | Required  |
| `model`       | VARCHAR(100)  | Required  |
| `price`       | NUMERIC(10,2) | Required  |
| `stock`       | INT           | Default 0 |
| `description` | TEXT          |           |
| `image_url`   | TEXT          | Cloudinary URL or static path |

### `carts`

| Column      | Type        | Notes                             |
| ----------- | ----------- | --------------------------------- |
| `id`      | SERIAL PK   |                                   |
| `user_id` | INT FK      | References `users(id)`          |
| `status`  | VARCHAR(20) | `'active'` or `'checked_out'` |

### `cart_items`

| Column       | Type          | Notes                                          |
| ------------ | ------------- | ---------------------------------------------- |
| `id`       | SERIAL PK     |                                                |
| `cart_id`  | INT FK        | References `carts(id)` — cascades on delete  |
| `phone_id` | INT FK        | References `phones(id)`                      |
| `quantity` | INT           | Required                                       |
| `price`    | NUMERIC(10,2) | Price at time of add                           |

### `orders`

| Column         | Type          | Notes                    |
| -------------- | ------------- | ------------------------ |
| `id`         | TEXT PK       | Stripe payment intent ID |
| `user_id`    | INT FK        | References `users(id)` |
| `status`     | VARCHAR(20)   | Default `'succeeded'`  |
| `total`      | NUMERIC(10,2) | Order total              |
| `created_at` | TIMESTAMP     | Order creation time      |

### `order_items`

| Column         | Type          | Notes                        |
| -------------- | ------------- | ---------------------------- |
| `id`         | SERIAL PK     |                              |
| `order_id`   | TEXT FK       | References `orders(id)`    |
| `phone_id`   | INT FK        | References `phones(id)`    |
| `phone_name` | TEXT          | Phone model at time of order |
| `quantity`   | INT           | Required                     |
| `price`      | NUMERIC(10,2) | Price at time of order       |

### `audit_logs` *(auto-created at startup)*

| Column          | Type         | Notes                             |
| --------------- | ------------ | --------------------------------- |
| `id`          | SERIAL PK    |                                   |
| `user_id`     | INT FK       | References `users(id)` — nullable |
| `action`      | VARCHAR(64)  | e.g. `phone.created`, `user.login_failed` |
| `resource_type` | VARCHAR(32) | e.g. `phone`, `user`, `order`   |
| `resource_id` | VARCHAR(64)  | ID of the affected resource       |
| `details`     | JSONB        | Action-specific metadata          |
| `ip_address`  | VARCHAR(45)  | Client IP (respects proxy headers) |
| `created_at`  | TIMESTAMPTZ  | Event timestamp                   |

---

## API Reference

Base URL: `http://localhost:8080`

All protected routes require the header:

```
Authorization: Bearer <jwt_token>
```

### Auth

| Method | Path               | Auth | Description                   |
| ------ | ------------------ | ---- | ----------------------------- |
| POST   | `/auth/register` | None | Create a new customer account |
| POST   | `/auth/login`    | None | Login and receive a JWT token |

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

| Method | Path             | Auth       | Description            |
| ------ | ---------------- | ---------- | ---------------------- |
| GET    | `/phones`      | None       | List all phones        |
| GET    | `/phones/{id}` | None       | Get a single phone     |
| POST   | `/phones`      | Admin only | Create a phone listing |
| PUT    | `/phones/{id}` | Admin only | Update a phone listing |
| DELETE | `/phones/{id}` | Admin only | Delete a phone listing |

---

### Upload

| Method | Path       | Auth       | Description                              |
| ------ | ---------- | ---------- | ---------------------------------------- |
| POST   | `/upload` | Admin only | Upload an image to Cloudinary, returns URL |

```
// multipart/form-data with field: image
// Response 200
{ "url": "https://res.cloudinary.com/..." }
```

---

### Cart

All cart routes require authentication.

| Method | Path               | Auth          | Description                 |
| ------ | ------------------ | ------------- | --------------------------- |
| GET    | `/cart`          | Authenticated | Get current user's cart     |
| POST   | `/cart`          | Authenticated | Add item to cart            |
| DELETE | `/cart/{itemId}` | Authenticated | Remove item from cart       |

---

### Payment

| Method | Path     | Auth          | Description                     |
| ------ | -------- | ------------- | ------------------------------- |
| POST   | `/pay` | Authenticated | Process Stripe payment for cart |

---

### Orders

| Method | Path             | Auth          | Description               |
| ------ | ---------------- | ------------- | ------------------------- |
| GET    | `/orders`      | Authenticated | List user's order history |
| GET    | `/orders/{id}` | Authenticated | Get order details by ID   |

---

### Audit Log

| Method | Path            | Auth       | Description                              |
| ------ | --------------- | ---------- | ---------------------------------------- |
| GET    | `/audit-logs` | Admin only | Recent audit events, newest first. Supports `?limit=` and `?offset=` |

---

## Getting Started

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) — the only requirement to run the backend and database

For frontend development only:

- **Node.js** 18 or later + npm — [nodejs.org](https://nodejs.org)

---

### 1. Configure environment variables

The default values in `backend/docker-compose.yml` work out of the box for local development. To change them edit the `environment` section under the `api` service:

| Variable           | Default                                            | Description                           |
| ------------------ | -------------------------------------------------- | ------------------------------------- |
| `DATABASE_URL`   | `postgres://postgres:postgres@db:5432/phones_db` | PostgreSQL connection string          |
| `JWT_SECRET`     | `change-me-in-production`                        | Secret key for signing JWT tokens     |
| `ADMIN_PASSWORD` | `adminpassword`                                  | Password for the seeded admin account |

---

### 2. Start with Docker

```bash
cd backend
docker compose up --build -d
```

This will:

1. Build the Go binary inside a Docker build stage
2. Start a PostgreSQL 16 container and automatically run `migration/001_init.sql`
3. Start the API container once the database is healthy
4. Seed an `admin` user using the `ADMIN_PASSWORD` environment variable
5. Create the `audit_logs` table if it does not already exist

Verify it is running:

```bash
docker compose logs -f api
```

You should see:

```
Connected to database
audit_logs table ready
Admin user seeded
Server running on :8080
```

| Service     | URL                   |
| ----------- | --------------------- |
| Backend API | http://localhost:8080 |
| PostgreSQL  | localhost:5432        |

**Stopping the containers:**

```bash
# Stop containers
docker compose down

# Stop and wipe the database for a fresh start
docker compose down -v
```

**Starting the frontend:**

```bash
cd frontend
npm install
npm run dev
```

| Service  | URL                   |
| -------- | --------------------- |
| Frontend | http://localhost:5173 |
| Backend  | http://localhost:8080 |

---

### 3. Running without Docker

Requires a local PostgreSQL instance.

```bash
cd backend
export DATABASE_URL="postgres://postgres:password@localhost:5432/phones_db"
export JWT_SECRET="your-secret-key"
export ADMIN_PASSWORD="adminpassword"

psql -U postgres -d phones_db -f migration/001_init.sql
go run ./cmd/api
```

---

## Environment Variables

### Backend

| Variable              | Required | Description                                                             |
| --------------------- | -------- | ----------------------------------------------------------------------- |
| `DATABASE_URL`      | Yes      | PostgreSQL DSN, e.g. `postgres://user:pass@localhost:5432/phones_db`  |
| `JWT_SECRET`        | Yes      | Secret used to sign JWT tokens — change this in production             |
| `ADMIN_PASSWORD`    | Yes      | Password for the seeded `admin` account                               |
| `STRIPE_SECRET_KEY` | Yes      | Stripe secret key for payment processing                               |
| `CLOUDINARY_URL`    | Yes      | Cloudinary DSN — format: `cloudinary://API_KEY:API_SECRET@CLOUD_NAME` |
| `REDIS_URL`         | No       | Redis connection URL — caching and rate limiting disabled if absent    |
| `STRIPE_CURRENCY`   | No       | Payment currency code (default: `sgd`)                               |
| `ALLOWED_ORIGIN`    | No       | CORS allowed origin (default: `*`)                                   |

### Frontend (build-time)

| Variable        | Required | Description                                            |
| --------------- | -------- | ------------------------------------------------------ |
| `VITE_API_URL` | No       | API base URL (default: `/api` — proxied through nginx) |

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
  - Add phones to cart (race-condition safe via `SELECT FOR UPDATE`)
  - Remove individual items from cart
  - Checkout cart with Stripe payment integration
  - View order history and order details

- **Admin**

  - All customer features
  - Create, edit, and delete phone listings from the Admin panel
  - Upload product images via Cloudinary (file picker + live preview, or paste a URL)
  - View a live audit log of all system events at `/audit-logs`
  - Role assigned at account seed time or via direct database update

- **General**

  - JWT-based stateless authentication
  - Role-based access control enforced on both backend middleware and frontend route guards
  - Password hashing with bcrypt
  - Audit logging — key events (logins, catalog changes, payments) recorded asynchronously to `audit_logs`
  - Redis caching for phone catalog (10-minute TTL, invalidated on write)
  - Sliding-window rate limiting via Redis (global + per-user limits on sensitive endpoints)
  - Responsive UI with dark-mode CSS variables
  - Fully containerised backend and database with Docker
  - Stock deduction only after successful Stripe payment

---

## Load Testing

[k6](https://k6.io) load tests live in `tests/load_test.js`. They cover two scenarios:

| Scenario | Endpoint | What it proves |
|---|---|---|
| Phone listing | `GET /phones` | Public catalog is reachable and returns a valid JSON array |
| Login throttle | `POST /auth/login` (wrong password) | Per-user rate limiter fires 429 after 5 failed attempts/min |

### Install k6

```bash
# Windows
winget install k6

# macOS
brew install k6
```

### Run locally (against docker-compose)

Start the backend first (`docker compose up -d` in `backend/`), then:

```bash
k6 run tests/load_test.js
```

### Run against the deployed Railway app

```bash
k6 run -e BASE_URL=https://ringr.up.railway.app/api tests/load_test.js
```

### Staged stress test (rate-limit demo)

To demonstrate the sliding-window throttle kicking in with `429 Too Many Requests`, swap the `options` block in `tests/load_test.js` — uncomment **Option 2** and comment out **Option 1**, then re-run:

```bash
k6 run -e BASE_URL=https://ringr.up.railway.app/api tests/load_test.js
```

The test ramps up to 1 000 virtual users over ~2 minutes. You will see 429 responses once the global (500 req/min) and per-user (5 login attempts/min) throttles are breached.

---

## Postman Collection

A ready-to-import Postman collection is provided at:

```
backend/phones-api.postman_collection.json
```

Import it into Postman, set the `baseUrl` collection variable to `http://localhost:8080`, then run the requests in order — the Login requests automatically save tokens to collection variables so subsequent requests are authenticated automatically.
