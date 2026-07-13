# GoBooking — Database Schema & Seed Guide

This document describes the database structure, table relationships, seed data setup, and demo accounts.

## 1. Overview

The schema is organized into 2 main groups:

| Group | Tables | Purpose |
|-------|--------|---------|
| **Domain** | `users`, `flights`, `hotels`, `room_types`, `inventory_units` | Core catalog: users, flights, hotels, room types, available units |
| **Booking** | `bookings`, `booking_items`, `payments`, `idempotency_keys` | Booking flow: reservations, items, payments, idempotency |

### Available Migrations

| # | File | Description |
|---|------|-------------|
| 1 | `000001_create_users_table` | Users table + email index |
| 2 | `000002_create_flights_table` | Flights table + origin/dest/departure index |
| 3 | `000003_create_hotels_table` | Hotels table + city index |
| 4 | `000004_create_room_types_table` | Room types table + FK to hotels |
| 5 | `000005_create_inventory_units_table` | Inventory units table + enum + index |
| 6 | `000006_create_bookings_table` | Bookings table + status enum |
| 7 | `000007_create_booking_items_table` | Booking items + FK to bookings/inventory |
| 8 | `000008_create_idempotency_keys_table` | Idempotency keys table |
| 9 | `000009_create_payments_table` | Payments table + status enum |
| 10 | `000010_add_booking_schedule` | Add `scheduled_at` to bookings |
| 11 | `000011_seed_fake_data` | Seed fake Vietnam data |

## 2. Entity Relationship Diagram

```
users (1) ----< bookings (1) ----< booking_items (N) ----< inventory_units (N)
                                                               |
flights (1) ----< inventory_units                           |
                                                               |
hotels (1) ----< room_types (1) ----< inventory_units      |
                                                               |
bookings (1) ----< payments (1)                           |
                                                               |
bookings (1) ----< idempotency_keys (1) ---------------------+
```

## 3. Table Details

### 3.1 `users`

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID PK | Primary key |
| `email` | VARCHAR(255) UNIQUE | Login email |
| `password_hash` | VARCHAR(255) | Hashed password |
| `role` | VARCHAR(50) DEFAULT 'user' | User role |
| `created_at` | TIMESTAMPTZ | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | Update timestamp |

### 3.2 `flights`

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID PK | Primary key |
| `code` | VARCHAR(50) UNIQUE | Flight code, e.g., `VN101` |
| `origin` | VARCHAR(100) | Departure airport code |
| `destination` | VARCHAR(100) | Arrival airport code |
| `departure_time` | TIMESTAMPTZ | Departure time |
| `arrival_time` | TIMESTAMPTZ | Arrival time |
| `base_price` | DECIMAL(12,2) | Base ticket price (VND) |
| `created_at` | TIMESTAMPTZ | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | Update timestamp |

### 3.3 `hotels`

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID PK | Primary key |
| `name` | VARCHAR(255) | Hotel name |
| `city` | VARCHAR(100) | City |
| `address` | VARCHAR(500) | Address |
| `description` | TEXT | Description |
| `created_at` | TIMESTAMPTZ | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | Update timestamp |

### 3.4 `room_types`

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID PK | Primary key |
| `hotel_id` | UUID FK -> hotels(id) | Hotel reference |
| `name` | VARCHAR(100) | Room type name |
| `base_price` | DECIMAL(12,2) | Price per night (VND) |
| `capacity` | INT DEFAULT 2 | Maximum occupancy |
| `created_at` | TIMESTAMPTZ | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | Update timestamp |

### 3.5 `inventory_units`

Central table managing availability for each seat/room.

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID PK | Primary key |
| `resource_type` | ENUM('flight_seat','hotel_room') | Resource type |
| `resource_id` | UUID | FK to `flights.id` or `room_types.id` |
| `unit_code` | VARCHAR(50) | Unit code, e.g., `J1A`, `RM-1` |
| `status` | ENUM('available','held','booked','cancelled') | Availability status |
| `held_until` | TIMESTAMPTZ | Hold expiry time |
| `version` | INT DEFAULT 1 | Used for optimistic locking |
| `created_at` | TIMESTAMPTZ | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | Update timestamp |

### 3.6 `bookings`

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID PK | Primary key |
| `user_id` | UUID FK -> users(id) | Booking user |
| `status` | ENUM('pending','confirmed','cancelled','expired') | Booking status |
| `idempotency_key` | VARCHAR(255) UNIQUE | Idempotency key |
| `total_amount` | DECIMAL(12,2) | Total amount (VND) |
| `expires_at` | TIMESTAMPTZ | Payment deadline |
| `scheduled_at` | TIMESTAMPTZ | Actual travel/check-in date |
| `created_at` | TIMESTAMPTZ | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | Update timestamp |

### 3.7 `booking_items`

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID PK | Primary key |
| `booking_id` | UUID FK -> bookings(id) | Booking reference |
| `inventory_unit_id` | UUID FK -> inventory_units(id) | Reserved unit |
| `price` | DECIMAL(12,2) | Price at booking time |
| `created_at` | TIMESTAMPTZ | Creation timestamp |

### 3.8 `payments`

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID PK | Primary key |
| `booking_id` | UUID FK -> bookings(id) | Booking reference |
| `status` | ENUM('pending','completed','failed','refunded') | Payment status |
| `amount` | DECIMAL(12,2) | Payment amount |
| `provider_ref` | VARCHAR(255) | Provider reference |
| `created_at` | TIMESTAMPTZ | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | Update timestamp |

### 3.9 `idempotency_keys`

| Column | Type | Description |
|--------|------|-------------|
| `key` | VARCHAR(255) PK | Idempotency key |
| `request_hash` | VARCHAR(64) | Request body hash |
| `response_body` | JSONB | Cached response |
| `status` | VARCHAR(20) | processing/completed/failed |
| `created_at` | TIMESTAMPTZ | Creation timestamp |
| `expires_at` | TIMESTAMPTZ | Cache expiry |

## 4. Indexes & Performance

| Index | Table | Purpose |
|-------|-------|---------|
| `idx_users_email` | users | Find user by email |
| `idx_flights_origin_dest_departure` | flights | Search by route + date |
| `idx_hotels_city` | hotels | Search by city |
| `idx_room_types_hotel_id` | room_types | Get room types for a hotel |
| `idx_inventory_units_resource` | inventory_units | Get units by flight/room |
| `idx_inventory_units_status` | inventory_units | Filter by status |
| `idx_inventory_units_held_until` | inventory_units | Find expired holds |
| `idx_bookings_user_id` | bookings | User booking history |
| `idx_bookings_status` | bookings | Filter by status |
| `idx_bookings_idempotency_key` | bookings | Check for duplicate key |
| `idx_bookings_expires_at` | bookings | Find expired bookings |
| `idx_booking_items_booking_id` | booking_items | Get items for a booking |
| `idx_booking_items_inventory_unit_id` | booking_items | Find booking by unit |
| `idx_payments_booking_id` | payments | Get payment for a booking |
| `idx_payments_status` | payments | Filter by payment status |
| `idx_idempotency_keys_expires_at` | idempotency_keys | Clean up expired keys |

## 5. Seed Data (Migration 000011)

Migration `000011_seed_fake_data` contains fake Vietnamese data for development and demo.

### 5.1 Users (3 accounts)

| Email | Password (demo) | Role |
|-------|----------------|------|
| `nguyen.van.a@example.com` | `hashed_demo_1` | user |
| `tran.thi.b@example.com` | `hashed_demo_2` | user |
| `le.van.c@example.com` | `hashed_demo_3` | user |

### 5.2 Flights (6 flights)

| Code | Route | Date | Time | Price (VND) |
|------|-------|------|------|-------------|
| VN101 | SGN → HAN | 2026-07-20 | 08:00 - 10:30 | 1,200,000 |
| VN102 | HAN → SGN | 2026-07-20 | 14:00 - 16:30 | 1,200,000 |
| VN103 | SGN → DAD | 2026-07-21 | 09:00 - 10:30 | 800,000 |
| VN104 | DAD → SGN | 2026-07-21 | 17:00 - 18:30 | 850,000 |
| VN105 | HAN → DAD | 2026-07-22 | 11:00 - 12:30 | 600,000 |
| VN201 | SGN → PQC | 2026-07-23 | 07:00 - 09:00 | 900,000 |

Each flight has **120 seats** (12 Business rows J1-J12 + 8 Economy rows 13-20, each row 6 seats A-F).

### 5.3 Hotels (3 hotels)

| Name | City | Address |
|------|------|---------|
| Saigon Hotel | Ho Chi Minh | 123 Dong Khoi, District 1 |
| Hanoi Hotel | Hanoi | 456 Old Quarter, Hoan Kiem |
| Da Nang Hotel | Da Nang | 789 My Khe Beach |

Each hotel has **2 room types**: Standard + Deluxe.

### 5.4 Room Types (6 types)

| Hotel | Type | Price/night (VND) | Capacity |
|-------|------|-------------------|----------|
| Saigon | Standard | 800,000 | 2 |
| Saigon | Deluxe | 1,800,000 | 3 |
| Hanoi | Standard | 700,000 | 2 |
| Hanoi | Deluxe | 1,600,000 | 3 |
| Da Nang | Standard | 600,000 | 2 |
| Da Nang | Deluxe | 1,400,000 | 4 |

Each type has **5 rooms** (RM-1 through RM-5).

### 5.5 Bookings (3 demo bookings)

| Booking | User | Status | Total (VND) | Notes |
|---------|------|--------|-------------|-------|
| BOOK-1 | nguyen.van.a | confirmed | 1,200,000 | Paid |
| BOOK-2 | tran.thi.b | pending | 800,000 | Room on hold |
| BOOK-3 | le.van.c | cancelled | 800,000 | Cancelled |

### 5.6 Total Seed Records

| Table | Count |
|-------|-------|
| users | 3 |
| flights | 6 |
| hotels | 3 |
| room_types | 6 |
| inventory_units (flight seats) | 720 |
| inventory_units (hotel rooms) | 30 |
| bookings | 3 |
| booking_items | 3 |
| payments | 1 |
| idempotency_keys | 2 |

## 6. Running Migrations

### 6.1 Install the Tool

```bash
# macOS/Linux
brew install golang-migrate

# Windows (PowerShell)
# Download binary from https://github.com/golang-migrate/migrate/releases
# Or use Chocolatey: choco install golang-migrate
```

### 6.2 Run All Migrations

```bash
migrate -path migrations -database "postgres://gobooking:gobooking@localhost:5432/gobooking?sslmode=disable" up
```

### 6.3 Rollback Last Migration

```bash
migrate -path migrations -database "postgres://gobooking:gobooking@localhost:5432/gobooking?sslmode=disable" down 1
```

### 6.4 Check Current Version

```bash
migrate -path migrations -database "postgres://gobooking:gobooking@localhost:5432/gobooking?sslmode=disable" version
```

## 7. Demo Accounts

After running `000011_seed_fake_data`, use these accounts to test the frontend:

```
Email: nguyen.van.a@example.com
→ Has 1 confirmed booking (VN101)

Email: tran.thi.b@example.com
→ Has 1 pending booking (Saigon Hotel)

Email: le.van.c@example.com
→ Has 1 cancelled booking (VN103)
```

> **Note:** This is demo data. Passwords are mock hashes (`hashed_demo_1/2/3`). The backend currently does not have a real login endpoint; the frontend uses mock auth. This data primarily provides real data for search/hold/booking API testing.

## 8. Technical Notes

1. **Fixed UUIDs**: Seed data uses fixed UUIDs (e.g., `aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa` for VN101) for easier debugging and reference.

2. **Timestamps**: Seed dates use UTC (`+00`). Bookings have `expires_at` set using `NOW() + INTERVAL '10 minutes'` to always be in the future.

3. **Inventory Status**:
   - `available`: unit can be booked
   - `held`: temporarily reserved for 10 minutes
   - `booked`: confirmed and paid
   - `cancelled`: cancelled

4. **Down Migration**: Deletes by specific UUID, not using `TRUNCATE` to avoid losing other users' data.

5. **Idempotency**: Seeds 2 idempotency keys for retry testing.

## 9. Useful Queries

```sql
-- View all flights with available seats
SELECT f.code, f.origin, f.destination, f.departure_time, f.base_price,
       COUNT(iu.id) FILTER (WHERE iu.status = 'available') AS seats_left
FROM flights f
LEFT JOIN inventory_units iu ON iu.resource_id = f.id AND iu.resource_type = 'flight_seat'
GROUP BY f.id
ORDER BY f.departure_time;

-- View available rooms for a hotel
SELECT h.name, rt.name AS room_type, rt.base_price,
       COUNT(iu.id) FILTER (WHERE iu.status = 'available') AS rooms_left
FROM hotels h
JOIN room_types rt ON rt.hotel_id = h.id
LEFT JOIN inventory_units iu ON iu.resource_id = rt.id AND iu.resource_type = 'hotel_room'
WHERE h.city = 'Ho Chi Minh'
GROUP BY h.id, rt.id;

-- View user's bookings
SELECT b.id, b.status, b.total_amount, b.created_at,
       bi.price, iu.unit_code, f.code
FROM bookings b
JOIN booking_items bi ON bi.booking_id = b.id
JOIN inventory_units iu ON iu.id = bi.inventory_unit_id
LEFT JOIN flights f ON f.id = iu.resource_id AND iu.resource_type = 'flight_seat'
WHERE b.user_id = '11111111-1111-1111-1111-111111111111'
ORDER BY b.created_at DESC;

-- Release expired holds (simulate worker)
UPDATE inventory_units
SET status = 'available', held_until = NULL, version = version + 1, updated_at = NOW()
WHERE status = 'held' AND held_until < NOW();

-- Count inventory by status
SELECT resource_type, status, COUNT(*) 
FROM inventory_units 
GROUP BY resource_type, status;
```
