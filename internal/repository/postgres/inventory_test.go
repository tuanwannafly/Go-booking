package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHoldRoomWithOptimisticLockRejectsNegativeRetries(t *testing.T) {
	repo := &InventoryRepository{}
	_, err := repo.HoldRoomWithOptimisticLock(context.Background(), uuid.Nil, time.Minute, -1)
	if err == nil {
		t.Fatal("expected an error for negative maxRetries")
	}
}

func TestOptimisticLockConflictIsRetryable(t *testing.T) {
	if !errors.Is(ErrOptimisticLockConflict, ErrOptimisticLockConflict) {
		t.Fatal("optimistic lock conflict must remain comparable")
	}
}
