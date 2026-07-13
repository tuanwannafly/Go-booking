DROP INDEX IF EXISTS idx_bookings_scheduled_at;
ALTER TABLE bookings DROP COLUMN IF EXISTS scheduled_at;
