# RestoBill — Progress Tracker

## Project Overview

Restaurant billing system built with Go (Gin) and PostgreSQL, following clean architecture patterns with domain/dto/ports/service/handler layering.

**Repository:** https://github.com/Abhiram123897/RestoBill

---

## Week 1: Database + Project Setup ✅

### Phase 1: Project Skeleton ✅
- Initialized Go module (`go mod init restobill`)
- Created folder structure: `cmd/`, `internal/`, `migrations/`
- Created `.gitignore`, `.env`, `Makefile`
- First commit + pushed to GitHub

### Phase 2: Docker + Database ✅
- Created `docker-compose.yml` (PostgreSQL 16)
- Switched to local PostgreSQL + pgAdmin for development
- Database `restobill` created and accessible

### Phase 3: Database Schema ✅
- 8 migration files created and run:
  - `001_create_users`
  - `002_create_categories`
  - `003_create_menu_items`
  - `004_create_restaurant_tables`
  - `005_create_orders`
  - `006_create_order_items`
  - `007_create_bills`
  - `008_create_payments`
- All tables verified in pgAdmin

### Phase 4: Server + DB Connection ✅
- `internal/logger/logger.go` — Zap structured logging
- `internal/errs/errors.go` — AppError with typed codes
- `internal/db/db.go` — pgxpool connection + PostgreSQL error helpers
- `internal/handler/health_handler.go` — Health endpoint pings DB
- `cmd/server/main.go` — Graceful shutdown
- Health endpoint tested: `GET /api/v1/health` → `{"status":"ok","database":"connected"}`

### Phase 5: Verify & Commit ✅
- All 8 tables verified
- Health endpoint working
- Pushed to GitHub

---

## Week 2: Menu Management API ✅

### Phase 1: Category Domain + Repository ✅
- `internal/domain/category.go` — Domain model
- `internal/dto/category_dto.go` — Request/Response DTOs
- `internal/repository/interfaces.go` — Repository interfaces (port pattern)
- `internal/repository/category_repo.go` — PostgreSQL implementation

### Phase 2: Category Service + Handler + Routes ✅
- `internal/service/category_service.go` — Business logic
- `internal/handler/category_handler.go` — HTTP handlers with UUID validation
- Routes wired in `internal/router/router.go`
- All CRUD endpoints tested via Postman

### Phase 3: Menu Item Domain + Repository ✅
- `internal/domain/menu_item.go` — Domain model
- `internal/dto/menu_item_dto.go` — Request/Response DTOs
- `internal/repository/menu_item_repo.go` — PostgreSQL implementation

### Phase 4: Menu Item Service + Handler + Routes ✅
- `internal/service/menu_item_service.go` — Business logic
- `internal/handler/menu_item_handler.go` — HTTP handlers with UUID validation
- Routes wired including filter by category + availability toggle
- All CRUD endpoints tested via Postman

### Phase 5: Seed Data + Verify + Commit ✅
- `seeds/seed.sql` — Sample data (4 categories, 12 menu items)
- Additional migrations:
  - `009_add_case_insensitive_category_name` — Case-insensitive unique names
  - `010_add_unique_menu_item_per_category` — Unique item name per category
- All endpoints verified

---

## Week 3: Order Management API ✅

### Phase 1: Restaurant Tables API ✅
- `internal/domain/restaurant_table.go` — Domain model
- `internal/dto/restaurant_table_dto.go` — Request/Response DTOs
- `internal/repository/restaurant_table_repo.go` — PostgreSQL implementation
- `internal/service/restaurant_table_service.go` — Business logic with status validation
- `internal/handler/restaurant_table_handler.go` — HTTP handlers
- All CRUD + status update endpoints tested via Postman

### Phase 2: Order Domain + Repository ✅
- `internal/domain/order.go` — Order + OrderItem domain models
- `internal/dto/order_dto.go` — Request/Response DTOs
- `internal/repository/order_repo.go` — PostgreSQL implementation with active order check

### Phase 3: Order Service + Handler + Routes ✅
- `internal/service/order_service.go` — Business logic:
  - Order type validation (dine_in requires table, takeaway/delivery no table)
  - Status transitions (open → served/billed/paid, can't go backward)
  - Price snapshot when adding items
  - Menu item availability check
  - Active order per table check (prevents duplicate orders on same table)
- `internal/handler/order_handler.go` — HTTP handlers with UUID + quantity validation
- Order filtering by status and order_type (`?status=open&order_type=dine_in`)

### Phase 4: Table Status Sync ✅
- Table auto-set to "occupied" when dine-in order created
- Table auto-set to "available" when order marked as "paid"
- No manual table status changes needed by waiter

### Phase 5: Verify + Commit ✅
- All order endpoints verified
- Edge cases tested (invalid status transitions, quantity validation, duplicate table orders)

---

## API Endpoints (Current)

### Health
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check + DB ping |

### Categories
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/categories` | Create category |
| GET | `/api/v1/categories` | List all categories |
| GET | `/api/v1/categories/:id` | Get category by ID |
| PUT | `/api/v1/categories/:id` | Update category |
| DELETE | `/api/v1/categories/:id` | Delete category |

### Menu Items
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/menu-items` | Create menu item |
| GET | `/api/v1/menu-items` | List all menu items |
| GET | `/api/v1/menu-items?category_id=xxx` | Filter by category |
| GET | `/api/v1/menu-items/:id` | Get menu item by ID |
| PUT | `/api/v1/menu-items/:id` | Update menu item |
| PATCH | `/api/v1/menu-items/:id/availability` | Toggle availability |
| DELETE | `/api/v1/menu-items/:id` | Delete menu item |

### Restaurant Tables
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tables` | Create table |
| GET | `/api/v1/tables` | List all tables |
| GET | `/api/v1/tables/:id` | Get table by ID |
| PUT | `/api/v1/tables/:id` | Update table |
| PATCH | `/api/v1/tables/:id/status` | Update table status |
| DELETE | `/api/v1/tables/:id` | Delete table |

### Orders
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/orders` | Create order |
| GET | `/api/v1/orders` | List all orders |
| GET | `/api/v1/orders?status=open` | Filter by status |
| GET | `/api/v1/orders?order_type=dine_in` | Filter by type |
| GET | `/api/v1/orders/:id` | Get order with items |
| PATCH | `/api/v1/orders/:id/status` | Update order status |
| POST | `/api/v1/orders/:id/items` | Add item to order |
| DELETE | `/api/v1/orders/:id/items/:itemId` | Remove item from order |

---

## Architecture

```
cmd/server/main.go          → Entry point + graceful shutdown
internal/
├── config/                  → Environment configuration
├── db/                      → Database connection + PG error helpers
├── logger/                  → Zap structured logging
├── errs/                    → AppError types (400, 404, 409, 500)
├── domain/                  → Domain models (internal)
├── dto/                     → Request/Response objects (API-facing)
├── repository/
│   ├── interfaces.go        → Port interfaces
│   ├── category_repo.go     → PostgreSQL implementation
│   ├── menu_item_repo.go    → PostgreSQL implementation
│   ├── restaurant_table_repo.go
│   └── order_repo.go
├── service/                 → Business logic layer
├── handler/                 → HTTP handlers (Gin)
├── router/                  → Route wiring
└── middleware/              → (Week 5 — Auth)
```

---

## Patterns Adopted

| Pattern | Used In |
|---------|---------|
| Domain/DTO separation | domain/, dto/ |
| Port interfaces (repos) | repository/interfaces.go |
| Service interface + struct | service/ |
| Handler struct holds service | handler/ |
| `*errs.AppError` return type | All layers |
| Error pass-through (repo→service→handler) | service/ |
| `AsPgError` + `pgerrcode` for PG errors | db/db.go, repository/ |
| UUID validation in handler | handler/ |
| Zap structured logger | logger/ |
| Graceful shutdown | main.go |

---

## Tech Stack

### Backend
| Layer | Technology |
|-------|-----------|
| Language | Go 1.25+ |
| Framework | Gin |
| Database | PostgreSQL 16 |
| DB Driver | pgx/v5 |
| Migrations | golang-migrate |
| Logger | Zap |
| UUID | google/uuid |
| Error Codes | jackc/pgerrcode |
| Auth | JWT (golang-jwt/jwt/v5) + bcrypt |

### Frontend (Week 6+)
| Concern | Technology | Purpose |
|---------|-----------|---------|
| Framework | React + Vite | Fast, component-based UI |
| Styling | Tailwind CSS + shadcn/ui | Dense POS tables/forms/modals |
| Server State | TanStack Query | Caching, refetching, loading/error states for menu/orders/bills/tables |
| Client State | Context API | Auth/user/role management |
| HTTP Client | axios (with interceptor) | Auto-attach JWT, handle 401 globally |

---

## Week 4: Billing & Payments API ✅

### Phase 1: Bill Domain + Repository ✅
- `internal/domain/bill.go` — Domain model
- `internal/dto/bill_dto.go` — Request/Response DTOs
- `internal/repository/bill_repo.go` — PostgreSQL implementation
- Additional migrations:
  - `011_add_cancelled_status_to_orders` — Added "cancelled" order status
  - `012_allow_rebilling_after_void` — Partial unique index for rebilling

### Phase 2: Bill Service + Handler ✅
- `internal/service/bill_service.go` — Business logic:
  - Auto-calculate subtotal from order items
  - Tax calculation (configurable rate)
  - Discount application with reason
  - Void bill (reverts order to open for correction)
  - Can't bill already billed/paid orders
  - Can't bill orders with no items
  - Regenerate bill after void
- `internal/handler/bill_handler.go` — HTTP handlers

### Phase 3: Payment Domain + Repository ✅
- `internal/domain/payment.go` — Domain model
- `internal/dto/payment_dto.go` — Request/Response DTOs
- `internal/repository/payment_repo.go` — PostgreSQL implementation

### Phase 4: Payment Service + Handler ✅
- `internal/service/payment_service.go` — Business logic:
  - Split payments supported (cash + card + upi)
  - Can't overpay (validates remaining balance)
  - Auto-marks bill as "paid" when fully paid
  - Auto-marks order as "paid"
  - Auto-frees table (dine-in)
  - Can't pay void/already paid bills
- `internal/handler/payment_handler.go` — HTTP handlers
- Rounded total to nearest integer (no decimals for customer)

### Phase 5: Verify + Commit ✅
- Full end-to-end flow tested (order → bill → payment → table freed)
- Split payment tested
- Error cases tested (overpay, void, invalid method)

---

## API Endpoints (Current)

### Health
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check + DB ping |

### Categories
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/categories` | Create category |
| GET | `/api/v1/categories` | List all categories |
| GET | `/api/v1/categories/:id` | Get category by ID |
| PUT | `/api/v1/categories/:id` | Update category |
| DELETE | `/api/v1/categories/:id` | Delete category |

### Menu Items
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/menu-items` | Create menu item |
| GET | `/api/v1/menu-items` | List all menu items |
| GET | `/api/v1/menu-items?category_id=xxx` | Filter by category |
| GET | `/api/v1/menu-items/:id` | Get menu item by ID |
| PUT | `/api/v1/menu-items/:id` | Update menu item |
| PATCH | `/api/v1/menu-items/:id/availability` | Toggle availability |
| DELETE | `/api/v1/menu-items/:id` | Delete menu item |

### Restaurant Tables
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tables` | Create table |
| GET | `/api/v1/tables` | List all tables |
| GET | `/api/v1/tables/:id` | Get table by ID |
| PUT | `/api/v1/tables/:id` | Update table |
| PATCH | `/api/v1/tables/:id/status` | Update table status |
| DELETE | `/api/v1/tables/:id` | Delete table |

### Orders
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/orders` | Create order |
| GET | `/api/v1/orders` | List all orders |
| GET | `/api/v1/orders?status=open` | Filter by status |
| GET | `/api/v1/orders?order_type=dine_in` | Filter by type |
| GET | `/api/v1/orders/:id` | Get order with items |
| PATCH | `/api/v1/orders/:id/status` | Update order status |
| POST | `/api/v1/orders/:id/items` | Add item to order |
| DELETE | `/api/v1/orders/:id/items/:itemId` | Remove item from order |

### Bills
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/bills` | Generate bill for an order |
| GET | `/api/v1/bills/:id` | Get bill by ID |
| GET | `/api/v1/orders/:id/bill` | Get bill for an order |
| PATCH | `/api/v1/bills/:id/discount` | Apply discount |
| PATCH | `/api/v1/bills/:id/void` | Void (cancel) bill |

### Payments
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/payments` | Record payment |
| GET | `/api/v1/bills/:id/payments` | Get payments for a bill |

---

## Week 5: Authentication & Authorization ✅

### Phase 1: User Registration + Password Hashing ✅
- `internal/domain/user.go` — User domain model
- `internal/dto/auth_dto.go` — Register, Login, AuthResponse DTOs
- `internal/repository/user_repo.go` — User DB operations (create, get by email/ID, get count)
- Password hashing with bcrypt

### Phase 2: JWT Login ✅
- `internal/service/auth_service.go` — Business logic:
  - Register with bcrypt password hashing
  - First user auto-becomes admin (bootstrap)
  - Only admin can register new users after that
  - Login verifies password → returns JWT token (24h expiry)
  - Token contains: user_id, email, role, exp
- `internal/handler/auth_handler.go` — Register, Login, GetMe endpoints

### Phase 3: Auth Middleware ✅
- `internal/middleware/auth.go` — JWT verification:
  - Extracts token from `Authorization: Bearer <token>`
  - Validates signature and expiry
  - Sets user_id, email, role in context
  - OptionalAuthMiddleware for register endpoint (first user has no token)
- Public routes: health, login
- Protected routes: everything else

### Phase 4: Role-Based Access Control ✅
- RoleMiddleware checks user's role against allowed roles
- Permissions:
  - Admin: everything (manage menu, users, void bills, all operations)
  - Cashier: create orders, generate bills, record payments, view all
  - Waiter: create orders, update order/table status, view all
- Waiter can't delete categories, generate bills, or register users
- Cashier can't void bills or manage menu

### Phase 5: Verify + Commit ✅
- Full auth flow tested (first user → admin → registers others)
- Token-based access verified
- Role restrictions verified (waiter blocked from admin actions)

---

## API Endpoints (Current)

### Public (no auth required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check + DB ping |
| POST | `/api/v1/auth/login` | Login → returns JWT token |
| POST | `/api/v1/auth/register` | Register (first user = admin, then admin-only) |

### Protected (token required)
| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/auth/me` | All | Get current user |
| GET | `/api/v1/categories` | All | List categories |
| POST | `/api/v1/categories` | Admin | Create category |
| PUT | `/api/v1/categories/:id` | Admin | Update category |
| DELETE | `/api/v1/categories/:id` | Admin | Delete category |
| GET | `/api/v1/menu-items` | All | List menu items |
| POST | `/api/v1/menu-items` | Admin | Create menu item |
| PUT | `/api/v1/menu-items/:id` | Admin | Update menu item |
| PATCH | `/api/v1/menu-items/:id/availability` | Admin | Toggle availability |
| DELETE | `/api/v1/menu-items/:id` | Admin | Delete menu item |
| GET | `/api/v1/tables` | All | List tables |
| PATCH | `/api/v1/tables/:id/status` | All | Update table status |
| POST | `/api/v1/tables` | Admin | Create table |
| PUT | `/api/v1/tables/:id` | Admin | Update table |
| DELETE | `/api/v1/tables/:id` | Admin | Delete table |
| POST | `/api/v1/orders` | All | Create order |
| GET | `/api/v1/orders` | All | List orders (filter by status/type) |
| GET | `/api/v1/orders/:id` | All | Get order with items |
| PATCH | `/api/v1/orders/:id/status` | All | Update order status |
| POST | `/api/v1/orders/:id/items` | All | Add item to order |
| DELETE | `/api/v1/orders/:id/items/:itemId` | All | Remove item from order |
| POST | `/api/v1/bills` | Admin, Cashier | Generate bill |
| GET | `/api/v1/bills/:id` | All | Get bill |
| GET | `/api/v1/orders/:id/bill` | All | Get bill for order |
| PATCH | `/api/v1/bills/:id/discount` | Admin, Cashier | Apply discount |
| PATCH | `/api/v1/bills/:id/void` | Admin | Void bill |
| POST | `/api/v1/payments` | Admin, Cashier | Record payment |
| GET | `/api/v1/bills/:id/payments` | All | Get payments for bill |

---

## Next: Week 6 — Frontend (React)

- Login page
- POS screen (order creation)
- Menu management (admin)
- Bill view & print
