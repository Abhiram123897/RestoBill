# RestoBill — Progress Tracker

## Project Overview

Restaurant billing system built with Go (Gin), PostgreSQL, following BlueDrop's clean architecture and Tharsis's PostgreSQL error handling patterns.

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
- `internal/logger/logger.go` — Zap structured logging (BlueDrop pattern)
- `internal/errs/errors.go` — AppError with typed codes (BlueDrop pattern)
- `internal/db/db.go` — pgxpool connection + PostgreSQL error helpers (Tharsis pattern)
- `internal/handler/health_handler.go` — Health endpoint pings DB
- `cmd/server/main.go` — Graceful shutdown (BlueDrop pattern)
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
- `internal/repository/interfaces.go` — Repository interfaces (BlueDrop ports pattern)
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

---

## Architecture (BlueDrop Pattern)

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
│   ├── interfaces.go        → Port interfaces (BlueDrop pattern)
│   ├── category_repo.go     → PostgreSQL implementation
│   └── menu_item_repo.go    → PostgreSQL implementation
├── service/                 → Business logic layer
├── handler/                 → HTTP handlers (Gin)
├── router/                  → Route wiring
└── middleware/              → (Week 5 — Auth)
```

---

## Patterns Adopted

| Pattern | Source | Used In |
|---------|--------|---------|
| Domain/DTO separation | BlueDrop | domain/, dto/ |
| Port interfaces (repos) | BlueDrop | repository/interfaces.go |
| Service interface + struct | BlueDrop | service/ |
| Handler struct holds service | BlueDrop | handler/ |
| `*errs.AppError` return type | BlueDrop | All layers |
| Error pass-through (repo→service→handler) | BlueDrop | service/ |
| `AsPgError` + `pgerrcode` | Tharsis | db/db.go, repository/ |
| UUID validation in handler | Tharsis | handler/ |
| Zap structured logger | BlueDrop | logger/ |
| Graceful shutdown | BlueDrop | main.go |

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.25 |
| Framework | Gin |
| Database | PostgreSQL 16 |
| DB Driver | pgx/v5 |
| Migrations | golang-migrate |
| Logger | Zap |
| UUID | google/uuid |
| Error Codes | jackc/pgerrcode |

---

## Next: Week 3 — Order Management API

- Create order (select table, add items)
- Update order (add/remove items)
- Order status flow (open → served → billed → paid)
- KOT (Kitchen Order Ticket) generation
