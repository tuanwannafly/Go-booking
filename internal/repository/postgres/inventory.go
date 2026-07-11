package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gobooking/internal/domain"
)

var (
	ErrInventoryNotFound      = errors.New("inventory unit not found")
	ErrInventoryNotAvailable  = errors.New("inventory unit not available")
	ErrOptimisticLockConflict = errors.New("optimistic lock conflict: version mismatch")
)

type InventoryRepository struct {
	db *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(ctx context.Context, unit *domain.InventoryUnit) error {
	query := `
		INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, held_until, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query, unit.ID, unit.ResourceType, unit.ResourceID, unit.UnitCode, unit.Status, unit.HeldUntil, unit.Version, unit.CreatedAt, unit.UpdatedAt)
	return err
}

func (r *InventoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.InventoryUnit, error) {
	query := `
		SELECT id, resource_type, resource_id, unit_code, status, held_until, version, created_at, updated_at
		FROM inventory_units WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)

	var unit domain.InventoryUnit
	var heldUntil sql.NullTime
	err := row.Scan(&unit.ID, &unit.ResourceType, &unit.ResourceID, &unit.UnitCode, &unit.Status, &heldUntil, &unit.Version, &unit.CreatedAt, &unit.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if heldUntil.Valid {
		unit.HeldUntil = &heldUntil.Time
	}

	return &unit, nil
}

func (r *InventoryRepository) GetByResource(ctx context.Context, resourceType domain.ResourceType, resourceID uuid.UUID) ([]domain.InventoryUnit, error) {
	query := `
		SELECT id, resource_type, resource_id, unit_code, status, held_until, version, created_at, updated_at
		FROM inventory_units WHERE resource_type = $1 AND resource_id = $2
	`
	rows, err := r.db.Query(ctx, query, resourceType, resourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []domain.InventoryUnit
	for rows.Next() {
		var unit domain.InventoryUnit
		var heldUntil sql.NullTime
		err := rows.Scan(&unit.ID, &unit.ResourceType, &unit.ResourceID, &unit.UnitCode, &unit.Status, &heldUntil, &unit.Version, &unit.CreatedAt, &unit.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if heldUntil.Valid {
			unit.HeldUntil = &heldUntil.Time
		}
		units = append(units, unit)
	}

	return units, nil
}

func (r *InventoryRepository) GetAvailableByResource(ctx context.Context, resourceType domain.ResourceType, resourceID uuid.UUID) ([]domain.InventoryUnit, error) {
	query := `
		SELECT id, resource_type, resource_id, unit_code, status, held_until, version, created_at, updated_at
		FROM inventory_units
		WHERE resource_type = $1 AND resource_id = $2 AND status = 'available'
	`
	rows, err := r.db.Query(ctx, query, resourceType, resourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []domain.InventoryUnit
	for rows.Next() {
		var unit domain.InventoryUnit
		var heldUntil sql.NullTime
		err := rows.Scan(&unit.ID, &unit.ResourceType, &unit.ResourceID, &unit.UnitCode, &unit.Status, &heldUntil, &unit.Version, &unit.CreatedAt, &unit.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if heldUntil.Valid {
			unit.HeldUntil = &heldUntil.Time
		}
		units = append(units, unit)
	}

	return units, nil
}

// HoldSeatWithPessimisticLock uses SELECT FOR UPDATE for flight seats (high contention)
func (r *InventoryRepository) HoldSeatWithPessimisticLock(ctx context.Context, seatID uuid.UUID, holdDuration time.Duration) (*domain.InventoryUnit, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// SELECT FOR UPDATE - pessimistic lock at database level
	query := `
		SELECT id, resource_type, resource_id, unit_code, status, held_until, version, created_at, updated_at
		FROM inventory_units
		WHERE id = $1 AND resource_type = 'flight_seat'
		FOR UPDATE
	`
	row := tx.QueryRow(ctx, query, seatID)

	var unit domain.InventoryUnit
	var heldUntil sql.NullTime
	err = row.Scan(&unit.ID, &unit.ResourceType, &unit.ResourceID, &unit.UnitCode, &unit.Status, &heldUntil, &unit.Version, &unit.CreatedAt, &unit.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, ErrInventoryNotFound
	}
	if err != nil {
		return nil, err
	}

	if heldUntil.Valid {
		unit.HeldUntil = &heldUntil.Time
	}

	if unit.Status != domain.InventoryStatusAvailable &&
		!(unit.Status == domain.InventoryStatusHeld && unit.HeldUntil != nil && unit.HeldUntil.Before(time.Now())) {
		return nil, ErrInventoryNotAvailable
	}

	// Update to held status
	heldUntilTime := time.Now().Add(holdDuration)
	updateQuery := `
		UPDATE inventory_units
		SET status = 'held', held_until = $1, version = version + 1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.Exec(ctx, updateQuery, heldUntilTime, seatID)
	if err != nil {
		return nil, err
	}

	unit.Status = domain.InventoryStatusHeld
	unit.HeldUntil = &heldUntilTime
	unit.Version++

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &unit, nil
}

// HoldRoomWithOptimisticLock uses optimistic locking for hotel rooms (lower contention)
func (r *InventoryRepository) HoldRoomWithOptimisticLock(ctx context.Context, roomID uuid.UUID, holdDuration time.Duration, maxRetries int) (*domain.InventoryUnit, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		unit, err := r.tryHoldRoomOptimistic(ctx, roomID, holdDuration)
		if err == nil {
			return unit, nil
		}

		lastErr = err
		if errors.Is(err, ErrOptimisticLockConflict) {
			// Retry on version conflict
			time.Sleep(time.Duration(attempt+1) * 10 * time.Millisecond)
			continue
		}

		// Non-retryable error
		return nil, err
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func (r *InventoryRepository) tryHoldRoomOptimistic(ctx context.Context, roomID uuid.UUID, holdDuration time.Duration) (*domain.InventoryUnit, error) {
	// First, get the current version
	query := `
		SELECT id, resource_type, resource_id, unit_code, status, held_until, version, created_at, updated_at
		FROM inventory_units WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, roomID)

	var unit domain.InventoryUnit
	var heldUntil sql.NullTime
	err := row.Scan(&unit.ID, &unit.ResourceType, &unit.ResourceID, &unit.UnitCode, &unit.Status, &heldUntil, &unit.Version, &unit.CreatedAt, &unit.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, ErrInventoryNotFound
	}
	if err != nil {
		return nil, err
	}

	if heldUntil.Valid {
		unit.HeldUntil = &heldUntil.Time
	}

	// Check if available
	if unit.Status != domain.InventoryStatusAvailable {
		// Check if hold expired
		if unit.HeldUntil != nil && unit.HeldUntil.Before(time.Now()) {
			// Expired hold, can proceed
		} else {
			return nil, ErrInventoryNotAvailable
		}
	}

	// Try to update with version check (optimistic lock)
	heldUntilTime := time.Now().Add(holdDuration)
	updateQuery := `
		UPDATE inventory_units
		SET status = 'held', held_until = $1, version = version + 1, updated_at = NOW()
		WHERE id = $2 AND version = $3
	`
	cmdTag, err := r.db.Exec(ctx, updateQuery, heldUntilTime, roomID, unit.Version)
	if err != nil {
		return nil, err
	}

	if cmdTag.RowsAffected() == 0 {
		return nil, ErrOptimisticLockConflict
	}

	unit.Status = domain.InventoryStatusHeld
	unit.HeldUntil = &heldUntilTime
	unit.Version++

	return &unit, nil
}

func (r *InventoryRepository) ReleaseHold(ctx context.Context, inventoryUnitIDs []uuid.UUID) error {
	if len(inventoryUnitIDs) == 0 {
		return nil
	}

	query := `
		UPDATE inventory_units
		SET status = 'available', held_until = NULL, version = version + 1, updated_at = NOW()
		WHERE id = ANY($1) AND status = 'held'
	`
	_, err := r.db.Exec(ctx, query, inventoryUnitIDs)
	return err
}

func (r *InventoryRepository) ReleaseExpiredHolds(ctx context.Context) (int64, error) {
	query := `
		UPDATE inventory_units
		SET status = 'available', held_until = NULL, version = version + 1, updated_at = NOW()
		WHERE status = 'held' AND held_until IS NOT NULL AND held_until < NOW()
	`
	cmdTag, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return cmdTag.RowsAffected(), nil
}

func (r *InventoryRepository) ConfirmBooking(ctx context.Context, inventoryUnitIDs []uuid.UUID) error {
	if len(inventoryUnitIDs) == 0 {
		return nil
	}

	query := `
		UPDATE inventory_units
		SET status = 'booked', held_until = NULL, version = version + 1, updated_at = NOW()
		WHERE id = ANY($1) AND status = 'held'
	`
	cmdTag, err := r.db.Exec(ctx, query, inventoryUnitIDs)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() != int64(len(inventoryUnitIDs)) {
		return errors.New("some inventory units were not in held status")
	}
	return nil
}

func (r *InventoryRepository) CancelBooking(ctx context.Context, inventoryUnitIDs []uuid.UUID) error {
	if len(inventoryUnitIDs) == 0 {
		return nil
	}

	query := `
		UPDATE inventory_units
		SET status = 'cancelled', held_until = NULL, version = version + 1, updated_at = NOW()
		WHERE id = ANY($1) AND status = 'booked'
	`
	_, err := r.db.Exec(ctx, query, inventoryUnitIDs)
	return err
}

func (r *InventoryRepository) ExpireBooking(ctx context.Context, inventoryUnitIDs []uuid.UUID) error {
	if len(inventoryUnitIDs) == 0 {
		return nil
	}

	query := `
		UPDATE inventory_units
		SET status = 'available', held_until = NULL, version = version + 1, updated_at = NOW()
		WHERE id = ANY($1) AND status IN ('held', 'booked')
	`
	_, err := r.db.Exec(ctx, query, inventoryUnitIDs)
	return err
}

func (r *InventoryRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.InventoryUnit, error) {
	if len(ids) == 0 {
		return []domain.InventoryUnit{}, nil
	}
	query := `SELECT id, resource_type, resource_id, unit_code, status, held_until, version, created_at, updated_at FROM inventory_units WHERE id = ANY($1)`
	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []domain.InventoryUnit
	for rows.Next() {
		var unit domain.InventoryUnit
		var heldUntil sql.NullTime
		if err := rows.Scan(&unit.ID, &unit.ResourceType, &unit.ResourceID, &unit.UnitCode, &unit.Status, &heldUntil, &unit.Version, &unit.CreatedAt, &unit.UpdatedAt); err != nil {
			return nil, err
		}
		if heldUntil.Valid {
			unit.HeldUntil = &heldUntil.Time
		}
		units = append(units, unit)
	}
	return units, rows.Err()
}
