-- Migration: 000007_create_booking_items_table.up.sql
CREATE TABLE booking_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    inventory_unit_id UUID NOT NULL REFERENCES inventory_units(id) ON DELETE CASCADE,
    price DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_booking_items_booking_id ON booking_items(booking_id);
CREATE INDEX idx_booking_items_inventory_unit_id ON booking_items(inventory_unit_id);