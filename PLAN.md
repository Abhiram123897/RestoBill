# RestoBill — Complete Mentorship Plan

## The Golden Rule: Backend First, Always

**Why?** The backend is the brain. The frontend is just a face. You can test your entire billing logic with Postman/curl before writing a single line of React. Real companies build APIs first, then multiple frontends consume them (web, mobile, tablet).

---

## Phase 0: Foundation (Day 1-2)

### What to do FIRST before writing any code:

1. **Set up Git repo** with proper `.gitignore`
2. **Design the database schema** on paper/whiteboard
3. **List your API endpoints** (what URLs will exist)
4. **Then** start coding the backend

### Mistakes to Avoid Early:

- ❌ Don't start with the frontend — you'll redesign it 5 times
- ❌ Don't build auth first — build core features, add auth later
- ❌ Don't try microservices — start monolith, split later
- ❌ Don't over-engineer — no caching, no message queues, no Kubernetes yet
- ❌ Don't skip database design — a bad schema ruins everything

---

## Architecture Decision: Monolith First

```
Week 1-8: Monolith (single Go binary + single PostgreSQL database)
Later:    Split into microservices ONLY if needed
```

**Why monolith first?**

- Simpler to develop, debug, and deploy
- One codebase, one database, one deployment
- You can always extract services later
- Even Uber and Netflix started as monoliths

**When to split into microservices (much later):**

- When different teams own different features
- When one part needs to scale independently (e.g., order processing vs reporting)
- When deployments of one feature break another

---

## Tech Stack

| Layer | Choice | Why |
|-------|--------|-----|
| **Frontend** | React (Vite) | Fast, component-based UI, great for interactive POS interfaces |
| **Backend** | Go + Gin | Fast, single binary, great concurrency, Gin is simple and popular |
| **Database** | PostgreSQL | Reliable for financial data, strong relational support |
| **Auth** | JWT | Stateless, simple, works well for REST APIs |
| **Receipts** | PDF generation (gofpdf) | For printable bills |
| **Docs** | Swagger (gin-swagger) | Professional API documentation |

---

## Weekly Roadmap

### Week 1: Database + Project Setup

#### Phase 1: Project Skeleton (Day 1)
- [ ] Initialize Go module
- [ ] Create folder structure (cmd/, internal/, migrations/)
- [ ] Create `.gitignore`, `.env`, `Makefile`
- [ ] `git init` + first commit

#### Phase 2: Docker + Database (Day 2)
- [ ] Write `docker-compose.yml` (PostgreSQL)
- [ ] Run `docker compose up` — database is live
- [ ] Connect to it and verify

#### Phase 3: Database Schema (Day 3)
- [ ] Write SQL migration files for all tables:
  - [ ] `users`
  - [ ] `categories`
  - [ ] `menu_items`
  - [ ] `tables`
  - [ ] `orders`
  - [ ] `order_items`
  - [ ] `bills`
  - [ ] `payments`
- [ ] Run migrations — tables created

#### Phase 4: Server + DB Connection (Day 4)
- [ ] Write `main.go` — starts Gin server
- [ ] Write `config.go` — loads env vars
- [ ] Connect to PostgreSQL using pgx
- [ ] Add `/api/v1/health` endpoint that pings the DB

#### Phase 5: Verify & Commit (Day 5)
- [ ] Test health endpoint with curl/Postman
- [ ] Verify all tables exist in DB
- [ ] Clean up, document, commit to Git

### Week 2: Menu Management API

#### Phase 1: Category Domain + Repository (Day 1)
- [ ] Create `internal/domain/category.go` — domain model
- [ ] Create `internal/dto/category_dto.go` — request/response DTOs
- [ ] Create `internal/repository/interfaces.go` — repository port (interface)
- [ ] Create `internal/repository/category_repo.go` — PostgreSQL implementation

#### Phase 2: Category Service + Handler (Day 2)
- [ ] Create `internal/service/category_service.go` — business logic
- [ ] Create `internal/handler/category_handler.go` — HTTP handlers
- [ ] Wire routes in `internal/router/router.go`
- [ ] Test CRUD with Postman (Create, Read, Update, Delete categories)

#### Phase 3: Menu Item Domain + Repository (Day 3)
- [ ] Create `internal/domain/menu_item.go` — domain model
- [ ] Create `internal/dto/menu_item_dto.go` — request/response DTOs
- [ ] Create `internal/repository/menu_item_repo.go` — PostgreSQL implementation

#### Phase 4: Menu Item Service + Handler (Day 4)
- [ ] Create `internal/service/menu_item_service.go` — business logic
- [ ] Create `internal/handler/menu_item_handler.go` — HTTP handlers
- [ ] Wire routes in `internal/router/router.go`
- [ ] Test CRUD with Postman (Create, Read, Update, Delete menu items)
- [ ] Test filtering menu items by category

#### Phase 5: Seed Data + Verify + Commit (Day 5)
- [ ] Create seed SQL script (sample categories + menu items)
- [ ] Run seed data
- [ ] Verify all endpoints with Postman
- [ ] Commit to Git: `feat: complete week 2 - menu management API`

### Week 3: Order Management API

#### Phase 1: Restaurant Table Domain + Repository (Day 1)
- [ ] Create `internal/domain/restaurant_table.go` — domain model
- [ ] Create `internal/dto/restaurant_table_dto.go` — request/response DTOs
- [ ] Add `RestaurantTableRepository` interface to `interfaces.go`
- [ ] Create `internal/repository/restaurant_table_repo.go` — PostgreSQL implementation
- [ ] Create `internal/service/restaurant_table_service.go` — business logic
- [ ] Create `internal/handler/restaurant_table_handler.go` — HTTP handlers
- [ ] Wire routes and test CRUD + status update with Postman

#### Phase 2: Order Domain + Repository (Day 2)
- [ ] Create `internal/domain/order.go` — domain model (Order + OrderItem)
- [ ] Create `internal/dto/order_dto.go` — request/response DTOs
- [ ] Add `OrderRepository` interface to `interfaces.go`
- [ ] Create `internal/repository/order_repo.go` — PostgreSQL implementation

#### Phase 3: Order Service + Handler (Day 3)
- [ ] Create `internal/service/order_service.go` — business logic (create order, add items, update status)
- [ ] Create `internal/handler/order_handler.go` — HTTP handlers
- [ ] Wire routes in router
- [ ] Test create order + add items with Postman

#### Phase 4: Order Status Flow + Item Management (Day 4)
- [ ] Implement order status transitions (open → served → billed → paid)
- [ ] Add/remove items from an existing order
- [ ] Update table status when order is created/completed
- [ ] Test full order lifecycle with Postman

#### Phase 5: Verify + Commit (Day 5)
- [ ] Test all order endpoints with Postman
- [ ] Verify table status changes correctly
- [ ] Update PROGRESS.md
- [ ] Commit to Git: `feat: complete week 3 - order management API`

### Week 4: Billing & Payments API

#### Phase 1: Bill Domain + Repository (Day 1)
- [ ] Create `internal/domain/bill.go` — domain model
- [ ] Create `internal/dto/bill_dto.go` — request/response DTOs
- [ ] Add `BillRepository` interface to `interfaces.go`
- [ ] Create `internal/repository/bill_repo.go` — PostgreSQL implementation

#### Phase 2: Bill Service + Handler (Day 2)
- [ ] Create `internal/service/bill_service.go` — business logic:
  - Generate bill from order (calculate subtotal from order items)
  - Tax calculation (configurable tax rate)
  - Discount application (amount or percentage)
  - Bill status management (pending/paid/void)
- [ ] Create `internal/handler/bill_handler.go` — HTTP handlers
- [ ] Wire routes and test bill generation with Postman

#### Phase 3: Payment Domain + Repository (Day 3)
- [ ] Create `internal/domain/payment.go` — domain model
- [ ] Create `internal/dto/payment_dto.go` — request/response DTOs
- [ ] Add `PaymentRepository` interface to `interfaces.go`
- [ ] Create `internal/repository/payment_repo.go` — PostgreSQL implementation

#### Phase 4: Payment Service + Handler (Day 4)
- [ ] Create `internal/service/payment_service.go` — business logic:
  - Record payment against a bill
  - Support multiple payment methods (cash/card/upi)
  - Auto-mark bill as paid when full amount received
  - Auto-update order status to paid
- [ ] Create `internal/handler/payment_handler.go` — HTTP handlers
- [ ] Wire routes and test full billing flow with Postman

#### Phase 5: Verify + Commit (Day 5)
- [ ] Test complete flow: order → bill → payment → table freed
- [ ] Test edge cases (bill for unpaid order, double payment, void bill)
- [ ] Update PROGRESS.md
- [ ] Commit to Git: `feat: complete week 4 - billing and payments API`

### Week 5: Authentication & Authorization

- [ ] JWT-based login
- [ ] User roles (admin, cashier, waiter)
- [ ] Protected routes
- [ ] Role-based access control

### Week 6: Frontend — Core Pages

- [ ] Login page
- [ ] POS screen (order creation)
- [ ] Menu management (admin)
- [ ] Bill view & print

### Week 7: Frontend — Advanced

- [ ] Table management UI
- [ ] Order status dashboard
- [ ] Reports page (daily sales)
- [ ] Receipt PDF generation

### Week 8: Polish & Deploy

- [ ] Docker setup
- [ ] Error handling & validation
- [ ] API documentation (Swagger)
- [ ] Deploy to cloud
- [ ] README & documentation

---

## Project Milestones

| Milestone | Deliverable | Status |
|-----------|-------------|--------|
| M1 | Project setup + DB schema + migrations running | ⬜ |
| M2 | Menu CRUD API working (tested via Postman) | ⬜ |
| M3 | Order creation + item management working | ⬜ |
| M4 | Bill generation + payment recording working | ⬜ |
| M5 | JWT auth + role-based access working | ⬜ |
| M6 | Frontend POS screen + menu management | ⬜ |
| M7 | Reports + receipts + table management | ⬜ |
| M8 | Dockerized + deployed + documented | ⬜ |

---

## Database Schema Design

```
┌─────────────────┐     ┌─────────────────┐
│   categories    │     │   menu_items    │
├─────────────────┤     ├─────────────────┤
│ id (PK)         │────<│ id (PK)         │
│ name            │     │ category_id (FK)│
│ display_order   │     │ name            │
│ created_at      │     │ price           │
│ updated_at      │     │ is_available    │
└─────────────────┘     │ created_at      │
                        │ updated_at      │
                        └─────────────────┘

┌─────────────────┐     ┌─────────────────┐
│    tables       │     │    orders       │
├─────────────────┤     ├─────────────────┤
│ id (PK)         │────<│ id (PK)         │
│ number          │     │ table_id (FK)   │
│ capacity        │     │ status          │ (open/served/billed/paid)
│ status          │     │ order_type      │ (dine_in/takeaway/delivery)
│ created_at      │     │ created_by (FK) │
└─────────────────┘     │ created_at      │
                        │ updated_at      │
                        └────────┬────────┘
                                 │
                        ┌────────┴────────┐
                        │  order_items    │
                        ├─────────────────┤
                        │ id (PK)         │
                        │ order_id (FK)   │
                        │ menu_item_id(FK)│
                        │ quantity        │
                        │ unit_price      │ (snapshot at order time)
                        │ notes           │ ("extra spicy")
                        │ created_at      │
                        └─────────────────┘

┌─────────────────┐     ┌─────────────────┐
│    bills        │     │   payments      │
├─────────────────┤     ├─────────────────┤
│ id (PK)         │────<│ id (PK)         │
│ order_id (FK)   │     │ bill_id (FK)    │
│ subtotal        │     │ amount          │
│ tax_amount      │     │ method          │ (cash/card/upi)
│ tax_rate        │     │ reference       │ (transaction ID)
│ discount_amount │     │ created_at      │
│ discount_reason │     └─────────────────┘
│ total           │
│ status          │ (pending/paid/void)
│ generated_by(FK)│
│ created_at      │
└─────────────────┘

┌─────────────────┐
│    users        │
├─────────────────┤
│ id (PK)         │
│ name            │
│ email           │
│ password_hash   │
│ role            │ (admin/cashier/waiter)
│ is_active       │
│ created_at      │
│ updated_at      │
└─────────────────┘
```

### Key Design Decisions:

- `unit_price` in `order_items` — snapshot the price at order time (menu prices can change later)
- `status` fields — track lifecycle of orders and bills
- `generated_by` / `created_by` — audit trail (who did what)
- Separate `bills` from `orders` — an order can exist without a bill (still being served)

---

## Backend Folder Structure

```
restobill/
├── cmd/
│   └── server/
│       └── main.go              ← Entry point
├── internal/
│   ├── config/
│   │   └── config.go           ← Load env vars
│   ├── domain/
│   │   ├── category.go         ← Domain models
│   │   ├── menu_item.go
│   │   ├── order.go
│   │   ├── bill.go
│   │   ├── payment.go
│   │   └── user.go
│   ├── dto/
│   │   ├── request.go          ← API request shapes
│   │   └── response.go         ← API response shapes
│   ├── handler/
│   │   ├── category_handler.go ← HTTP handlers
│   │   ├── menu_handler.go
│   │   ├── order_handler.go
│   │   ├── bill_handler.go
│   │   └── auth_handler.go
│   ├── middleware/
│   │   ├── auth.go             ← JWT verification
│   │   └── audit.go            ← Request logging
│   ├── repository/
│   │   ├── interfaces.go       ← Repository contracts
│   │   ├── category_repo.go    ← PostgreSQL implementations
│   │   ├── menu_repo.go
│   │   ├── order_repo.go
│   │   ├── bill_repo.go
│   │   └── user_repo.go
│   ├── service/
│   │   ├── menu_service.go     ← Business logic
│   │   ├── order_service.go
│   │   ├── bill_service.go
│   │   └── auth_service.go
│   └── router/
│       └── router.go           ← Route definitions
├── migrations/
│   ├── 001_create_users.up.sql
│   ├── 001_create_users.down.sql
│   ├── 002_create_categories.up.sql
│   └── ...
├── .env
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## REST API Design

### Menu Management

```
GET    /api/v1/categories              ← List all categories
POST   /api/v1/categories              ← Create category (admin)
PUT    /api/v1/categories/:id          ← Update category (admin)
DELETE /api/v1/categories/:id          ← Delete category (admin)

GET    /api/v1/menu-items              ← List all items (filter by category)
POST   /api/v1/menu-items              ← Create item (admin)
PUT    /api/v1/menu-items/:id          ← Update item (admin)
PATCH  /api/v1/menu-items/:id/availability ← Toggle available/unavailable
```

### Order Management

```
POST   /api/v1/orders                  ← Create new order
GET    /api/v1/orders                  ← List orders (filter by status/table)
GET    /api/v1/orders/:id              ← Get order details
PUT    /api/v1/orders/:id/items        ← Add/update items in order
PATCH  /api/v1/orders/:id/status       ← Update order status
DELETE /api/v1/orders/:id/items/:itemId ← Remove item from order
```

### Billing

```
POST   /api/v1/bills                   ← Generate bill for an order
GET    /api/v1/bills/:id               ← Get bill details
PATCH  /api/v1/bills/:id/discount      ← Apply discount
POST   /api/v1/bills/:id/payments      ← Record payment
GET    /api/v1/bills/:id/receipt        ← Get printable receipt
```

### Tables

```
GET    /api/v1/tables                  ← List all tables with status
PATCH  /api/v1/tables/:id/status       ← Update table status
```

### Auth

```
POST   /api/v1/auth/login              ← Login, get JWT
POST   /api/v1/auth/register           ← Create user (admin only)
GET    /api/v1/auth/me                 ← Get current user info
```

### Reports

```
GET    /api/v1/reports/daily-sales     ← Today's sales summary
GET    /api/v1/reports/item-sales      ← Item-wise sales
```

---

## Authentication Flow (JWT)

```
1. User sends POST /api/v1/auth/login with {email, password}
2. Server verifies password hash (bcrypt)
3. Server generates JWT token with {user_id, role, exp}
4. Client stores token in localStorage
5. Client sends token in header: Authorization: Bearer <token>
6. Middleware verifies token on every protected request
7. Middleware checks role for admin-only routes
```

### Role-Based Access:

| Role | Can Do |
|------|--------|
| **admin** | Everything — manage menu, users, view reports, void bills |
| **cashier** | Create orders, generate bills, record payments |
| **waiter** | Create orders, update order status, view menu |

---

## Frontend Pages (React)

| Page | Who Uses It | Priority |
|------|-------------|----------|
| Login | Everyone | Week 6 |
| POS Screen (order creation) | Cashier/Waiter | Week 6 |
| Menu Management | Admin | Week 6 |
| Bill/Receipt View | Cashier | Week 6 |
| Table Map | Waiter | Week 7 |
| Order Status Dashboard | Waiter/Kitchen | Week 7 |
| Reports Dashboard | Admin/Manager | Week 7 |
| User Management | Admin | Week 7 |

---

## Docker Setup

```yaml
# docker-compose.yml
services:
  db:
    image: postgres:16
    environment:
      POSTGRES_DB: restobill
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DB_HOST: db
      DB_PORT: 5432
      DB_NAME: restobill
      DB_USER: postgres
      DB_PASSWORD: postgres
      JWT_SECRET: your-secret-key
    depends_on:
      - db

volumes:
  pgdata:
```

```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
COPY migrations ./migrations
EXPOSE 8080
CMD ["./server"]
```

---

## GitHub Workflow & Best Practices

### Branch Strategy

```
main          ← production-ready code
└── develop   ← integration branch
    ├── feature/menu-crud
    ├── feature/order-management
    ├── feature/billing
    └── fix/tax-calculation
```

### Commit Message Convention

```
feat: add menu item CRUD endpoints
fix: correct tax calculation for takeaway orders
docs: add API documentation
refactor: extract bill calculation to service layer
test: add order creation unit tests
chore: update dependencies
```

### PR Workflow

1. Create feature branch from `develop`
2. Make small, focused commits
3. Push and create PR
4. Write description of what changed and why
5. Merge to `develop` → test → merge to `main`

### GitHub Actions CI (`.github/workflows/ci.yml`)

```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_DB: test_db
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
        ports:
          - 5432:5432
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - run: go test ./...
      - run: go build ./cmd/server
```

---

## How Real Restaurant Billing Systems Work

1. **Waiter** takes order on tablet → creates order linked to table
2. **KOT prints** in kitchen automatically
3. Kitchen prepares food, marks items as ready
4. Waiter serves food, updates order status
5. Customer asks for bill → **cashier generates bill**
6. System calculates: subtotal + tax - discounts = total
7. Customer pays (cash/card/UPI) → payment recorded
8. **Receipt prints** → table freed for next customer
9. End of day → manager checks **daily sales report**

**Key insight:** The order and the bill are separate events. An order can have items added over time (customer orders dessert later). The bill is generated only when the customer asks for it.

---

## How to Make This Resume-Worthy

1. **Clean code** — proper folder structure, interfaces, separation of concerns
2. **Tests** — unit tests for service layer, integration tests for API
3. **Docker** — one command to run everything (`docker compose up`)
4. **Swagger docs** — professional API documentation
5. **CI/CD** — GitHub Actions for build + test on every PR
6. **README** — screenshots, architecture diagram, setup instructions
7. **Real features** — don't just CRUD; add tax logic, discounts, reports, receipts
8. **Security** — JWT auth, role-based access, input validation, SQL injection prevention
9. **Production touches** — graceful shutdown, health checks, structured logging, error handling
10. **Git history** — clean commits, PRs, branches (shows you work professionally)

---

## Recommended Libraries

| Purpose | Library |
|---------|---------|
| Router | `github.com/gin-gonic/gin` |
| Database driver | `github.com/jackc/pgx/v5` |
| Migrations | `github.com/golang-migrate/migrate/v4` |
| JWT | `github.com/golang-jwt/jwt/v5` |
| Validation | `github.com/go-playground/validator/v10` |
| Config | `github.com/joho/godotenv` |
| Logging | `go.uber.org/zap` |
| UUID | `github.com/google/uuid` |
| PDF receipts | `github.com/jung-kurt/gofpdf` |
| Swagger | `github.com/swaggo/gin-swagger` |
| Testing | `github.com/stretchr/testify` |
| Password hashing | `golang.org/x/crypto/bcrypt` |

---

## Scaling to Microservices (Future)

When the monolith grows too large, split along domain boundaries:

```
Monolith → Split into:
├── menu-service        (menu CRUD, categories)
├── order-service       (order lifecycle, KOT)
├── billing-service     (bill generation, payments)
├── auth-service        (login, JWT, user management)
└── reporting-service   (sales reports, analytics)
```

Each gets its own database, communicates via REST or message queues (RabbitMQ/Kafka).

**Don't do this until you have a working monolith.**

---

## Deployment Options

| Option | Best For | Cost | Complexity |
|--------|----------|------|------------|
| **Local network** | Single restaurant, no internet needed | Free | Low |
| **Railway / Render** | Quick deploy, managed infrastructure | Free tier → ~$10-20/mo | Low |
| **VPS (DigitalOcean)** | Good balance of control and simplicity | ~$5-12/mo | Medium |
| **AWS** | Production, scalable, full control | ~$15-50/mo | Medium-High |

**Recommended path:** Build locally → Deploy to Railway/Render for demo → Move to AWS for production.

---

## Insights from Analyzed Projects

### From Tharsis (Go + Chi + PostgreSQL):

- ResourceMetadata pattern (ID, CreatedAt, UpdatedAt, Version) for all models
- Optimistic locking for concurrent order modifications
- Transaction handling via context
- Database migration system (golang-migrate)

### From BlueDrop (Go + Gin + DynamoDB):

- Clean architecture (domain/dto/ports/service)
- Audit logging middleware (track who did what)
- Channel-based async processing
- Custom error handling (`AppError`)
- Domain ↔ DTO conversion pattern
- Graceful shutdown

---

## Quick Reference: Build Order

```
1. docker-compose.yml (PostgreSQL)     ← Get DB running
2. migrations/ (schema)                ← Create tables
3. internal/domain/ (models)           ← Define data structures
4. internal/repository/ (DB access)    ← Talk to database
5. internal/service/ (business logic)  ← Core logic
6. internal/handler/ (HTTP handlers)   ← API endpoints
7. internal/router/ (routes)           ← Wire it all together
8. cmd/server/main.go (entry point)    ← Start the server
9. Test with Postman                   ← Verify everything works
10. Add auth middleware                ← Secure the API
11. Build frontend                     ← React UI
12. Dockerize everything               ← One-command setup
13. Deploy                             ← Go live
```

---

## Next Step

Say **"let's start"** to begin building the RestoBill project from step 1.
