-- Migration: 000009_create_payments_table.down.sql
DROP TABLE IF EXISTS payments;
DROP TYPE IF EXISTS payment_status;