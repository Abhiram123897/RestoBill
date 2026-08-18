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

---

## Next: Week 4 — Billing & Payments API

- Generate bill from order
- Tax calculation
- Discount application
- Payment recording
- Bill finalization
