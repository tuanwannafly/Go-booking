-- Migration: 000006_create_bookings_table.down.sql
DROP TABLE IF EXISTS bookings;
DROP TYPE IF EXISTS booking_status;