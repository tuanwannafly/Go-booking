-- Migration: 000006_create_bookings_table.up.sql
CREATE TYPE booking_status AS ENUM ('pending', 'confirmed', 'cancelled', 'expired');

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status booking_status NOT NULL DEFAULT 'pending',
    idempotency_key VARCHAR(255) UNIQUE,
    total_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_idempotency_key ON bookings(idempotency_key) WHERE idempotency_key IS NOT NULL;
CREATE INDEX idx_bookings_expires_at ON bookings(expires_at) WHERE expires_at IS NOT NULL;