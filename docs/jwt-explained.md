# JWT Authentication - Complete Explanation

## Without JWT (No Security)

Imagine a simple API with no authentication:

```
Client                              Server
  │                                    │
  │── GET /api/v1/categories ─────────►│
  │                                    │  Anyone can access!
  │◄── [{id: 1, name: "Starters"}] ──│
  │                                    │
  │── DELETE /api/v1/categories/1 ────►│
  │                                    │  Anyone can delete! 😱
  │◄── "deleted" ─────────────────────│
```

**Problem:** Anyone can do anything. A random person could delete all your menu items.

---

## With JWT (Secured)

```
Client                              Server
  │                                    │
  │── DELETE /api/v1/categories/1 ────►│
  │   (no token)                       │
  │◄── 401 "unauthorized" ───────────│  Blocked! ✅
  │                                    │
  │── DELETE /api/v1/categories/1 ────►│
  │   Authorization: Bearer eyJ...     │  Token present → verify → check role
  │◄── "deleted" ─────────────────────│  Admin → allowed ✅
```

---

## The Complete Flow (Step by Step)

### Step 1: User Logs In

**Request:**
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "abhiram@restaurant.com",
  "password": "mypassword123"
}
```

**What server does:**
1. Finds user in database by email
2. Compares password with stored hash
3. If valid → creates a JWT token
4. Returns the token to the client

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWJjLTEyMyIsImVtYWlsIjoiYWJoaXJhbUByZXN0YXVyYW50LmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcyNDYwMDAwMH0.xyz_signature",
  "user": {
    "id": "abc-123",
    "name": "Abhiram",
    "email": "abhiram@restaurant.com",
    "role": "admin"
  }
}
```

**The token decoded:**
```
HEADER:    {"alg": "HS256", "typ": "JWT"}
PAYLOAD:   {"user_id": "abc-123", "email": "abhiram@restaurant.com", "role": "admin", "exp": 1724600000}
SIGNATURE: HMACSHA256(header + payload, "restobill-dev-secret-change-in-production")
```

---

### Step 2: Client Stores the Token

The client (browser, Postman, mobile app) saves this token in memory/localStorage.

---

### Step 3: Client Makes Authenticated Requests

Every request to a protected endpoint must include the token in the **Authorization header**:

```
GET /api/v1/categories
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWJjLTEyMyIsImVtYWlsIjoiYWJoaXJhbUByZXN0YXVyYW50LmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcyNDYwMDAwMH0.xyz_signature
```

**Format:** `Authorization: Bearer <your_token_here>`

---

### Step 4: Server Receives the Request

Here's exactly what happens in your code when a request hits a protected route:

```
Request arrives at: POST /api/v1/categories
                         │
                         ▼
              ┌─────────────────────┐
              │  AuthMiddleware      │ ← Runs FIRST
              │                     │
              │  1. Read header     │
              │  2. Extract token   │
              │  3. Verify token    │
              │  4. Set user info   │
              └──────────┬──────────┘
                         │ (if valid)
                         ▼
              ┌─────────────────────┐
              │  RoleMiddleware      │ ← Runs SECOND
              │                     │
              │  1. Get role from   │
              │     context         │
              │  2. Check if role   │
              │     is allowed      │
              └──────────┬──────────┘
                         │ (if allowed)
                         ▼
              ┌─────────────────────┐
              │  Handler            │ ← Runs LAST
              │                     │
              │  Process the actual │
              │  request            │
              └─────────────────────┘
```

---

## How the Router Connects Everything

**File: `internal/router/router.go`**

```go
// PUBLIC routes - anyone can access (no middleware)
api.GET("/health", healthHandler.GetHealth)
api.POST("/auth/login", authHandler.Login)

// PROTECTED routes - token required
protected := api.Group("")
protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))  // ← All routes below need a valid token
{
    // Any authenticated user can read
    protected.GET("/categories", categoryHandler.GetAllCategories)

    // Only admin can create (extra RoleMiddleware check)
    protected.POST("/categories", middleware.RoleMiddleware("admin"), categoryHandler.CreateCategory)
}
```

**What `protected.Use(middleware.AuthMiddleware(...))` means:**

Every request to routes inside this group will FIRST pass through AuthMiddleware. If the token is invalid or missing → 401 response, request never reaches the handler.

**What `middleware.RoleMiddleware("admin")` means:**

After AuthMiddleware passes, this checks if the user's role (from the token) matches "admin". If not → 403 response.

---

## The Middleware Chain (Visual)

### Example 1: GET /api/v1/categories (any authenticated user)

```
Request: GET /categories
Header:  Authorization: Bearer eyJ...

Route definition:
  protected.GET("/categories", categoryHandler.GetAllCategories)
            │                        │
            │                        └── Handler (final destination)
            │
            └── protected group has: AuthMiddleware

Execution order:
  1. AuthMiddleware → verifies token → extracts user_id, role → ✅ passes
  2. categoryHandler.GetAllCategories → returns categories
```

### Example 2: POST /api/v1/categories (admin only)

```
Request: POST /categories
Header:  Authorization: Bearer eyJ... (waiter's token, role="waiter")

Route definition:
  protected.POST("/categories", middleware.RoleMiddleware("admin"), categoryHandler.CreateCategory)
             │                       │                                    │
             │                       │                                    └── Handler
             │                       └── RoleMiddleware (checks role)
             │
             └── protected group has: AuthMiddleware

Execution order:
  1. AuthMiddleware → verifies token → role="waiter" → ✅ passes
  2. RoleMiddleware("admin") → is "waiter" == "admin"? NO → ❌ 403 Forbidden
  3. Handler never runs
```

### Example 3: POST /api/v1/categories (admin token)

```
Request: POST /categories
Header:  Authorization: Bearer eyJ... (admin's token, role="admin")

Execution order:
  1. AuthMiddleware → verifies token → role="admin" → ✅ passes
  2. RoleMiddleware("admin") → is "admin" == "admin"? YES → ✅ passes
  3. categoryHandler.CreateCategory → creates the category → returns 201
```

### Example 4: No token at all

```
Request: GET /categories
Header:  (no Authorization header)

Execution order:
  1. AuthMiddleware → no header found → ❌ 401 "authorization header required"
  2. Nothing else runs
```

---

## The JWT Secret Key Explained

### What is it?

```env
JWT_SECRET=restobill-dev-secret-change-in-production
```

It's a **password that only the server knows**. It's used for two things:

### 1. Signing (creating tokens)

When a user logs in, the server creates a token and **signs** it with the secret:

```go
func (s *authService) generateToken(user *domain.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    user.Role,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))  // ← SECRET used here to SIGN
}
```

The signature = `HMACSHA256(header + "." + payload, SECRET_KEY)`

### 2. Verifying (checking tokens)

When a request comes in, the server **verifies** the token with the same secret:

```go
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(jwtSecret), nil  // ← Same SECRET used here to VERIFY
})
```

### Why is it important?

Without the secret, you can't create valid signatures. If someone tries to:
1. Create a fake token → signature won't match → server rejects it
2. Modify the payload (e.g., change role from "waiter" to "admin") → signature won't match → server rejects it

**Analogy:** The secret is like the stamp of a notary. Only the notary has the real stamp. If someone forges a document without the stamp, it's immediately recognized as fake.

### What happens if the secret is leaked?

Anyone with the secret can:
- Create tokens for any user
- Create tokens with any role (admin)
- Bypass all security

That's why the `.env` file says `change-in-production` — in production, use a long, random string.

---

## Without JWT vs With JWT - Side by Side

### Without JWT:

```go
// router.go - No protection
r.GET("/categories", categoryHandler.GetAllCategories)
r.POST("/categories", categoryHandler.CreateCategory)   // Anyone can create!
r.DELETE("/categories/:id", categoryHandler.DeleteCategory)  // Anyone can delete!
```

```
Request:
  POST /categories
  Body: {"name": "Starters"}

  → Goes directly to handler
  → Category created
  → No way to know WHO did it or IF they should be allowed
```

### With JWT:

```go
// router.go - Protected
protected := api.Group("")
protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
{
    protected.GET("/categories", categoryHandler.GetAllCategories)
    protected.POST("/categories", middleware.RoleMiddleware("admin"), categoryHandler.CreateCategory)
    protected.DELETE("/categories/:id", middleware.RoleMiddleware("admin"), categoryHandler.DeleteCategory)
}
```

```
Request:
  POST /categories
  Authorization: Bearer eyJ...
  Body: {"name": "Starters"}

  → AuthMiddleware: Is token valid? Who is this?
  → RoleMiddleware: Is this user an admin?
  → Only then → handler runs
```

---

## How Headers Work in HTTP Requests

An HTTP request has 3 parts:

```
┌─────────────────────────────────────────┐
│ REQUEST LINE                             │
│ POST /api/v1/categories HTTP/1.1         │
├─────────────────────────────────────────┤
│ HEADERS (metadata about the request)     │
│                                          │
│ Host: localhost:8080                      │
│ Content-Type: application/json           │
│ Authorization: Bearer eyJhbGciOi...      │  ← Token goes here
│ Accept: application/json                 │
├─────────────────────────────────────────┤
│ BODY (the actual data)                   │
│                                          │
│ {"name": "Starters"}                     │
└─────────────────────────────────────────┘
```

**Headers are like the envelope of a letter:**
- `Host` → address (where to send)
- `Content-Type` → what's inside (JSON, XML, etc.)
- `Authorization` → your ID card (proves who you are)
- `Accept` → what format you want back

**The server reads the Authorization header:**
```go
authHeader := c.GetHeader("Authorization")
// Returns: "Bearer eyJhbGciOiJIUzI1NiJ9..."

parts := strings.Split(authHeader, " ")
// parts[0] = "Bearer"
// parts[1] = "eyJhbGciOiJIUzI1NiJ9..."  ← This is the token
```

---

## Testing with Postman / curl

### Login (get a token):

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "abhiram@restaurant.com", "password": "mypassword123"}'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYWJjLTEyMyIsInJvbGUiOiJhZG1pbiJ9.signature",
  "user": {"name": "Abhiram", "role": "admin"}
}
```

### Use the token:

```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYWJjLTEyMyIsInJvbGUiOiJhZG1pbiJ9.signature" \
  -d '{"name": "Starters"}'
```

### Without token (will fail):

```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Starters"}'
```

Response:
```json
{"message": "authorization header required", "type": "Authentication error"}
```

---

## Summary: Where the Secret Key is Used

```
.env file
  │
  │  JWT_SECRET=restobill-dev-secret-change-in-production
  │
  ▼
config.go (loads it)
  │
  │  cfg.JWTSecret = getEnv("JWT_SECRET", "")
  │
  ▼
router.go (passes it to middleware and service)
  │
  ├──► AuthService (uses secret to CREATE tokens)
  │      │
  │      └── generateToken():
  │            token.SignedString([]byte(s.jwtSecret))
  │
  └──► AuthMiddleware (uses secret to VERIFY tokens)
         │
         └── jwt.Parse(tokenString, func(...) {
               return []byte(jwtSecret), nil
             })
```

**Same secret, two uses:**
1. **Sign** (create) → when user logs in
2. **Verify** (check) → when user makes a request

If the signing secret ≠ verifying secret → token rejected. That's why both must use the same value from `.env`.

---

## Real-World Analogy

Think of a **nightclub**:

1. **Login** = Showing your ID at the entrance → bouncer gives you a **wristband** (token)
2. **JWT Secret** = The special ink/pattern on the wristband that only the club can make
3. **Authorization Header** = Showing your wristband when ordering at the bar
4. **AuthMiddleware** = Bartender checking if your wristband is real (not fake)
5. **RoleMiddleware** = VIP area bouncer checking if your wristband says "VIP" (admin)

- **No wristband** → "You can't be here" (401)
- **Fake wristband** → "This is not real" (401)
- **Real wristband, not VIP** → "You can't enter VIP area" (403)
- **Real VIP wristband** → "Welcome!" (200)
