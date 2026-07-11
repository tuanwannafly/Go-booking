package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"gobooking/internal/domain"
)

var (
	ErrIdempotencyKeyNotFound = errors.New("idempotency key not found")
	ErrIdempotencyKeyConflict = errors.New("idempotency key conflict: same key with different request")
)

type IdempotencyRepository struct {
	db *pgxpool.Pool
}

func NewIdempotencyRepository(db *pgxpool.Pool) *IdempotencyRepository {
	return &IdempotencyRepository{db: db}
}

func (r *IdempotencyRepository) Create(ctx context.Context, key *domain.IdempotencyKey) error {
	responseBody, err := json.Marshal(key.ResponseBody)
	if err != nil {
		return fmt.Errorf("failed to marshal response body: %w", err)
	}

	query := `
		INSERT INTO idempotency_keys (key, request_hash, response_body, status, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = r.db.Exec(ctx, query, key.Key, key.RequestHash, responseBody, key.Status, key.CreatedAt, key.ExpiresAt)
	if err != nil {
		// Check for unique constraint violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return ErrIdempotencyKeyConflict
		}
		return err
	}
	return nil
}

func (r *IdempotencyRepository) Get(ctx context.Context, key string) (*domain.IdempotencyKey, error) {
	query := `
		SELECT key, request_hash, response_body, status, created_at, expires_at
		FROM idempotency_keys WHERE key = $1
	`
	row := r.db.QueryRow(ctx, query, key)

	var ik domain.IdempotencyKey
	var responseBody []byte
	var expiresAt sql.NullTime
	err := row.Scan(&ik.Key, &ik.RequestHash, &responseBody, &ik.Status, &ik.CreatedAt, &expiresAt)
	if err == pgx.ErrNoRows {
		return nil, ErrIdempotencyKeyNotFound
	}
	if err != nil {
		return nil, err
	}

	if expiresAt.Valid {
		ik.ExpiresAt = expiresAt.Time
	}

	// Unmarshal response body
	if len(responseBody) > 0 {
		if err := json.Unmarshal(responseBody, &ik.ResponseBody); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	// Check if expired
	if time.Now().After(ik.ExpiresAt) {
		return nil, ErrIdempotencyKeyNotFound
	}

	return &ik, nil
}

func (r *IdempotencyRepository) Update(ctx context.Context, key *domain.IdempotencyKey) error {
	responseBody, err := json.Marshal(key.ResponseBody)
	if err != nil {
		return fmt.Errorf("failed to marshal response body: %w", err)
	}

	query := `
		UPDATE idempotency_keys
		SET response_body = $1, status = $2
		WHERE key = $3
	`
	_, err = r.db.Exec(ctx, query, responseBody, key.Status, key.Key)
	return err
}

func (r *IdempotencyRepository) DeleteExpired(ctx context.Context) (int64, error) {
	query := `DELETE FROM idempotency_keys WHERE expires_at < NOW()`
	cmdTag, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return cmdTag.RowsAffected(), nil
}
