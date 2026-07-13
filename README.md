# GoBooking — Flight & Hotel Reservation Service

A production-ready Go backend service demonstrating **concurrency-safe booking**, **idempotency**, and **distributed locking** patterns.

## Key Features

- **Anti-Overbooking Protection**
  - Flight seats: Pessimistic locking via Redis distributed lock + `SELECT FOR UPDATE`
  - Hotel rooms: Optimistic locking via version column with automatic retry

- **Idempotency**: `Idempotency-Key` header prevents double-charging on client retry

- **Concurrency Verified**: 100-goroutine race-condition tests proving no overbooking

- **Clean Architecture**: Handler → Service → Repository → Database/Redis

- **Observable**: Structured logging, health checks, OpenAPI documentation

## Architecture

```
HTTP Handler → Service (Use Case) → Repository → PostgreSQL / Redis
```

### Locking Strategy

| Resource | Strategy | Rationale |
|----------|----------|-----------|
| Flight Seat | Pessimistic (Redis lock + `SELECT FOR UPDATE`) | Burst contention during flash sales |
| Hotel Room | Optimistic (`version` column + retry) | Contention spread over days, lower peak |

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.22+ |
| HTTP Framework | Gin |
| Database | PostgreSQL 16 (pgx/v5) |
| Cache/Lock | Redis 7 + Redsync |
| Migrations | golang-migrate |
| Testing | testing + testify + `-race` |
| Container | Docker Compose |
| CI/CD | GitHub Actions |
| API Docs | Swaggo (OpenAPI 3) |

## Project Structure

```
gobooking/
├── cmd/api/main.go                 # Application entry point
├── internal/
│   ├── config/                     # Viper configuration
│   ├── domain/                     # Domain entities + interfaces
│   ├── service/                    # Business logic
│   ├── repository/
│   │   ├── postgres/               # PostgreSQL implementations
│   │   └── redis/                  # Redis client + Redsync
│   ├── http/
│   │   ├── handler/                # Gin HTTP handlers
│   │   ├── middleware/            # Logging, rate-limiting, idempotency, auth
│   │   └── router.go
│   └── worker/                     # Background jobs (hold release, expiry)
├── migrations/                     # SQL migrations (up/down)
├── test/
│   ├── concurrency/                # Race detector tests
│   └── integration/                # Full flow tests
├── docker-compose.yml
├── Dockerfile
├── .github/workflows/ci.yml
└── README.md
```

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.22+ (for local development without Docker)

### Run with Docker Compose (Recommended)

```bash
git clone https://github.com/tuanwannafly/Go-booking.git
cd Go-booking
docker compose up --build
```

- API available at http://localhost:8080
- Swagger UI at http://localhost:8080/swagger/index.html

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
| `DB_USER` | gobooking | Database user |
| `DB_PASSWORD` | gobooking | Database password |
| `DB_NAME` | gobooking | Database name |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6379 | Redis port |
| `JWT_SECRET` | dev-secret | JWT signing key |
| `LOG_LEVEL` | debug | Log level |
| `SERVER_PORT` | 8080 | HTTP port |

## API Endpoints

| Method | Path | Auth | Idempotent | Description |
|--------|------|------|------------|-------------|
| GET | `/healthz` | - | - | Health check |
| GET | `/flights/search` | - | - | Search flights |
| GET | `/hotels/search` | - | - | Search hotels |
| POST | `/flights/{id}/seats/{seatId}/hold` | User | - | Hold seat (pessimistic lock) |
| POST | `/hotels/rooms/{roomId}/hold` | User | - | Hold room (optimistic lock) |
| POST | `/bookings` | User | `Idempotency-Key` | Create booking |
| POST | `/bookings/{id}/confirm` | User | - | Confirm booking + mock payment |
| POST | `/bookings/{id}/cancel` | User | - | Cancel booking + refund calculation |
| POST | `/bookings/{id}/schedule` | User | - | Reschedule booking |
| GET | `/bookings/{id}` | User | - | Get booking details |
| GET | `/swagger/index.html` | - | - | API documentation |

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

Retrying the same request with the same `Idempotency-Key` returns the original booking with HTTP `200` and `X-Idempotency-Replayed: true`. Reusing the key with a different body returns HTTP `422`.

**Confirm Booking**

```bash
curl -X POST http://localhost:8080/bookings/{bookingId}/confirm \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"payment_method":"mock"}'
```

**Schedule a Booking**

```bash
# Set scheduled_at during booking creation
curl -X POST http://localhost:8080/bookings \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: unique-key-123" \
  -d '{
        "items":[{"inventory_unit_id":"<held-seat-id>"}],
        "scheduled_at":"2026-09-15T08:00:00Z"
      }'

# Reschedule an existing booking
curl -X POST http://localhost:8080/bookings/{bookingId}/schedule \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"scheduled_at":"2026-09-16T08:00:00Z"}'
```

## Testing

### Run All Tests (with Race Detector)

```bash
go test -race -cover ./...
```

### Concurrency Tests

```bash
docker compose up -d postgres redis
set GBOOKING_INTEGRATION=1
go test -race -v ./test/concurrency/...
```

The concurrency suite creates isolated UUID-based fixtures, synchronizes all goroutines before the hold attempt, and verifies both final inventory status and the single version increment.

**Expected output:**
```
=== RUN   TestHoldSeat_ConcurrentRequests
100 goroutines, 1 succeeded, 99 rejected (ErrSeatUnavailable)
--- PASS: TestHoldSeat_ConcurrentRequests (0.42s)
```

### Integration Tests

```bash
go test -v ./test/integration/...
```

## Database Schema

### Core Tables

```sql
-- Inventory units (seats and rooms)
inventory_units (
  id UUID PK,
  resource_type ENUM('flight_seat','hotel_room'),
  resource_id UUID,      -- flight_id or room_type_id
  unit_code VARCHAR,     -- "A1", "101"
  status ENUM('available','held','booked','cancelled'),
  held_until TIMESTAMPTZ,
  version INT DEFAULT 1  -- for optimistic locking
)

-- Idempotency tracking
idempotency_keys (
  key VARCHAR(255) PK,
  request_hash VARCHAR(64),
  response_body JSONB,
  status VARCHAR(20),
  expires_at TIMESTAMPTZ
)

-- Bookings
bookings (
  id UUID PK,
  user_id UUID,
  status booking_status,     -- pending | confirmed | cancelled | expired
  idempotency_key VARCHAR(255) UNIQUE,
  total_amount DECIMAL(12,2),
  expires_at TIMESTAMPTZ,     -- payment window expiry
  scheduled_at TIMESTAMPTZ,   -- actual travel date
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
)
```

## Background Workers

| Worker | Interval | Action |
|--------|----------|--------|
| HoldReleaseWorker | 30s | `status=held AND held_until < NOW() → available` |
| BookingExpireWorker | 60s | `status=pending AND expires_at < NOW() → expired + release inventory` |

## API Documentation

Swagger UI is available at `http://localhost:8080/swagger/index.html` after running the server.

To generate documentation:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs
```

## CI/CD Pipeline

GitHub Actions (`.github/workflows/ci.yml`):

1. **Lint**: `golangci-lint`
2. **Test**: `go vet`, `go test -race -cover ./...`
3. **Build**: `docker build -t gobooking:ci .`

## Future Improvements

- [ ] Prometheus metrics (`/metrics`) + Grafana dashboard
- [ ] JWT auth with roles (admin/user)
- [ ] k6 load test simulating flash sale
- [ ] Circuit breaker for payment provider

## License

MIT
