# Ringr Mobile — Frontend

React + Vite frontend for the Ringr Mobile phone e-commerce application.

---

## Tech Stack

| Tool              | Version  | Purpose                              |
|-------------------|----------|--------------------------------------|
| React             | 19       | UI framework                         |
| React Router      | 7        | Client-side routing                  |
| Vite              | 8        | Build tool + dev server              |
| Tailwind CSS      | v4       | Utility-first styling                |
| shadcn/ui         | Base UI  | Pre-built accessible components      |
| Lucide React      | —        | Icon library                         |
| Sonner            | —        | Toast notifications                  |

---

## Project Structure

```
frontend/
├── index.html
├── vite.config.js           # Tailwind plugin, @ alias, API proxy
├── components.json          # shadcn configuration
├── jsconfig.json            # @ path alias for IDE
├── src/
│   ├── App.jsx              # Router setup and top-level layout
│   ├── main.jsx             # React entry point
│   ├── index.css            # Tailwind + shadcn CSS variables
│   ├── api/
│   │   └── client.js        # All API calls (auth, phones, cart)
│   ├── components/
│   │   ├── Navbar.jsx       # Top navigation bar
│   │   ├── ProtectedRoute.jsx  # Auth + admin route guards
│   │   └── ui/              # shadcn components (button, card, etc.)
│   ├── context/
│   │   └── AuthContext.jsx  # Global auth state (JWT in localStorage)
│   ├── lib/
│   │   └── utils.js         # cn() helper (clsx + tailwind-merge)
│   └── pages/
│       ├── Home.jsx         # Phone catalog with search
│       ├── PhoneDetail.jsx  # Single phone + add to cart
│       ├── Login.jsx        # Login form
│       ├── Register.jsx     # Registration form with address fields
│       ├── Cart.jsx         # Cart management + checkout
│       └── Admin.jsx        # Admin panel — phone CRUD
└── package.json
```

---

## Pages & Routes

| Route          | Page          | Access         | Description                        |
|----------------|---------------|----------------|------------------------------------|
| `/`            | Home          | Public         | Browse and search all phones       |
| `/phones/:id`  | PhoneDetail   | Public         | Phone details, add to cart         |
| `/login`       | Login         | Public         | Log in with username + password    |
| `/register`    | Register      | Public         | Create a new customer account      |
| `/cart`        | Cart          | Authenticated  | View cart, remove items, checkout  |
| `/admin`       | Admin         | Admin only     | Create, edit, and delete phones    |

---

## Getting Started

### Install dependencies

```bash
npm install
```

### Development

```bash
npm run dev
```

Starts at `http://localhost:5173`. All `/api/*` requests are proxied to the Go backend at `http://127.0.0.1:8080` — make sure the backend is running first.

### Build for production

```bash
npm run build
# Output → dist/
```

### Preview production build

```bash
npm run preview
```

### Lint

```bash
npm run lint
```

---

## API Proxy

The Vite dev server is configured to proxy API calls so the frontend never needs to know the backend's port:

```
/api/auth/login  →  http://127.0.0.1:8080/auth/login
/api/phones      →  http://127.0.0.1:8080/phones
/api/cart        →  http://127.0.0.1:8080/cart
```

This is configured in `vite.config.js`. In production, point your reverse proxy (nginx, etc.) to do the same.

---

## Auth Flow

1. User logs in → backend returns a JWT.
2. JWT is stored in `localStorage` and decoded client-side to read `username`, `role`, and `exp`.
3. `AuthContext` exposes `isAuthenticated`, `isAdmin`, `user`, `login`, `logout`, and `register`.
4. Every API request automatically attaches `Authorization: Bearer <token>` via `client.js`.
5. `ProtectedRoute` and `AdminRoute` redirect unauthenticated / unauthorized users away from protected pages.

---

## Adding shadcn Components

This project uses the **Base UI** variant of shadcn (not Radix UI). To add a new component:

```bash
npx shadcn@latest add <component-name>
```

> **Note:** Base UI does not support the `asChild` prop. Use `buttonVariants()` directly on a `<Link>` element instead of `<Button asChild>`.
