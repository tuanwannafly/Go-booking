package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gobooking/internal/domain"
	"gobooking/internal/repository/postgres"
	"gobooking/internal/repository/redis"
)

type HoldService struct {
	inventoryRepo *postgres.InventoryRepository
	redisClient   *redis.RedisClient
	holdDuration  time.Duration
}

func NewHoldService(
	inventoryRepo *postgres.InventoryRepository,
	redisClient *redis.RedisClient,
	holdDurationMinutes int,
) *HoldService {
	return &HoldService{
		inventoryRepo: inventoryRepo,
		redisClient:   redisClient,
		holdDuration:  time.Duration(holdDurationMinutes) * time.Minute,
	}
}

// HoldSeat uses pessimistic locking (Redis distributed lock + SELECT FOR UPDATE)
// for flight seats due to high contention during flash sales
func (s *HoldService) HoldSeat(ctx context.Context, seatID uuid.UUID, holdDurationMinutes int) (*domain.InventoryUnit, error) {
	holdDuration := time.Duration(holdDurationMinutes) * time.Minute
	if holdDuration <= 0 {
		holdDuration = s.holdDuration
	}

	// Acquire distributed lock first (fast failure at application level)
	lockName := fmt.Sprintf("hold:seat:%s", seatID.String())
	lock := s.redisClient.NewLock(lockName,
		redis.WithExpiry(holdDuration+30*time.Second), // lock expires slightly after hold
		redis.WithTries(3),
		redis.WithRetryDelay(50*time.Millisecond),
	)

	if err := lock.Lock(ctx); err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer lock.Unlock(ctx)

	// Then use pessimistic lock at database level
	return s.inventoryRepo.HoldSeatWithPessimisticLock(ctx, seatID, holdDuration)
}

// HoldRoom uses optimistic locking for hotel rooms (lower contention)
func (s *HoldService) HoldRoom(ctx context.Context, roomID uuid.UUID, holdDurationMinutes int) (*domain.InventoryUnit, error) {
	holdDuration := time.Duration(holdDurationMinutes) * time.Minute
	if holdDuration <= 0 {
		holdDuration = s.holdDuration
	}

	// For hotel rooms, we use optimistic locking with retries
	// No distributed lock needed as contention is lower and spread over time
	return s.inventoryRepo.HoldRoomWithOptimisticLock(ctx, roomID, holdDuration, 3)
}

func (s *HoldService) ReleaseHold(ctx context.Context, inventoryUnitID uuid.UUID) error {
	return s.inventoryRepo.ReleaseHold(ctx, []uuid.UUID{inventoryUnitID})
}

func (s *HoldService) ReleaseExpiredHolds(ctx context.Context) (int64, error) {
	return s.inventoryRepo.ReleaseExpiredHolds(ctx)
}
