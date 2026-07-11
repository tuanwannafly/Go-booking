-- Migration: 000005_create_inventory_units_table.up.sql
CREATE TYPE resource_type AS ENUM ('flight_seat', 'hotel_room');
CREATE TYPE inventory_status AS ENUM ('available', 'held', 'booked', 'cancelled');

CREATE TABLE inventory_units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_type resource_type NOT NULL,
    resource_id UUID NOT NULL,
    unit_code VARCHAR(50) NOT NULL,
    status inventory_status NOT NULL DEFAULT 'available',
    held_until TIMESTAMPTZ,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(resource_type, resource_id, unit_code)
);

CREATE INDEX idx_inventory_units_resource ON inventory_units(resource_type, resource_id);
CREATE INDEX idx_inventory_units_status ON inventory_units(status);
CREATE INDEX idx_inventory_units_held_until ON inventory_units(held_until) WHERE held_until IS NOT NULL;