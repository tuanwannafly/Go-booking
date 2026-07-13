package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidInterval = errors.New("worker interval must be greater than zero")

type Worker interface {
	Start(ctx context.Context) error
	Stop() error
	Name() string
}

type WorkerService struct {
	workers   []Worker
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	startOnce sync.Once
	stopOnce  sync.Once
}

func NewWorkerService(workers ...Worker) *WorkerService {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerService{
		workers: workers,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (s *WorkerService) Start(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	s.startOnce.Do(func() {
		s.ctx, s.cancel = context.WithCancel(ctx)
		for _, w := range s.workers {
			s.wg.Add(1)
			go func(worker Worker) {
				defer s.wg.Done()
				if err := worker.Start(s.ctx); err != nil && !errors.Is(err, context.Canceled) {
					fmt.Printf("Worker %s error: %v\n", worker.Name(), err)
				}
			}(w)
		}
	})
	return nil
}

func (s *WorkerService) Stop() error {
	s.stopOnce.Do(func() {
		if s.cancel != nil {
			s.cancel()
		}
		for _, w := range s.workers {
			if err := w.Stop(); err != nil {
				fmt.Printf("Worker %s shutdown error: %v\n", w.Name(), err)
			}
		}
		s.wg.Wait()
	})
	return nil
}

// HoldReleaseWorker releases expired holds
type HoldReleaseWorker struct {
	db       *pgxpool.Pool
	interval time.Duration
	ticker   *time.Ticker
	stopCh   chan struct{}
	stopOnce sync.Once
}

func NewHoldReleaseWorker(db *pgxpool.Pool, interval time.Duration) *HoldReleaseWorker {
	return &HoldReleaseWorker{
		db:       db,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

func (w *HoldReleaseWorker) Name() string {
	return "HoldReleaseWorker"
}

func (w *HoldReleaseWorker) Start(ctx context.Context) error {
	if w.interval <= 0 {
		return ErrInvalidInterval
	}
	w.ticker = time.NewTicker(w.interval)
	defer w.ticker.Stop()
	w.releaseExpiredHolds(ctx)

	for {
		select {
		case <-w.ticker.C:
			w.releaseExpiredHolds(ctx)
		case <-w.stopCh:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (w *HoldReleaseWorker) Stop() error {
	w.stopOnce.Do(func() { close(w.stopCh) })
	return nil
}

func (w *HoldReleaseWorker) releaseExpiredHolds(ctx context.Context) {
	if w.db == nil {
		return
	}
	query := `
		UPDATE inventory_units
		SET status = 'available', held_until = NULL, version = version + 1, updated_at = NOW()
		WHERE status = 'held' AND held_until IS NOT NULL AND held_until < NOW()
	`

	cmdTag, err := w.db.Exec(ctx, query)
	if err != nil {
		fmt.Printf("Failed to release expired holds: %v\n", err)
		return
	}

	if cmdTag.RowsAffected() > 0 {
		fmt.Printf("Released %d expired holds\n", cmdTag.RowsAffected())
	}
}

// BookingExpireWorker expires pending bookings
type BookingExpireWorker struct {
	db       *pgxpool.Pool
	interval time.Duration
	ticker   *time.Ticker
	stopCh   chan struct{}
	stopOnce sync.Once
}

func NewBookingExpireWorker(db *pgxpool.Pool, interval time.Duration) *BookingExpireWorker {
	return &BookingExpireWorker{
		db:       db,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

func (w *BookingExpireWorker) Name() string {
	return "BookingExpireWorker"
}

func (w *BookingExpireWorker) Start(ctx context.Context) error {
	if w.interval <= 0 {
		return ErrInvalidInterval
	}
	w.ticker = time.NewTicker(w.interval)
	defer w.ticker.Stop()

	for {
		select {
		case <-w.ticker.C:
			w.expirePendingBookings(ctx)
		case <-w.stopCh:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (w *BookingExpireWorker) Stop() error {
	w.stopOnce.Do(func() { close(w.stopCh) })
	return nil
}

func (w *BookingExpireWorker) expirePendingBookings(ctx context.Context) {
	// First get expired bookings
	query := `
		SELECT id FROM bookings
		WHERE status = 'pending' AND expires_at IS NOT NULL AND expires_at < NOW()
	`

	rows, err := w.db.Query(ctx, query)
	if err != nil {
		fmt.Printf("Failed to query expired bookings: %v\n", err)
		return
	}
	defer rows.Close()

	var bookingIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			fmt.Printf("Failed to scan booking ID: %v\n", err)
			continue
		}
		bookingIDs = append(bookingIDs, id)
	}

	if len(bookingIDs) == 0 {
		return
	}

	// Get inventory unit IDs for these bookings
	placeholders := ""
	args := make([]interface{}, len(bookingIDs))
	for i, id := range bookingIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query = fmt.Sprintf(`
		SELECT inventory_unit_id FROM booking_items
		WHERE booking_id IN (%s)
	`, placeholders)

	rows, err = w.db.Query(ctx, query, args...)
	if err != nil {
		fmt.Printf("Failed to query booking items: %v\n", err)
		return
	}
	defer rows.Close()

	var inventoryUnitIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			continue
		}
		inventoryUnitIDs = append(inventoryUnitIDs, id)
	}

	// Release inventory units
	if len(inventoryUnitIDs) > 0 {
		placeholders = ""
		args = make([]interface{}, len(inventoryUnitIDs))
		for i, id := range inventoryUnitIDs {
			if i > 0 {
				placeholders += ","
			}
			placeholders += fmt.Sprintf("$%d", i+1)
			args[i] = id
		}

		query = fmt.Sprintf(`
			UPDATE inventory_units
			SET status = 'available', held_until = NULL, version = version + 1, updated_at = NOW()
			WHERE id IN (%s) AND status IN ('held', 'booked')
		`, placeholders)

		_, err = w.db.Exec(ctx, query, args...)
		if err != nil {
			fmt.Printf("Failed to release inventory: %v\n", err)
		}
	}

	// Update booking status to expired
	placeholders = ""
	args = make([]interface{}, len(bookingIDs))
	for i, id := range bookingIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query = fmt.Sprintf(`
		UPDATE bookings
		SET status = 'expired', updated_at = NOW()
		WHERE id IN (%s)
	`, placeholders)

	cmdTag, err := w.db.Exec(ctx, query, args...)
	if err != nil {
		fmt.Printf("Failed to expire bookings: %v\n", err)
		return
	}

	if cmdTag.RowsAffected() > 0 {
		fmt.Printf("Expired %d pending bookings\n", cmdTag.RowsAffected())
	}
}

// IdempotencyCleanupWorker cleans up expired idempotency keys
type IdempotencyCleanupWorker struct {
	db       *pgxpool.Pool
	interval time.Duration
	ticker   *time.Ticker
	stopCh   chan struct{}
	stopOnce sync.Once
}

func NewIdempotencyCleanupWorker(db *pgxpool.Pool, interval time.Duration) *IdempotencyCleanupWorker {
	return &IdempotencyCleanupWorker{
		db:       db,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

func (w *IdempotencyCleanupWorker) Name() string {
	return "IdempotencyCleanupWorker"
}

func (w *IdempotencyCleanupWorker) Start(ctx context.Context) error {
	if w.interval <= 0 {
		return ErrInvalidInterval
	}
	w.ticker = time.NewTicker(w.interval)
	defer w.ticker.Stop()

	for {
		select {
		case <-w.ticker.C:
			w.cleanupExpiredKeys(ctx)
		case <-w.stopCh:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (w *IdempotencyCleanupWorker) Stop() error {
	w.stopOnce.Do(func() { close(w.stopCh) })
	return nil
}

func (w *IdempotencyCleanupWorker) cleanupExpiredKeys(ctx context.Context) {
	query := `DELETE FROM idempotency_keys WHERE expires_at < NOW()`
	cmdTag, err := w.db.Exec(ctx, query)
	if err != nil {
		fmt.Printf("Failed to cleanup idempotency keys: %v\n", err)
		return
	}

	if cmdTag.RowsAffected() > 0 {
		fmt.Printf("Cleaned up %d expired idempotency keys\n", cmdTag.RowsAffected())
	}
}
