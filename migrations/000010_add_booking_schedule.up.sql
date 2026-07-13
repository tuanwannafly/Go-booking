ALTER TABLE bookings ADD COLUMN scheduled_at TIMESTAMPTZ;
CREATE INDEX idx_bookings_scheduled_at ON bookings(scheduled_at) WHERE scheduled_at IS NOT NULL;
