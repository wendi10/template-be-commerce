# 🛒 E-Commerce Backend — Go + MySQL

Production-ready e-commerce REST API built with Go, MySQL, and Clean Architecture. Supports a public storefront, admin dashboard, and payment gateway integration (Doku).

---

## 📦 Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.21 |
| Database | MySQL 8.4 |
| Router | Chi v5 |
| Auth | JWT (HS256) |
| Password | bcrypt |
| ORM / DB  | GORM v2 (MySQL driver) |
| Migration | goose |
| Logging | zap |
| Validation | go-playground/validator |
| Decimal | shopspring/decimal |
| Container | Docker + docker-compose |

---

## 🏗 Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go              # Entry point & DI wiring
├── config/
│   └── config.go                # Config loader (env vars)
├── internal/
│   ├── domain/                  # Pure domain models + DTOs
│   │   ├── user.go
│   │   ├── customer.go
│   │   ├── product.go
│   │   ├── banner.go
│   │   ├── promo.go
│   │   ├── cart.go
│   │   ├── order.go
│   │   └── payment.go
│   ├── repository/
│   │   ├── interfaces.go        # Repository interfaces
│   │   └── mysql/               # MySQL implementations
│   ├── usecase/                 # Business logic
│   │   ├── auth_usecase.go
│   │   ├── customer_usecase.go
│   │   ├── product_usecase.go
│   │   ├── banner_usecase.go
│   │   ├── promo_usecase.go
│   │   ├── cart_usecase.go
│   │   ├── order_usecase.go
│   │   └── payment_usecase.go
│   ├── handler/
│   │   └── http/                # HTTP handlers + router
│   ├── middleware/              # JWT auth, logging, recovery
│   └── payment/
│       ├── gateway.go           # Payment gateway interface
│       └── doku/                # Doku provider implementation
├── migrations/                  # goose SQL migrations
├── pkg/
│   ├── apperrors/               # Typed application errors
│   ├── jwt/                     # JWT manager
│   ├── logger/                  # Zap logger
│   ├── pagination/              # Pagination helpers
│   ├── response/                # HTTP response envelope
│   ├── slug/                    # Slug generator
│   └── validator/               # Struct validation
├── Dockerfile
├── docker-compose.yml
└── env.example
```

---

## 🚀 Quick Start

### Prerequisites
- Go 1.22+
- Docker & Docker Compose

### 1. Clone & configure
```bash
git clone <repo-url>
cd template-be-commerce

# Copy and edit the env file
cp env.example .env
```

### 2. Run with Docker Compose
```bash
docker-compose up --build
```

### 3. Run locally (without Docker)
```bash
# Start only the database
docker-compose up mysql -d

# Download dependencies
go mod tidy

# Run the API
go run ./cmd/api
```

The API will be available at `http://localhost:8080`.

---

## 🔑 Default Admin Credentials

After migrations run, a seed admin account is created:
- **Email:** `admin@example.com`
- **Password:** `Admin@1234`

> ⚠️ Change these credentials immediately in production.

---

## 🗄 Database ERD

```
users
  id (PK)           uuid
  email             varchar UNIQUE
  password_hash     text
  role              enum(customer|admin)
  is_active         boolean
  created_at        timestamptz
  updated_at        timestamptz
  deleted_at        timestamptz (soft delete)

customers
  id (PK)           uuid
  user_id (FK)      → users.id
  first_name        varchar
  last_name         varchar
  phone             varchar
  avatar            text
  ...timestamps

addresses
  id (PK)           uuid
  customer_id (FK)  → customers.id
  label             varchar
  recipient_name    varchar
  phone             varchar
  address_line1/2   varchar
  city, province, postal_code
  is_default        boolean
  ...timestamps

categories
  id (PK)           uuid
  name, slug        varchar
  description       text
  parent_id (FK)    → categories.id  (self-ref for subcategories)
  is_active         boolean
  ...timestamps

products
  id (PK)           uuid
  category_id (FK)  → categories.id
  name, slug        varchar
  description       text
  price             numeric(15,2)
  weight            numeric(10,3)
  stock             integer
  is_active         boolean
  ...timestamps

product_images
  id (PK)           uuid
  product_id (FK)   → products.id
  url               text
  alt_text          varchar
  is_primary        boolean
  sort_order        integer

banners
  id (PK)           uuid
  title, subtitle   varchar
  image_url         text
  link_url          text
  is_active         boolean
  sort_order        integer
  start_date/end_date  timestamptz
  ...timestamps

promo_codes
  id (PK)           uuid
  code              varchar UNIQUE
  name, description varchar/text
  type              enum(percentage|fixed)
  value             numeric
  min_purchase      numeric
  max_discount      numeric
  usage_limit       integer  (0 = unlimited)
  used_count        integer
  start_date/end_date  timestamptz
  is_active         boolean
  ...timestamps

cart_items
  id (PK)           uuid
  customer_id (FK)  → customers.id
  product_id (FK)   → products.id
  quantity          integer
  UNIQUE(customer_id, product_id)
  ...timestamps

orders
  id (PK)           uuid
  customer_id (FK)  → customers.id
  address_id (FK)   → addresses.id
  promo_code_id (FK)→ promo_codes.id (nullable)
  order_number      varchar UNIQUE
  status            enum(pending|waiting_payment|paid|processing|shipped|completed|cancelled)
  sub_total         numeric
  discount_amount   numeric
  shipping_cost     numeric
  total_amount      numeric
  notes             text
  ...timestamps

order_items
  id (PK)           uuid
  order_id (FK)     → orders.id
  product_id (FK)   → products.id
  product_name      varchar  (snapshot)
  product_price     numeric  (snapshot)
  quantity          integer
  total_price       numeric

payments
  id (PK)             uuid
  order_id (FK)       → orders.id
  payment_method      varchar
  payment_provider    enum(doku)
  amount              numeric
  status              enum(pending|success|failed|expired|refunded)
  transaction_id      varchar
  payment_url         text
  callback_data       text (raw JSON)
  paid_at             timestamptz
  expired_at          timestamptz
  ...timestamps
```

---

## 📡 API Routes

### Auth
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/auth/register` | — | Register customer |
| POST | `/api/v1/auth/login` | — | Customer login |
| POST | `/api/v1/auth/refresh` | — | Refresh token |
| POST | `/api/v1/admin/auth/login` | — | Admin login |

### Customer Profile
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/me` | Customer | Get profile |
| PUT | `/api/v1/me` | Customer | Update profile |
| GET | `/api/v1/me/addresses` | Customer | List addresses |
| POST | `/api/v1/me/addresses` | Customer | Create address |
| PUT | `/api/v1/me/addresses/{id}` | Customer | Update address |
| DELETE | `/api/v1/me/addresses/{id}` | Customer | Delete address |
| PATCH | `/api/v1/me/addresses/{id}/default` | Customer | Set default |

### Products (Public)
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/products` | — | List products (paginated) |
| GET | `/api/v1/products/{id}` | — | Get product by ID |
| GET | `/api/v1/products/slug/{slug}` | — | Get product by slug |
| GET | `/api/v1/categories` | — | List categories |

### Banners & Promos (Public)
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/banners` | — | List active banners |
| POST | `/api/v1/promos/validate` | — | Validate promo code |

### Cart
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/cart` | Customer | Get cart + summary |
| POST | `/api/v1/cart` | Customer | Add item |
| PUT | `/api/v1/cart/{itemID}` | Customer | Update quantity |
| DELETE | `/api/v1/cart/{itemID}` | Customer | Remove item |

### Orders
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/orders` | Customer | Create order from cart |
| GET | `/api/v1/orders` | Customer | List my orders |
| GET | `/api/v1/orders/{id}` | Customer | Get order detail |
| POST | `/api/v1/orders/{id}/cancel` | Customer | Cancel order |

### Payments
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/payments` | Customer | Create payment |
| GET | `/api/v1/payments/order/{orderID}` | Customer | Get payment by order |
| POST | `/api/v1/payments/callback/{provider}` | — (signed) | Payment webhook |

### Admin — Products
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/admin/products` | Admin | List all products |
| POST | `/api/v1/admin/products` | Admin | Create product |
| PUT | `/api/v1/admin/products/{id}` | Admin | Update product |
| DELETE | `/api/v1/admin/products/{id}` | Admin | Delete product |
| POST | `/api/v1/admin/products/{id}/images` | Admin | Add image |
| DELETE | `/api/v1/admin/products/{id}/images/{imageID}` | Admin | Delete image |
| GET | `/api/v1/admin/categories` | Admin | List categories |
| POST | `/api/v1/admin/categories` | Admin | Create category |
| PUT | `/api/v1/admin/categories/{id}` | Admin | Update category |
| DELETE | `/api/v1/admin/categories/{id}` | Admin | Delete category |

### Admin — Banners & Promos
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/admin/banners` | Admin | List banners |
| POST | `/api/v1/admin/banners` | Admin | Create banner |
| PUT | `/api/v1/admin/banners/{id}` | Admin | Update banner |
| DELETE | `/api/v1/admin/banners/{id}` | Admin | Delete banner |
| GET | `/api/v1/admin/promos` | Admin | List promos |
| POST | `/api/v1/admin/promos` | Admin | Create promo |
| PUT | `/api/v1/admin/promos/{id}` | Admin | Update promo |
| DELETE | `/api/v1/admin/promos/{id}` | Admin | Delete promo |

### Admin — Orders & Reports
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/admin/orders` | Admin | List all orders |
| GET | `/api/v1/admin/orders/{id}` | Admin | Order detail |
| PATCH | `/api/v1/admin/orders/{id}/status` | Admin | Update order status |
| GET | `/api/v1/admin/customers` | Admin | List customers |
| GET | `/api/v1/admin/reports/sales` | Admin | Sales report |

---

## 📨 Example Requests & Responses

### Register Customer
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "Secure@1234",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+62812345678"
}
```
```json
{
  "success": true,
  "message": "registration successful",
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci...",
    "expires_at": "2026-03-05T15:00:00Z",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john@example.com",
      "role": "customer"
    }
  }
}
```

### Add to Cart
```http
POST /api/v1/cart
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "product_id": "a1b2c3d4-...",
  "quantity": 2
}
```

### Checkout (Create Order)
```http
POST /api/v1/orders
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "address_id": "a1b2c3d4-...",
  "promo_code": "SAVE10",
  "notes": "Please ring the bell"
}
```
```json
{
  "success": true,
  "message": "order created",
  "data": {
    "id": "...",
    "order_number": "ORD-1709654321000",
    "status": "pending",
    "sub_total": "150000.00",
    "discount_amount": "15000.00",
    "shipping_cost": "15000.00",
    "total_amount": "150000.00",
    "items": [...]
  }
}
```

### Create Payment
```http
POST /api/v1/payments
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "order_id": "...",
  "payment_method": "virtual_account",
  "provider": "doku"
}
```
```json
{
  "success": true,
  "message": "payment created",
  "data": {
    "order": { ... },
    "payment": { "status": "pending", ... },
    "payment_url": "https://checkout.doku.com/..."
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": 400,
    "message": "validation error",
    "detail": "email is required; password must be at least 8 characters"
  }
}
```

---

## 🔒 Order Status Flow

```
pending → waiting_payment → paid → processing → shipped → completed
   │               │          │
   └───────────────┴──────────┴──→ cancelled
```

---

## 💳 Payment Gateway

The payment system uses an **abstraction interface** (`payment.Gateway`) so new providers can be added without touching business logic:

```go
type Gateway interface {
    CreateTransaction(ctx, req) (*TransactionResponse, error)
    HandleCallback(ctx, payload, headers) (*CallbackResult, error)
    ProviderName() domain.PaymentProvider
}
```

Currently supported:
- **Doku** — `internal/payment/doku/`

To add a new provider (e.g., Midtrans):
1. Create `internal/payment/midtrans/midtrans.go`
2. Implement the `Gateway` interface
3. Register in `cmd/api/main.go` gateways map
4. Add the enum value to `domain.PaymentProvider`

---

## 🌍 Environment Variables

See `env.example` for all available configuration options.

Key variables:
- `JWT_ACCESS_SECRET` / `JWT_REFRESH_SECRET` — **must** be changed in production
- `DOKU_CLIENT_ID` / `DOKU_SECRET_KEY` — from your Doku dashboard
- `PAYMENT_CALLBACK_URL` — publicly accessible URL for payment webhooks

---

## 🧪 Running Tests

```bash
go test ./...
```

---

## 🐳 Docker

### Build image
```bash
docker build -t ecommerce-api .
```

### Full stack
```bash
# Start all services
docker-compose up -d

# With pgAdmin (dev profile)
docker-compose --profile dev up -d

# View logs
docker-compose logs -f api

# Stop
docker-compose down
```

---

## 📈 Scalability Notes

- **Connection pooling** via pgx pool (configurable `DB_MAX_OPEN_CONNS`)
- **Soft deletes** on all major entities preserve data integrity
- **Stateless JWT** auth enables horizontal scaling behind a load balancer
- **Payment abstraction** allows adding providers without changing core logic
- **Separate admin/public routes** allow independent rate limiting and auth policies
- All DB queries use **parameterised statements** (no SQL injection risk)
- **Paginated list endpoints** with configurable limits (max 100)
