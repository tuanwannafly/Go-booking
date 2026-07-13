# GoBooking — Flight/Hotel Reservation Core Service

A production-ready Go service demonstrating **concurrency-safe booking**, **idempotency**, and **distributed locking** — built for a Backend Golang Intern portfolio.

##  Project Goals

- **Anti-overbooking**: Pessimistic lock (Redis + `SELECT FOR UPDATE`) for flight seats; Optimistic lock (version column) for hotel rooms
- **Idempotency**: `Idempotency-Key` header prevents double-charge on retry
- **Concurrency proof**: 100-goroutine test with `go test -race` proving no overbooking
- **Clean architecture**: Handler → Service → Repository → DB/Redis
- **Observable**: Structured logging, health checks, OpenAPI docs

##  Architecture

```
HTTP Handler → Service (Use Case) → Repository → PostgreSQL / Redis
```

**Locking Strategy** (interview talking point):
| Resource | Strategy | Why |
|----------|----------|-----|
| Flight Seat | Pessimistic (Redis lock + `SELECT FOR UPDATE`) | Burst contention during flash sales |
| Hotel Room | Optimistic (`version` column + retry) | Contention spread over days, lower peak |

##  Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.22+ |
| HTTP | Gin |
| Database | PostgreSQL 16 (pgx/v5) |
| Cache/Lock | Redis 7 + Redsync |
| Migration | golang-migrate |
| Test | testing + testify + `-race` |
| Container | Docker Compose |
| CI | GitHub Actions |
| Docs | Swaggo (OpenAPI 3) |

##  Project Structure

```
gobooking/
├── cmd/api/main.go                 # Entry point
├── internal/
│   ├── config/                     # Viper config
│   ├── domain/                     # Entities + interfaces
│   ├── service/                    # Business logic
│   ├── repository/
│   │   ├── postgres/               # PG implementations
│   │   └── redis/                  # Redis client + Redsync
│   ├── http/
│   │   ├── handler/                # Gin handlers
│   │   ├── middleware/             # Logging, rate-limit, idempotency, auth
│   │   └── router.go
│   └── worker/                     # Background jobs (hold release, booking expiry)
├── migrations/                     # SQL migrations (up/down)
├── test/
│   ├── concurrency/                # Race detector tests
│   └── integration/                # Full flow tests
├── docker-compose.yml
├── Dockerfile
├── .github/workflows/ci.yml
└── README.md
```

##  Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.22+ (for local dev without Docker)

### Run with Docker Compose (Recommended)

```bash
# Clone and start
git clone https://github.com/tuanwannafly/Go-booking.git
cd Go-booking
docker compose up --build

# API available at http://localhost:8080
# Swagger UI at http://localhost:8080/swagger/index.html
```

### Run Locally

```bash
# Start dependencies only
docker compose up -d postgres redis

# Run migrations
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path migrations -database "postgres://gobooking:gobooking@localhost:5432/gobooking?sslmode=disable" up

# Run server
go run ./cmd/api
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | development | Environment |
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | gobooking | DB user |
| `DB_PASSWORD` | gobooking | DB password |
| `DB_NAME` | gobooking | DB name |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6379 | Redis port |
| `JWT_SECRET` | dev-secret | JWT signing key |
| `LOG_LEVEL` | debug | Log level |
| `SERVER_PORT` | 8080 | HTTP port |

##  API Endpoints

| Method | Path | Auth | Idempotent | Description |
|--------|------|------|------------|-------------|
| GET | `/healthz` | - | - | Health check |
| GET | `/flights/search` | - | - | Search flights |
| GET | `/hotels/search` | - | - | Search hotels |
| POST | `/flights/{id}/seats/{seatId}/hold` | User | - | Hold seat (pessimistic) |
| POST | `/hotels/rooms/{roomId}/hold` | User | - | Hold room (optimistic) |
| POST | `/bookings` | User |  `Idempotency-Key` | Create booking |
| POST | `/bookings/{id}/confirm` | User | - | Confirm + mock payment |
| POST | `/bookings/{id}/cancel` | User | - | Cancel + refund calc |
| POST | `/bookings/{id}/schedule` | User | - | Reschedule booking (set `scheduled_at`) |
| GET | `/bookings/{id}` | User | - | Get booking details |
| GET | `/swagger/index.html` | - | - | API docs |

### Example Requests

**Search Flights**
```bash
curl "http://localhost:8080/flights/search?origin=SGN&destination=HAN&date=2026-07-15"
```

**Hold Seat (Pessimistic Lock)**
```bash
curl -X POST http://localhost:8080/flights/{flightId}/seats/{seatId}/hold \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"hold_duration": 10}'
```

**Create Booking (Idempotent)**
```bash
curl -X POST http://localhost:8080/bookings \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: unique-key-123" \
  -d '{"items":[{"inventory_unit_id":"<held-seat-id>"}]}'
```

Retrying the same request with the same `Idempotency-Key` returns the original
booking with HTTP `200` and `X-Idempotency-Replayed: true`. Reusing the key
with a different body returns HTTP `422`.

**Confirm Booking**
```bash
curl -X POST http://localhost:8080/bookings/{bookingId}/confirm \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"payment_method":"mock"}'
```

**Schedule a Booking** (set or update the actual travel date)
```bash
# Provide scheduled_at during booking creation
curl -X POST http://localhost:8080/bookings \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: unique-key-123" \
  -d '{
        "items":[{"inventory_unit_id":"<held-seat-id>"}],
        "scheduled_at":"2026-09-15T08:00:00Z"
      }'

# Or reschedule an existing booking
curl -X POST http://localhost:8080/bookings/{bookingId}/schedule \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"scheduled_at":"2026-09-16T08:00:00Z"}'
```

##  Testing

### Run All Tests (with Race Detector)
```bash
go test -race -cover ./...
```

### Concurrency Tests (The "Showcase" Tests)
```bash
# Run with Docker dependencies
docker compose up -d postgres redis
set GBOOKING_INTEGRATION=1
go test -race -v ./test/concurrency/...

# Expected output:
# === RUN   TestHoldSeat_ConcurrentRequests
# 100 goroutines, 1 succeeded, 99 rejected (ErrSeatUnavailable)
# --- PASS: TestHoldSeat_ConcurrentRequests (0.42s)
```

The concurrency suite creates isolated UUID-based fixtures, synchronizes all
goroutines before the hold attempt, and verifies both final inventory status
and the single version increment.

### Integration Tests
```bash
go test -v ./test/integration/...
```

##  Database Schema (Core Tables)

```sql
-- Core locking table
inventory_units (
  id UUID PK,
  resource_type ENUM('flight_seat','hotel_room'),
  resource_id UUID,      -- flight_id or room_type_id
  unit_code VARCHAR,     -- "A1", "101"
  status ENUM('available','held','booked','cancelled'),
  held_until TIMESTAMPTZ,
  version INT DEFAULT 1  -- for optimistic lock
)

-- Idempotency
idempotency_keys (
  key VARCHAR(255) PK,
  request_hash VARCHAR(64),
  response_body JSONB,
  status VARCHAR(20),
  expires_at TIMESTAMPTZ
)

-- Bookings (scheduled_at = real travel date, separate from expires_at = payment window)
bookings (
  id UUID PK,
  user_id UUID,
  status booking_status,  -- pending | confirmed | cancelled | expired
  idempotency_key VARCHAR(255) UNIQUE,
  total_amount DECIMAL(12,2),
  expires_at TIMESTAMPTZ,  -- when pending payment auto-expires
  scheduled_at TIMESTAMPTZ,-- real travel date (reschedulable)
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
)
```

##  Background Workers

| Worker | Interval | Action |
|--------|----------|--------|
| HoldReleaseWorker | 30s | `status=held AND held_until < NOW() → available` |
| BookingExpireWorker | 60s | `status=pending AND expires_at < NOW() → expired + release inventory` |

##  API Documentation

Swagger UI available at `http://localhost:8080/swagger/index.html` after running the server.

Generate docs:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs
```

##  CI/CD Pipeline

GitHub Actions (`.github/workflows/ci.yml`):
1. **Lint**: `golangci-lint`
2. **Test**: `go vet`, `go test -race -cover ./...`
3. **Build**: `docker build -t gobooking:ci .`

##  Interview Talking Points

1. **Pessimistic vs Optimistic Locking**
   - Flight seats: Redis distributed lock + `SELECT FOR UPDATE` (serializable isolation)
   - Hotel rooms: `version` column, `UPDATE ... WHERE version = $x`, retry 3x
   - Why different? Contention pattern: flight = burst (flash sale), hotel = spread over days

2. **Idempotency Implementation**
   - Middleware extracts `Idempotency-Key`
   - Hash request body → check `idempotency_keys` table
   - If exists + same hash → return cached response
   - If exists + different hash → 422 Conflict
   - Else process → store response → return

3. **Concurrency Proof**
   - `go test -race ./test/concurrency/...`
   - 100 goroutines × 1 seat → exactly 1 success
   - No data races detected by race detector

4. **Distributed Lock (Redsync)**
   - Redis-based mutex with auto-expiry
   - Acquired before DB transaction
   - Prevents thundering herd on DB

5. **Clean Architecture Benefits**
   - Handlers only handle HTTP
   - Services contain business rules (testable without HTTP)
   - Repositories abstract DB (swappable, mockable)
   - Domain entities have no external deps

##  Stretch Goals (Sprint 5)

- [ ] Prometheus metrics (`/metrics`) + Grafana dashboard
- [ ] JWT auth with roles (admin/user)
- [ ] k6 load test simulating flash sale
- [ ] Circuit breaker for payment provider

##  License

MIT — use freely for learning/portfolio.
