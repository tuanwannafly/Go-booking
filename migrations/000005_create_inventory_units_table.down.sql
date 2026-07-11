-- Migration: 000005_create_inventory_units_table.down.sql
DROP TABLE IF EXISTS inventory_units;
DROP TYPE IF EXISTS inventory_status;
DROP TYPE IF EXISTS resource_type;