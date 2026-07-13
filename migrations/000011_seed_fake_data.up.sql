-- Migration: 000011_seed_fake_data.up.sql
-- Seed dữ liệu fake tiếng Việt cho development/demo.
-- Dùng UUID cố định để dễ debug và rollback an toàn.

-- ============================================================
-- 1. USERS
-- ============================================================
INSERT INTO users (id, email, password_hash, role, created_at, updated_at) VALUES
  ('11111111-1111-1111-1111-111111111111', 'nguyen.van.a@example.com', '$2a$10$Ngq1xHxJQ8qHxJQ8qHxJQOw2X9Y7Z6a5b4c3d2e1f0a9b', 'user', '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('22222222-2222-2222-2222-222222222222', 'tran.thi.b@example.com', '$2a$10$Pmr2yIyKR9rIyKR9rIyKRP3Y8Z0A7B6C5D4E3F2G1H0I', 'user', '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('33333333-3333-3333-3333-333333333333', 'le.van.c@example.com', '$2a$10$Qns3zJzLS0sJzLS0sJzLSQ4Z9A1B8C7D6E5F4G3H2I1J0K', 'user', '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00');

-- ============================================================
-- 2. FLIGHTS
-- ============================================================
INSERT INTO flights (id, code, origin, destination, departure_time, arrival_time, base_price, created_at, updated_at) VALUES
  ('10000000-1000-1000-1000-100000000001', 'VN101', 'SGN', 'HAN', '2026-07-20 08:00:00+00', '2026-07-20 10:30:00+00', 1200000, '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('10000000-1000-1000-1000-100000000002', 'VN102', 'HAN', 'SGN', '2026-07-20 14:00:00+00', '2026-07-20 16:30:00+00', 1200000, '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('10000000-1000-1000-1000-100000000003', 'VN103', 'SGN', 'DAD', '2026-07-21 09:00:00+00', '2026-07-21 10:30:00+00', 800000, '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('10000000-1000-1000-1000-100000000004', 'VN104', 'DAD', 'SGN', '2026-07-21 17:00:00+00', '2026-07-21 18:30:00+00', 850000, '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('10000000-1000-1000-1000-100000000005', 'VN105', 'HAN', 'DAD', '2026-07-22 11:00:00+00', '2026-07-22 12:30:00+00', 600000, '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('10000000-1000-1000-1000-100000000006', 'VN201', 'SGN', 'PQC', '2026-07-23 07:00:00+00', '2026-07-23 09:00:00+00', 900000, '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00');

-- ============================================================
-- 3. HOTELS
-- ============================================================
INSERT INTO hotels (id, name, city, address, description, created_at, updated_at) VALUES
  ('20000000-2000-2000-2000-200000000001', 'Khách sạn Sài Gòn', 'Ho Chi Minh', '123 Đồng Khởi, Quận 1', 'Khách sạn trung tâm thành phố với view sông Sài Gòn', '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('20000000-2000-2000-2000-200000000002', 'Khách sạn Hà Nội', 'Hanoi', '456 Phố Cổ, Hoàn Kiếm', 'Khách sạn kiểu cổ điển gần Hồ Gươm', '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('20000000-2000-2000-2000-200000000003', 'Khách sạn Đà Nẵng', 'Da Nang', '789 Bãi biển Mỹ Khê', 'Resort biển với hồ bơi vô cực', '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00');

-- ============================================================
-- 4. ROOM TYPES
-- ============================================================
INSERT INTO room_types (id, hotel_id, name, base_price, capacity, created_at, updated_at) VALUES
  ('30000000-3000-3000-3000-300000000001', '20000000-2000-2000-2000-200000000001', 'Standard', 800000, 2, '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('30000000-3000-3000-3000-300000000002', '20000000-2000-2000-2000-200000000001', 'Deluxe', 1800000, 3, '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('30000000-3000-3000-3000-300000000003', '20000000-2000-2000-2000-200000000002', 'Standard', 700000, 2, '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('30000000-3000-3000-3000-300000000004', '20000000-2000-2000-2000-200000000002', 'Deluxe', 1600000, 3, '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('30000000-3000-3000-3000-300000000005', '20000000-2000-2000-2000-200000000003', 'Standard', 600000, 2, '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00'),
  ('30000000-3000-3000-3000-300000000006', '20000000-2000-2000-2000-200000000003', 'Deluxe', 1400000, 4, '2026-06-15 00:00:00+00', '2026-06-15 00:00:00+00');

-- ============================================================
-- 5. INVENTORY UNITS
-- ============================================================

-- 5.1 Flight seats: 20 rows x 6 seats = 120 seats/flight
-- Rows 1-12: Business (J1A..J12F), Rows 13-20: Economy (13A..20F)
INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
SELECT
    gen_random_uuid(),
    'flight_seat',
    '10000000-1000-1000-1000-100000000001'::uuid,
    CASE WHEN s.row_num <= 12 THEN 'J' || s.row_num || l.letter ELSE s.row_num || l.letter END,
    'available',
    1
FROM generate_series(1, 20) AS s(row_num)
CROSS JOIN (VALUES ('A'),('B'),('C'),('D'),('E'),('F')) AS l(letter);

INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
SELECT
    gen_random_uuid(),
    'flight_seat',
    '10000000-1000-1000-1000-100000000002'::uuid,
    CASE WHEN s.row_num <= 12 THEN 'J' || s.row_num || l.letter ELSE s.row_num || l.letter END,
    'available',
    1
FROM generate_series(1, 20) AS s(row_num)
CROSS JOIN (VALUES ('A'),('B'),('C'),('D'),('E'),('F')) AS l(letter);

INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
SELECT
    gen_random_uuid(),
    'flight_seat',
    '10000000-1000-1000-1000-100000000003'::uuid,
    CASE WHEN s.row_num <= 12 THEN 'J' || s.row_num || l.letter ELSE s.row_num || l.letter END,
    'available',
    1
FROM generate_series(1, 20) AS s(row_num)
CROSS JOIN (VALUES ('A'),('B'),('C'),('D'),('E'),('F')) AS l(letter);

INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
SELECT
    gen_random_uuid(),
    'flight_seat',
    '10000000-1000-1000-1000-100000000004'::uuid,
    CASE WHEN s.row_num <= 12 THEN 'J' || s.row_num || l.letter ELSE s.row_num || l.letter END,
    'available',
    1
FROM generate_series(1, 20) AS s(row_num)
CROSS JOIN (VALUES ('A'),('B'),('C'),('D'),('E'),('F')) AS l(letter);

INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
SELECT
    gen_random_uuid(),
    'flight_seat',
    '10000000-1000-1000-1000-100000000005'::uuid,
    CASE WHEN s.row_num <= 12 THEN 'J' || s.row_num || l.letter ELSE s.row_num || l.letter END,
    'available',
    1
FROM generate_series(1, 20) AS s(row_num)
CROSS JOIN (VALUES ('A'),('B'),('C'),('D'),('E'),('F')) AS l(letter);

INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
SELECT
    gen_random_uuid(),
    'flight_seat',
    '10000000-1000-1000-1000-100000000006'::uuid,
    CASE WHEN s.row_num <= 12 THEN 'J' || s.row_num || l.letter ELSE s.row_num || l.letter END,
    'available',
    1
FROM generate_series(1, 20) AS s(row_num)
CROSS JOIN (VALUES ('A'),('B'),('C'),('D'),('E'),('F')) AS l(letter);

-- 5.2 Hotel rooms: 5 rooms/room_type
INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
SELECT gen_random_uuid(), 'hotel_room', '30000000-3000-3000-3000-300000000001'::uuid, 'RM-' || g.n, 'available', 1
FROM generate_series(1, 5) AS g(n);

INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
SELECT gen_random_uuid(), 'hotel_room', '30000000-3000-3000-3000-300000000002'::uuid, 'RM-' || g.n, 'available', 1
FROM generate_series(1, 5) AS g(n);

INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
SELECT gen_random_uuid(), 'hotel_room', '30000000-3000-3000-3000-300000000003'::uuid, 'RM-' || g.n, 'available', 1
FROM generate_series(1, 5) AS g(n);

INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
SELECT gen_random_uuid(), 'hotel_room', '30000000-3000-3000-3000-300000000004'::uuid, 'RM-' || g.n, 'available', 1
FROM generate_series(1, 5) AS g(n);

INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
SELECT gen_random_uuid(), 'hotel_room', '30000000-3000-3000-3000-300000000005'::uuid, 'RM-' || g.n, 'available', 1
FROM generate_series(1, 5) AS g(n);

INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
SELECT gen_random_uuid(), 'hotel_room', '30000000-3000-3000-3000-300000000006'::uuid, 'RM-' || g.n, 'available', 1
FROM generate_series(1, 5) AS g(n);

-- ============================================================
-- 6. BOOKINGS
-- ============================================================
INSERT INTO bookings (id, user_id, status, idempotency_key, total_amount, expires_at, scheduled_at, created_at, updated_at) VALUES
  ('40000000-4000-4000-4000-400000000001', '11111111-1111-1111-1111-111111111111', 'confirmed', NULL, 1200000, NOW() + INTERVAL '10 minutes', '2026-07-20 08:00:00+00', '2026-07-10 08:00:00+00', '2026-07-10 08:05:00+00'),
  ('40000000-4000-4000-4000-400000000002', '22222222-2222-2222-2222-222222222222', 'pending', 'demo-idempotency-key-2', 800000, NOW() + INTERVAL '10 minutes', '2026-07-25 00:00:00+00', '2026-07-10 09:00:00+00', '2026-07-10 09:00:00+00'),
  ('40000000-4000-4000-4000-400000000003', '33333333-3333-3333-3333-333333333333', 'cancelled', NULL, 800000, NOW() + INTERVAL '10 minutes', '2026-07-21 09:00:00+00', '2026-07-10 07:00:00+00', '2026-07-10 07:10:00+00');

-- ============================================================
-- 7. BOOKING ITEMS
-- ============================================================
INSERT INTO booking_items (id, booking_id, inventory_unit_id, price, created_at)
SELECT '50000000-5000-5000-5000-500000000001'::uuid, '40000000-4000-4000-4000-400000000001'::uuid, iu.id, 1200000, '2026-07-10 08:00:00+00'
FROM inventory_units iu
WHERE iu.resource_id = '10000000-1000-1000-1000-100000000001'::uuid
  AND iu.unit_code = 'J1A'
LIMIT 1;

INSERT INTO booking_items (id, booking_id, inventory_unit_id, price, created_at)
SELECT '50000000-5000-5000-5000-500000000002'::uuid, '40000000-4000-4000-4000-400000000002'::uuid, iu.id, 800000, '2026-07-10 09:00:00+00'
FROM inventory_units iu
WHERE iu.resource_id = '30000000-3000-3000-3000-300000000001'::uuid
  AND iu.unit_code = 'RM-1'
LIMIT 1;

INSERT INTO booking_items (id, booking_id, inventory_unit_id, price, created_at)
SELECT '50000000-5000-5000-5000-500000000003'::uuid, '40000000-4000-4000-4000-400000000003'::uuid, iu.id, 800000, '2026-07-10 07:00:00+00'
FROM inventory_units iu
WHERE iu.resource_id = '10000000-1000-1000-1000-100000000003'::uuid
  AND iu.unit_code = 'J1A'
LIMIT 1;

-- ============================================================
-- 8. PAYMENTS
-- ============================================================
INSERT INTO payments (id, booking_id, status, amount, provider_ref, created_at, updated_at) VALUES
  ('60000000-6000-6000-4000-600000000001', '40000000-4000-4000-4000-400000000001', 'completed', 1200000, 'MOCK-PAY-001', '2026-07-10 08:05:00+00', '2026-07-10 08:05:00+00');

-- ============================================================
-- 9. IDEMPOTENCY KEYS
-- ============================================================
INSERT INTO idempotency_keys (key, request_hash, response_body, status, created_at, expires_at) VALUES
  ('demo-idempotency-key-1', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa', '{"booking_id":"40000000-4000-4000-4000-400000000001"}', 'completed', '2026-07-10 08:00:00+00', NOW() + INTERVAL '24 hours'),
  ('demo-idempotency-key-2', 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb', '{"booking_id":"40000000-4000-4000-4000-400000000002"}', 'processing', '2026-07-10 09:00:00+00', NOW() + INTERVAL '24 hours');
