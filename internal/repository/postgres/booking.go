package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gobooking/internal/domain"
)

var (
	ErrBookingNotFound      = errors.New("booking not found")
	ErrBookingNotPending    = errors.New("booking is not in pending status")
	ErrIdempotencyKeyExists = errors.New("idempotency key already exists with different request")
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	query := `
		INSERT INTO bookings (id, user_id, status, idempotency_key, total_amount, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query, booking.ID, booking.UserID, booking.Status, booking.IdempotencyKey, booking.TotalAmount, booking.ExpiresAt, booking.CreatedAt, booking.UpdatedAt)
	return err
}

func (r *BookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	query := `
		SELECT id, user_id, status, idempotency_key, total_amount, expires_at, created_at, updated_at
		FROM bookings WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)

	var booking domain.Booking
	var expiresAt sql.NullTime
	err := row.Scan(&booking.ID, &booking.UserID, &booking.Status, &booking.IdempotencyKey, &booking.TotalAmount, &expiresAt, &booking.CreatedAt, &booking.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if expiresAt.Valid {
		booking.ExpiresAt = &expiresAt.Time
	}

	return &booking, nil
}

func (r *BookingRepository) GetByIDWithItems(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	booking, err := r.GetByID(ctx, id)
	if err != nil || booking == nil {
		return booking, err
	}

	items, err := NewBookingItemRepository(r.db).GetByBookingID(ctx, id)
	if err != nil {
		return nil, err
	}
	booking.Items = items

	return booking, nil
}

func (r *BookingRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]domain.Booking, int64, error) {
	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM bookings WHERE user_id = $1`
	var total int64
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []domain.Booking{}, 0, nil
	}

	query := `
		SELECT id, user_id, status, idempotency_key, total_amount, expires_at, created_at, updated_at
		FROM bookings WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var booking domain.Booking
		var expiresAt sql.NullTime
		err := rows.Scan(&booking.ID, &booking.UserID, &booking.Status, &booking.IdempotencyKey, &booking.TotalAmount, &expiresAt, &booking.CreatedAt, &booking.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		if expiresAt.Valid {
			booking.ExpiresAt = &expiresAt.Time
		}
		bookings = append(bookings, booking)
	}

	return bookings, total, nil
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus) error {
	query := `UPDATE bookings SET status = $1, updated_at = NOW() WHERE id = $2`
	cmdTag, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrBookingNotFound
	}
	return nil
}

func (r *BookingRepository) UpdateExpiresAt(ctx context.Context, id uuid.UUID, expiresAt *time.Time) error {
	query := `UPDATE bookings SET expires_at = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, expiresAt, id)
	return err
}

func (r *BookingRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Booking, error) {
	query := `
		SELECT id, user_id, status, idempotency_key, total_amount, expires_at, created_at, updated_at
		FROM bookings WHERE idempotency_key = $1
	`
	row := r.db.QueryRow(ctx, query, key)

	var booking domain.Booking
	var expiresAt sql.NullTime
	err := row.Scan(&booking.ID, &booking.UserID, &booking.Status, &booking.IdempotencyKey, &booking.TotalAmount, &expiresAt, &booking.CreatedAt, &booking.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if expiresAt.Valid {
		booking.ExpiresAt = &expiresAt.Time
	}

	return &booking, nil
}

func (r *BookingRepository) GetExpiredPendingBookings(ctx context.Context) ([]domain.Booking, error) {
	query := `
		SELECT id, user_id, status, idempotency_key, total_amount, expires_at, created_at, updated_at
		FROM bookings
		WHERE status = 'pending' AND expires_at IS NOT NULL AND expires_at < NOW()
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var booking domain.Booking
		var expiresAt sql.NullTime
		err := rows.Scan(&booking.ID, &booking.UserID, &booking.Status, &booking.IdempotencyKey, &booking.TotalAmount, &expiresAt, &booking.CreatedAt, &booking.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if expiresAt.Valid {
			booking.ExpiresAt = &expiresAt.Time
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

type BookingItemRepository struct {
	db *pgxpool.Pool
}

func NewBookingItemRepository(db *pgxpool.Pool) *BookingItemRepository {
	return &BookingItemRepository{db: db}
}

func (r *BookingItemRepository) Create(ctx context.Context, item *domain.BookingItem) error {
	query := `
		INSERT INTO booking_items (id, booking_id, inventory_unit_id, price, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, item.ID, item.BookingID, item.InventoryUnitID, item.Price, item.CreatedAt)
	return err
}

func (r *BookingItemRepository) GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]domain.BookingItem, error) {
	query := `
		SELECT id, booking_id, inventory_unit_id, price, created_at
		FROM booking_items WHERE booking_id = $1
	`
	rows, err := r.db.Query(ctx, query, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.BookingItem
	for rows.Next() {
		var item domain.BookingItem
		err := rows.Scan(&item.ID, &item.BookingID, &item.InventoryUnitID, &item.Price, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *BookingItemRepository) GetInventoryUnitIDsByBookingID(ctx context.Context, bookingID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT inventory_unit_id FROM booking_items WHERE booking_id = $1`
	rows, err := r.db.Query(ctx, query, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
