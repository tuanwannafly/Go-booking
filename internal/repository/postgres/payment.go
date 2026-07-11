package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gobooking/internal/domain"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (id, booking_id, status, amount, provider_ref, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query, payment.ID, payment.BookingID, payment.Status, payment.Amount, payment.ProviderRef, payment.CreatedAt, payment.UpdatedAt)
	return err
}

func (r *PaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	query := `SELECT id, booking_id, status, amount, provider_ref, created_at, updated_at FROM payments WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var payment domain.Payment
	err := row.Scan(&payment.ID, &payment.BookingID, &payment.Status, &payment.Amount, &payment.ProviderRef, &payment.CreatedAt, &payment.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, ErrPaymentNotFound
	}
	return &payment, err
}

func (r *PaymentRepository) GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*domain.Payment, error) {
	query := `SELECT id, booking_id, status, amount, provider_ref, created_at, updated_at FROM payments WHERE booking_id = $1`
	row := r.db.QueryRow(ctx, query, bookingID)

	var payment domain.Payment
	err := row.Scan(&payment.ID, &payment.BookingID, &payment.Status, &payment.Amount, &payment.ProviderRef, &payment.CreatedAt, &payment.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &payment, err
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus, providerRef string) error {
	query := `UPDATE payments SET status = $1, provider_ref = $2, updated_at = NOW() WHERE id = $3`
	cmdTag, err := r.db.Exec(ctx, query, status, providerRef, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrPaymentNotFound
	}
	return nil
}

func (r *PaymentRepository) UpdateStatusByBookingID(ctx context.Context, bookingID uuid.UUID, status domain.PaymentStatus, providerRef string) error {
	query := `UPDATE payments SET status = $1, provider_ref = $2, updated_at = NOW() WHERE booking_id = $3`
	cmdTag, err := r.db.Exec(ctx, query, status, providerRef, bookingID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrPaymentNotFound
	}
	return nil
}
