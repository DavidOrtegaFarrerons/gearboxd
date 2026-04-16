> 🚗 **Live API:** [`https://api.gearboxd.davidortegafarrerons.com/v1/healthcheck`](https://api.gearboxd.davidortegafarrerons.com/v1/healthcheck)

> 🏎️ **Live Website:** [`https://gearboxd.davidortegafarrerons.com`](https://gearboxd.davidortegafarrerons.com)
# Gearboxd

A production-grade REST API built in Go, a Letterboxd-style car discovery platform. Users can browse a catalogue of iconic cars, register accounts, authenticate, and interact with the collection. Built as a portfolio project following the patterns in *Let's Go Further* by Alex Edwards.

---

## Tech Stack

- **Language:** Go (stdlib `net/http`, no framework)
- **Database:** PostgreSQL with `golang-migrate`
- **Driver:** `pq` (raw SQL, no ORM)
- **Routing:** `httprouter`
- **Search:** PostgreSQL `ILIKE` pattern matching
- **Decimal precision:** `shopspring/decimal`
- **Local dev:** Docker & Docker Compose
- **API docs:** Bruno

---

## Project Structure

```
.
├── cmd/
│   ├── api/
│   │   ├── main.go           # Entry point, config flags, application setup
│   │   ├── server.go         # HTTP server with graceful shutdown
│   │   ├── routes.go         # Route registration
│   │   ├── middleware.go     # Rate limiting, authentication, CORS
│   │   ├── cars.go           # Car handlers
│   │   ├── users.go          # User registration & activation handlers
│   │   ├── tokens.go         # Authentication token handlers
│   │   ├── healthcheck.go    # Healthcheck handler
│   │   ├── helpers.go        # JSON read/write helpers, pagination
│   │   ├── errors.go         # Centralised error responses
│   │   └── context.go        # Request context helpers
│   └── seed/
│       └── main.go           # Database seeder (30+ cars)
└── internal/
    ├── data/
    │   ├── models.go         # Models struct (dependency injection)
    │   ├── cars.go           # CarStore interface & PostgresCarStore
    │   ├── users.go          # User store & SQL queries
    │   ├── tokens.go         # Token store & SQL queries
    │   ├── permissions.go    # Permission scopes
    │   ├── filters.go        # Filtering, sorting, pagination logic
    │   └── metadata.go       # Pagination metadata envelope
    ├── validator/
    │   └── validator.go      # Input validation helpers
    ├── mailer/
    │   └── mailer.go         # Email sending (background goroutines)
    └── assert/
        └── assert.go         # Test assertion helpers
```

---

## API Endpoints

### Healthcheck

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| `GET` | `/v1/healthcheck` | Server status and environment info | None |

### Cars

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| `GET` | `/v1/cars` | List cars with filtering, sorting & pagination | None |
| `POST` | `/v1/cars` | Create a new car | `cars:write` |
| `GET` | `/v1/cars/:id` | Get a single car by ID | None |
| `PATCH` | `/v1/cars/:id` | Partially update a car | `cars:write` |
| `DELETE` | `/v1/cars/:id` | Delete a car | `cars:write` |

**Query parameters for `GET /v1/cars`:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `make` | string | Filter by manufacturer |
| `model` | string | Filter by model name |
| `year` | integer | Filter by year |
| `gearbox` | string | `manual` or `automatic` |
| `drivetrain` | string | e.g. `RWD`, `FWD`, `AWD` |
| `fuel` | string | e.g. `gas`, `diesel`, `hybrid`, `electric` |
| `horsepower_min` | integer | Minimum horsepower |
| `horsepower_max` | integer | Maximum horsepower |
| `price_min` | integer | Minimum price (new) |
| `price_max` | integer | Maximum price (new) |
| `page` | integer | Page number (default: 1) |
| `page_size` | integer | Results per page (default: 20, max: 100) |
| `sort` | string | Sort field; prefix with `-` for descending (e.g. `-year`) |

### Users

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| `POST` | `/v1/users` | Register a new user | None |
| `PUT` | `/v1/users/activated` | Activate account with token | None |

### Tokens

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| `POST` | `/v1/tokens/authentication` | Get a Bearer authentication token | None |

---

## Running Locally

**Prerequisites:** Docker & Make

```bash
# Copy and fill in environment variables
cp .env.example .env

# Build and start the server
make build

# (Optional) Seed the database (requires the container to be running)
make seed
```

---

## Testing

Tests are written alongside implementation.

```bash
make test
```

**Test patterns used:**

- Table-driven tests throughout for validators, helpers, and handlers
- Store interface pattern (`CarStore`, `MockCarStore`) — no real database required for handler tests
- `httptest.NewRecorder` for middleware and handler tests
- `httprouter` params injected via `context.WithValue` in handler tests
- Custom `internal/assert` package for clean test assertions

---

## Key Design Decisions

**No ORM.** Raw SQL keeps queries explicit and avoids hidden N+1 problems. `database/sql` with `pgx` is idiomatic for production Go services.

**Store interface pattern.** Each data layer is defined as an interface (`CarStore`, `PostgresCarStore`, `MockCarStore`). Handlers depend on the interface, making them fully testable without a real database.

**Sentinel values for optional range filters.** `(column >= $n OR $n = 0)` short-circuits cleanly when no value is provided, keeping query construction simple without dynamic SQL.

**Optimistic concurrency via version field.** Car rows carry a `version` integer. `UPDATE ... WHERE id = $1 AND version = $2` prevents lost updates from concurrent admin edits — if the version has changed, the update affects zero rows and the handler returns a `409 Conflict`.

**Stateful Bearer tokens.** Auth tokens are stored in the database with an expiry timestamp. This enables server-side revocation, a capability JWT cannot offer without additional infrastructure.

**Separate seed binary.** The seeder at `cmd/seed` is a standalone binary that only requires a database connection. It uses a wipe-and-reinsert strategy wrapped in a transaction for idempotent, atomic seeding.

---

## Permissions

| Scope | Who has it | What it allows |
|-------|-----------|----------------|
| `cars:read` | All activated users | Read the car catalogue |
| `cars:write` | Admin users | Create, update, delete cars |

---

## License

MIT