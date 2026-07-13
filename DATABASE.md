# GoBooking — Database Schema & Seed Guide

Tài liệu này mô tả cấu trúc database, mối quan hệ giữa các bảng, cách seed dữ liệu fake phát triển/demo, và các tài khoản demo có sẵn sau khi chạy migration.

## 1. Tổng quan

Schema được chia thành 2 nhóm chính:

| Nhóm | Bảng | Mục đích |
|------|------|----------|
| **Domain** | `users`, `flights`, `hotels`, `room_types`, `inventory_units` | Danh mục cốt lõi: người dùng, chuyến bay, khách sạn, loại phòng, unit khả dụng |
| **Booking** | `bookings`, `booking_items`, `payments`, `idempotency_keys` | Luồng đặt chỗ: đặt chỗ, item, thanh toán, chống trùng |

Các migration hiện có:

| # | File | Nội dung |
|---|------|----------|
| 1 | `000001_create_users_table` | Bảng users + index email |
| 2 | `000002_create_flights_table` | Bảng flights + index origin/dest/departure |
| 3 | `000003_create_hotels_table` | Bảng hotels + index city |
| 4 | `000004_create_room_types_table` | Bảng room_types + FK -> hotels |
| 5 | `000005_create_inventory_units_table` | Bảng inventory_units + enum + index |
| 6 | `000006_create_bookings_table` | Bảng bookings + enum status |
| 7 | `000007_create_booking_items_table` | Bảng booking_items + FK bookings/inventory |
| 8 | `000008_create_idempotency_keys_table` | Bảng idempotency_keys |
| 9 | `000009_create_payments_table` | Bảng payments + enum status |
| 10 | `000010_add_booking_schedule` | Thêm `scheduled_at` cho bookings |
| 11 | `000011_seed_fake_data` | Seed dữ liệu fake Việt Nam |

## 2. Sơ đồ quan hệ

```text
users (1) ----< bookings (1) ----< booking_items (N) ----< inventory_units (N)
                                                               |
flights (1) ----< inventory_units                           |
                                                               |
hotels (1) ----< room_types (1) ----< inventory_units        |
                                                               |
bookings (1) ----< payments (1)                              |
                                                               |
bookings (1) ----< idempotency_keys (1) ---------------------+
```

## 3. Chi tiết từng bảng

### 3.1 `users`

| Cột | Kiểu | Mô tả |
|-----|------|-------|
| `id` | UUID PK | Khóa chính |
| `email` | VARCHAR(255) UNIQUE | Email đăng nhập |
| `password_hash` | VARCHAR(255) | Mật khẩu đã hash |
| `role` | VARCHAR(50) DEFAULT 'user' | Vai trò |
| `created_at` | TIMESTAMPTZ | Thời gian tạo |
| `updated_at` | TIMESTAMPTZ | Thời gian cập nhật |

### 3.2 `flights`

| Cột | Kiểu | Mô tả |
|-----|------|-------|
| `id` | UUID PK | Khóa chính |
| `code` | VARCHAR(50) UNIQUE | Mã chuyến bay, VD: `VN101` |
| `origin` | VARCHAR(100) | Mã sân bay đi |
| `destination` | VARCHAR(100) | Mã sân bay đến |
| `departure_time` | TIMESTAMPTZ | Giờ khởi hành |
| `arrival_time` | TIMESTAMPTZ | Giờ hạ cánh |
| `base_price` | DECIMAL(12,2) | Giá vé gốc (VND) |
| `created_at` | TIMESTAMPTZ | Thời gian tạo |
| `updated_at` | TIMESTAMPTZ | Thời gian cập nhật |

### 3.3 `hotels`

| Cột | Kiểu | Mô tả |
|-----|------|-------|
| `id` | UUID PK | Khóa chính |
| `name` | VARCHAR(255) | Tên khách sạn |
| `city` | VARCHAR(100) | Thành phố |
| `address` | VARCHAR(500) | Địa chỉ |
| `description` | TEXT | Mô tả |
| `created_at` | TIMESTAMPTZ | Thời gian tạo |
| `updated_at` | TIMESTAMPTZ | Thời gian cập nhật |

### 3.4 `room_types`

| Cột | Kiểu | Mô tả |
|-----|------|-------|
| `id` | UUID PK | Khóa chính |
| `hotel_id` | UUID FK -> hotels(id) | Khách sạn |
| `name` | VARCHAR(100) | Tên loại phòng |
| `base_price` | DECIMAL(12,2) | Giá/đêm (VND) |
| `capacity` | INT DEFAULT 2 | Số người tối đa |
| `created_at` | TIMESTAMPTZ | Thời gian tạo |
| `updated_at` | TIMESTAMPTZ | Thời gian cập nhật |

### 3.5 `inventory_units`

Bảng trung tâm quản lý khả dụng của từng seat/room.

| Cột | Kiểu | Mô tả |
|-----|------|-------|
| `id` | UUID PK | Khóa chính |
| `resource_type` | ENUM('flight_seat','hotel_room') | Loại resource |
| `resource_id` | UUID | FK đến `flights.id` hoặc `room_types.id` |
| `unit_code` | VARCHAR(50) | Mã unit, VD: `J1A`, `RM-1` |
| `status` | ENUM('available','held','booked','cancelled') | Trạng thái |
| `held_until` | TIMESTAMPTZ | Thời gian hết hold |
| `version` | INT DEFAULT 1 | Dùng cho optimistic lock |
| `created_at` | TIMESTAMPTZ | Thời gian tạo |
| `updated_at` | TIMESTAMPTZ | Thời gian cập nhật |

### 3.6 `bookings`

| Cột | Kiểu | Mô tả |
|-----|------|-------|
| `id` | UUID PK | Khóa chính |
| `user_id` | UUID FK -> users(id) | Người đặt |
| `status` | ENUM('pending','confirmed','cancelled','expired') | Trạng thái |
| `idempotency_key` | VARCHAR(255) UNIQUE | Key chống trùng |
| `total_amount` | DECIMAL(12,2) | Tổng tiền (VND) |
| `expires_at` | TIMESTAMPTZ | Hạn thanh toán |
| `scheduled_at` | TIMESTAMPTZ | Ngày đi/check-in thực tế |
| `created_at` | TIMESTAMPTZ | Thời gian tạo |
| `updated_at` | TIMESTAMPTZ | Thời gian cập nhật |

### 3.7 `booking_items`

| Cột | Kiểu | Mô tả |
|-----|------|-------|
| `id` | UUID PK | Khóa chính |
| `booking_id` | UUID FK -> bookings(id) | Đặt chỗ |
| `inventory_unit_id` | UUID FK -> inventory_units(id) | Unit được đặt |
| `price` | DECIMAL(12,2) | Giá tại thời điểm đặt |
| `created_at` | TIMESTAMPTZ | Thời gian tạo |

### 3.8 `payments`

| Cột | Kiểu | Mô tả |
|-----|------|-------|
| `id` | UUID PK | Khóa chính |
| `booking_id` | UUID FK -> bookings(id) | Đặt chỗ |
| `status` | ENUM('pending','completed','failed','refunded') | Trạng thái |
| `amount` | DECIMAL(12,2) | Số tiền |
| `provider_ref` | VARCHAR(255) | Mã tham chiếu provider |
| `created_at` | TIMESTAMPTZ | Thời gian tạo |
| `updated_at` | TIMESTAMPTZ | Thời gian cập nhật |

### 3.9 `idempotency_keys`

| Cột | Kiểu | Mô tả |
|-----|------|-------|
| `key` | VARCHAR(255) PK | Key idempotency |
| `request_hash` | VARCHAR(64) | Hash request body |
| `response_body` | JSONB | Response đã cache |
| `status` | VARCHAR(20) | processing/completed/failed |
| `created_at` | TIMESTAMPTZ | Thời gian tạo |
| `expires_at` | TIMESTAMPTZ | Hạn cache |

## 4. Index & Performance

| Index | Bảng | Mục đích |
|-------|------|----------|
| `idx_users_email` | users | Tìm user theo email |
| `idx_flights_origin_dest_departure` | flights | Tìm chuyến bay theo route + ngày |
| `idx_hotels_city` | hotels | Tìm khách sạn theo thành phố |
| `idx_room_types_hotel_id` | room_types | Lấy room types của 1 khách sạn |
| `idx_inventory_units_resource` | inventory_units | Lấy units theo flight/room |
| `idx_inventory_units_status` | inventory_units | Lọc theo trạng thái |
| `idx_inventory_units_held_until` | inventory_units | Tìm hold hết hạn |
| `idx_bookings_user_id` | bookings | Lịch sử đặt chỗ của user |
| `idx_bookings_status` | bookings | Lọc đặt chỗ theo trạng thái |
| `idx_bookings_idempotency_key` | bookings | Kiểm tra key trùng |
| `idx_bookings_expires_at` | bookings | Tìm booking hết hạn |
| `idx_booking_items_booking_id` | booking_items | Lấy items của 1 booking |
| `idx_booking_items_inventory_unit_id` | booking_items | Tìm booking theo unit |
| `idx_payments_booking_id` | payments | Lấy payment của 1 booking |
| `idx_payments_status` | payments | Lọc payment theo trạng thái |
| `idx_idempotency_keys_expires_at` | idempotency_keys | Dọn key hết hạn |

## 5. Seed Data (Migration 000011)

Migration `000011_seed_fake_data` chứa dữ liệu fake tiếng Việt để phát triển và demo.

### 5.1 Users (3 tài khoản)

| Email | Password (demo) | Role |
|-------|----------------|------|
| `nguyen.van.a@example.com` | `hashed_demo_1` | user |
| `tran.thi.b@example.com` | `hashed_demo_2` | user |
| `le.van.c@example.com` | `hashed_demo_3` | user |

### 5.2 Flights (6 chuyến)

| Mã | Route | Ngày | Giờ | Giá (VND) |
|----|-------|------|-----|-----------|
| VN101 | SGN → HAN | 2026-07-20 | 08:00 - 10:30 | 1,200,000 |
| VN102 | HAN → SGN | 2026-07-20 | 14:00 - 16:30 | 1,200,000 |
| VN103 | SGN → DAD | 2026-07-21 | 09:00 - 10:30 | 800,000 |
| VN104 | DAD → SGN | 2026-07-21 | 17:00 - 18:30 | 850,000 |
| VN105 | HAN → DAD | 2026-07-22 | 11:00 - 12:30 | 600,000 |
| VN201 | SGN → PQC | 2026-07-23 | 07:00 - 09:00 | 900,000 |

Mỗi chuyến có **120 ghế** (12 hàng Business J1-J12 + 8 hàng Economy 13-20, mỗi hàng 6 ghế A-F).

### 5.3 Hotels (3 khách sạn)

| Tên | Thành phố | Địa chỉ |
|-----|-----------|---------|
| Khách sạn Sài Gòn | Ho Chi Minh | 123 Đồng Khởi, Quận 1 |
| Khách sạn Hà Nội | Hanoi | 456 Phố Cổ, Hoàn Kiếm |
| Khách sạn Đà Nẵng | Da Nang | 789 Bãi biển Mỹ Khê |

Mỗi khách sạn có **2 loại phòng**: Standard + Deluxe.

### 5.4 Room Types (6 loại)

| Khách sạn | Loại | Giá/đêm (VND) | Sức chứa |
|-----------|------|---------------|----------|
| Sài Gòn | Standard | 800,000 | 2 |
| Sài Gòn | Deluxe | 1,800,000 | 3 |
| Hà Nội | Standard | 700,000 | 2 |
| Hà Nội | Deluxe | 1,600,000 | 3 |
| Đà Nẵng | Standard | 600,000 | 2 |
| Đà Nẵng | Deluxe | 1,400,000 | 4 |

Mỗi loại có **5 phòng** (RM-1 đến RM-5).

### 5.5 Bookings (3 đặt chỗ demo)

| Booking | User | Status | Tổng (VND) | Ghi chú |
|---------|------|--------|-----------|---------|
| BOOK-1 | nguyen.van.a | confirmed | 1,200,000 | Đã thanh toán |
| BOOK-2 | tran.thi.b | pending | 800,000 | Đang giữ phòng |
| BOOK-3 | le.van.c | cancelled | 800,000 | Đã hủy |

### 5.6 Tổng số record seed

| Bảng | Số lượng |
|------|----------|
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

## 6. Cách chạy migration

### 6.1 Cài đặt tool

```bash
# macOS/Linux
brew install golang-migrate

# Windows (PowerShell)
# Tải binary từ https://github.com/golang-migrate/migrate/releases
# Hoặc dùng Chocolatey: choco install golang-migrate
```

### 6.2 Chạy tất cả migrations

```bash
migrate -path migrations -database "postgres://gobooking:gobooking@localhost:5432/gobooking?sslmode=disable" up
```

### 6.3 Rollback migration cuối

```bash
migrate -path migrations -database "postgres://gobooking:gobooking@localhost:5432/gobooking?sslmode=disable" down 1
```

### 6.4 Kiểm tra version hiện tại

```bash
migrate -path migrations -database "postgres://gobooking:gobooking@localhost:5432/gobooking?sslmode=disable" version
```

## 7. Demo Accounts

Sau khi chạy `000011_seed_fake_data`, bạn có thể dùng các tài khoản sau để test frontend:

```
Email: nguyen.van.a@example.com
→ Đã có 1 booking confirmed (VN101)

Email: tran.thi.b@example.com
→ Đã có 1 booking pending (khách sạn Sài Gòn)

Email: le.van.c@example.com
→ Đã có 1 booking cancelled (VN103)
```

> **Lưu ý:** Đây là demo data. Mật khẩu là hash giả (`hashed_demo_1/2/3`). Backend hiện tại chưa có endpoint đăng nhập thực, frontend dùng mock auth. Dữ liệu này chủ yếu để search/hold/booking API có data thật để test.

## 8. Lưu ý kỹ thuật

1. **UUID cố định**: Seed data dùng UUID cố định (ví dụ: `aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa` cho VN101) để dễ debug và tham chiếu.
2. **Timestamps**: Ngày giờ seed dùng UTC (`+00`). Các booking có `expires_at` dùng `NOW() + INTERVAL '10 minutes'` để luôn trong tương lai khi chạy migration.
3. **Inventory status**:
   - `available`: unit có thể đặt
   - `held`: đang được giữ 10 phút
   - `booked`: đã xác nhận thanh toán
   - `cancelled`: đã hủy
4. **Down migration**: Xóa theo UUID cụ thể, không dùng `TRUNCATE` để tránh mất data của user khác.
5. **Idempotency**: Seed 2 idempotency keys để test retry không trùng.

## 9. Query hữu ích

```sql
-- Xem tất cả chuyến bay còn ghế
SELECT f.code, f.origin, f.destination, f.departure_time, f.base_price,
       COUNT(iu.id) FILTER (WHERE iu.status = 'available') AS seats_left
FROM flights f
LEFT JOIN inventory_units iu ON iu.resource_id = f.id AND iu.resource_type = 'flight_seat'
GROUP BY f.id
ORDER BY f.departure_time;

-- Xem phòng trống của 1 khách sạn
SELECT h.name, rt.name AS room_type, rt.base_price,
       COUNT(iu.id) FILTER (WHERE iu.status = 'available') AS rooms_left
FROM hotels h
JOIN room_types rt ON rt.hotel_id = h.id
LEFT JOIN inventory_units iu ON iu.resource_id = rt.id AND iu.resource_type = 'hotel_room'
WHERE h.city = 'Ho Chi Minh'
GROUP BY h.id, rt.id;

-- Xem booking của user
SELECT b.id, b.status, b.total_amount, b.created_at,
       bi.price, iu.unit_code, f.code
FROM bookings b
JOIN booking_items bi ON bi.booking_id = b.id
JOIN inventory_units iu ON iu.id = bi.inventory_unit_id
LEFT JOIN flights f ON f.id = iu.resource_id AND iu.resource_type = 'flight_seat'
WHERE b.user_id = '11111111-1111-1111-1111-111111111111'
ORDER BY b.created_at DESC;

-- Dọn hold hết hạn (giả lập worker)
UPDATE inventory_units
SET status = 'available', held_until = NULL, version = version + 1, updated_at = NOW()
WHERE status = 'held' AND held_until < NOW();

-- Đếm tổng inventory đang hold
SELECT resource_type, status, COUNT(*) 
FROM inventory_units 
GROUP BY resource_type, status;
```
